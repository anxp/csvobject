package csvobject

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

type HeaderSType struct {
	fieldsHierarchy []string
	scalarType      string
}

// ImportRawCSV imports CSV data represented as a one huge string in array of objects.
// Returns slice, where each element represents structured (parsed to structure) data row.
// Type of expected object (structure into which we will parse) can be passed by passing null-object as first argument of this method.
func ImportRawCSV(csvRowNullObject interface{}, csvRawData string) ([]interface{}, error) {
	csvRows := strings.Split(csvRawData, "\n")
	csvRows = filterEmptyStrings(csvRows)

	return ImportStringsCSV(csvRowNullObject, csvRows)
}

// ImportStringsCSV imports CSV data represented as a slice of CSV strings.
// Returns slice, where each element represents structured (parsed to structure) data row.
// Type of expected object (structure into which we will parse) can be passed by passing null-object as first argument of this method.
func ImportStringsCSV(csvRowNullObject interface{}, csvRows []string) ([]interface{}, error) {
	csvRows = filterEmptyStrings(csvRows)

	if len(csvRows) == 0 {
		return nil, errors.New("empty input data")
	}

	headerLine := csvRows[0]
	columnHeaders := strings.Split(headerLine, ",")
	headerTypes := make([]HeaderSType, len(columnHeaders))
	csvOutputDataStructure := make([]interface{}, len(csvRows)-1) // -1 because first line is header line

	for i := 0; i < len(columnHeaders); i++ {
		header := strings.Split(columnHeaders[i], "|")
		fieldsHierarchy := strings.Split(header[0], ".")
		scalarType := header[1]

		headerTypes[i] = HeaderSType{
			fieldsHierarchy: fieldsHierarchy,
			scalarType:      scalarType,
		}
	}

	for i := 1; i < len(csvRows); i++ {
		currentLineValues := strings.Split(csvRows[i], ",")

		if len(currentLineValues) != len(columnHeaders) {
			return nil, fmt.Errorf("csv data damaged, expected [%d] columns, found [%d], at line [%d]", len(columnHeaders), len(currentLineValues), i)
		}

		// Magic explanation:
		// https://stackoverflow.com/questions/63421976/panic-reflect-call-of-reflect-value-fieldbyname-on-interface-value

		// Create a FRESH COPY of line object for each new line
		csvLineObject := csvRowNullObject

		// reflectedValue is the interface{}
		reflectedValue := reflect.ValueOf(&csvLineObject).Elem()

		// Allocate a temporary variable with type of the struct.
		// reflectedValue.Elem() is the vale contained in the interface.
		tmpLine := reflect.New(reflectedValue.Elem().Type()).Elem()

		for colNo, cellValue := range currentLineValues {
			// At this point we have a VALUE in cellValue variable in STRING format
			// We need to convert this value to desired format.
			fieldsHierarchy := headerTypes[colNo].fieldsHierarchy
			scalarType := headerTypes[colNo].scalarType

			valueTyped, err := convertStringToTypedValue(cellValue, scalarType)

			if err != nil {
				return nil, err
			}

			tmpCell := tmpLine

			for j := 0; j < len(fieldsHierarchy); j++ {
				tmpCell = tmpCell.FieldByName(fieldsHierarchy[j])
			}

			tmpCell.Set(reflect.ValueOf(valueTyped))
		}

		// Set the interface to the modified struct value.
		reflectedValue.Set(tmpLine)
		csvOutputDataStructure[i-1] = csvLineObject
	}

	return csvOutputDataStructure, nil
}

// ExportFullDataToCSV exports array of structures to CSV data. Each structure in array represents single CSV line.
func ExportFullDataToCSV(csvRows []interface{}) (string, error) {
	if len(csvRows) == 0 {
		return "", errors.New("no input data")
	}

	var content string

	for i, structuredRow := range csvRows {
		header, payload, err := structureToCSVString(structuredRow, "")

		if err != nil {
			return "", err
		}

		if i == 0 {
			content = header + "\n"
		}

		content += payload + "\n"
	}

	return content, nil
}

// ExportDataRowToCSV exports ONLY DATA from given object (structure) to CSV line.
// EOL character ("\n") is optional.
func ExportDataRowToCSV(csvRow interface{}, addEOL bool) (string, error) {
	_, contentRow, err := structureToCSVString(csvRow, "")

	if err != nil {
		return "", err
	}

	if contentRow == "" {
		return "", errors.New("invalid input data")
	}

	if addEOL {
		contentRow += "\n"
	}

	return contentRow, nil
}

// ExportHeaderToCSV exports ONLY HEADER / FIELD NAMES from given object (structure) to CSV line.
// EOL character ("\n") is optional.
func ExportHeaderToCSV(csvRow interface{}, addEOL bool) (string, error) {
	header, _, err := structureToCSVString(csvRow, "")

	if err != nil {
		return "", err
	}

	if header == "" {
		return "", errors.New("invalid input data")
	}

	if addEOL {
		header += "\n"
	}

	return header, nil
}

func structureToCSVString(csvRowObject interface{}, parentFieldName string) (string, string, error) {
	reflectedValue := reflect.ValueOf(csvRowObject)

	fieldNames := make([]string, reflectedValue.NumField())
	fieldValues := make([]string, reflectedValue.NumField())

	for i := 0; i < reflectedValue.NumField(); i++ {
		nameStr := fmt.Sprintf("%s", reflectedValue.Type().Field(i).Name)

		if parentFieldName != "" {
			nameStr = parentFieldName + "." + nameStr
		}

		var err error
		valueStr := ""

		if reflectedValue.Field(i).Kind() == reflect.Struct {
			nameStr, valueStr, err = structureToCSVString(reflectedValue.Field(i).Interface(), nameStr)

			if err != nil {
				return "", "", nil
			}
		} else if reflectedValue.Field(i).Kind().String() == reflectedValue.Field(i).Type().String() || // Check for standard scalar types (int, float, etc)
			(reflectedValue.Field(i).Kind().String() == "ptr" && reflectedValue.Field(i).Type().String() == "*big.Int") { // Check for big.Int type

			valueStr = fmt.Sprintf("%v", reflectedValue.Field(i).Interface())
			nameStr = nameStr + "|" + reflectedValue.Field(i).Type().String() // Why Type()? For int and float type and kind are the same; but for big.Int type = "*big.Int", kind = "ptr"

		} else {
			return "", "", fmt.Errorf(
				"only standard scalar types allowed, this one is not standard; Type = \"%s\", Kind = \"%s\"",
				reflectedValue.Field(i).Type().String(),
				reflectedValue.Field(i).Kind().String(),
			)
		}

		fieldNames[i] = nameStr
		fieldValues[i] = valueStr
	}

	return strings.Join(fieldNames, ","), strings.Join(fieldValues, ","), nil
}

func convertStringToTypedValue(strValue string, desiredType string) (interface{}, error) {

	var result interface{}
	var err error

	switch desiredType {

	case "string":
		result, err = strValue, nil

	case "int8":
		result, err = strconv.ParseInt(strValue, 10, 8)

	case "int":
		result, err = strconv.Atoi(strValue)

	case "int32":
		result, err = strconv.ParseInt(strValue, 10, 32)

	case "int64":
		result, err = strconv.ParseInt(strValue, 10, 64)

	case "float32":
		result, err = strconv.ParseFloat(strValue, 32)

	case "float64":
		result, err = strconv.ParseFloat(strValue, 64)

	case "bool":
		result, err = strconv.ParseBool(strValue)

	case "*big.Int":
		ok := false

		result, ok = big.NewInt(0).SetString(strValue, 10)
		if !ok {
			return nil, errors.New("failed to convert string to big.Int value")
		}

	default:
		return nil, errors.New("unsupported value type: \"" + desiredType + "\"")
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func filterEmptyStrings(lines []string) []string {
	result := make([]string, 0, len(lines))

	for _, str := range lines {
		if str != "" {
			result = append(result, str)
		}
	}

	return result
}
