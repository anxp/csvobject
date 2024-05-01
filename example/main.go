package main

import (
	"csvobject"
	"fmt"
)

type Student struct {
	StudentID      int64
	FirstName      string
	LastName       string
	EntryYear      int
	Speciality     string
	AdditionalData AdditionalData
}

type AdditionalData struct {
	IsCovidVaccinated bool
	HasDriveLicense   bool
	IsPaidStudent     bool
}

func main() {
	rawCsvData := "StudentID|int64," +
		"FirstName|string," +
		"LastName|string," +
		"EntryYear|int," +
		"Speciality|string," +
		"AdditionalData.IsCovidVaccinated|bool," +
		"AdditionalData.HasDriveLicense|bool," +
		"AdditionalData.IsPaidStudent|bool" +
		"\n" +
		"1,John,Smith,2010,Physics,true,true,false" +
		"\n" +
		"2,Max,Mustermann,2011,Electronics,false,true,true" +
		"\n" +
		"3,John,Davis,2012,Mathematics,true,false,true"

	tmpStudents, err := csvobject.ImportRawCSV(Student{}, rawCsvData)

	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	ourStudents := make([]Student, len(tmpStudents))

	for i, student := range tmpStudents {
		ourStudents[i] = student.(Student)
	}

	fmt.Printf("%v", ourStudents)

	return
}
