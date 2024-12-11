package main

import (
	"fmt"
	"github.com/anxp/csvobject"
	"math/big"
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
	Stipend           *big.Int // Yeah, our students are really rich
}

func main() {
	rawCsvData := "StudentID|int64," +
		"FirstName|string," +
		"LastName|string," +
		"EntryYear|int," +
		"Speciality|string," +
		"AdditionalData.IsCovidVaccinated|bool," +
		"AdditionalData.HasDriveLicense|bool," +
		"AdditionalData.IsPaidStudent|bool," +
		"AdditionalData.Stipend|*big.Int\n" +
		"1,John,Smith,2010,Physics,true,true,false,1000000000000000000\n" + // John have stipend 1 ETH in Wei (1 ETH = 1*10^18 Wei)
		"2,Max,Mustermann,2011,Electronics,false,true,true,2000000000000000000\n" + // Max have stipend 2 ETH in Wei
		"3,Peter,Davis,2012,Mathematics,true,false,true,3000000000000000000" // Peter have stipend 3 ETH in Wei

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
