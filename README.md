# CsvObject
## 📖 Description
CsvObject package allows to dump structured data to CSV ~~file~~ string, 
and parses CSV data from string to slice of structures.

In order to minimize dependencies, this package does not do direct write to file (filesystem). 
It only deals with (CSV) strings.

### Limitations:
CSV data can be represented as an arbitrary structure, 
where fields can represent other structures (with unlimited nesting), 
**but every "row" should be "flat", or in other words every last field in nested structure 
chain should be scalar and have standard type**: 
- string
- int8
- int
- int32
- int64
- float32
- float64
- bool

## 📦 Installation
To make CsvObject available in your project, you can run the following command.
Make sure to run this command inside your project, when you're using go modules 😉

```sh 
go get github.com/anxp/csvobject
```

## 👀 Example
For example, if we have CSV data like this 
(note that the AdditionalData field is another nested structure):

| StudentID&#124;string | FirstName&#124;string | LastName&#124;string | EntryYear&#124;int | Speciality&#124;string | AdditionalData.IsCovidVaccinated&#124;bool | AdditionalData.HasDriveLicense&#124;bool | AdditionalData.IsPaidStudent&#124;bool |
|-----------------------|-----------------------|----------------------|--------------------|------------------------|--------------------------------------------|------------------------------------------|----------------------------------------|
| 1                     | John                  | Smith                | 2010               | Physics                | true                                       | true                                     | false                                  |
| 2                     | Max                   | Mustermann           | 2011               | Electronics            | false                                      | true                                     | true                                   |
| 3                     | John                  | Davis                | 2012               | Mathematics            | true                                       | false                                    | true                                   |

*Please note, data above ☝ should be standard CSV data (columns separated by comma, 
rows terminated by \\n\). Here they rendered as a table just for clarity only*

Then, we can parse such data in structured format like this (screenshot from GoLand debugger): 

![CSV data parsed to structure](screenshot.png?raw=true "Title")