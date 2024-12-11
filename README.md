# CsvObject
## 📖 Description
CsvObject package allows to dump structured data to CSV ~~file~~ string, 
and vise versa - parse CSV data from string to slice of structures.

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
- *big.Int (this type added 11.12.204 to expand the application of the module to crypto projects)

## 📦 Installation
To make CsvObject available in your project, you can run the following command.
Make sure to run this command inside your project, when you're using go modules 😉

```sh 
go get github.com/anxp/csvobject
```

## 👀 Example
For example, if we have in our code some structured data like:

```
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
```

And we have series of such data (slice):

![CSV data parsed to structure](screenshot.png?raw=true "Title")

And we want to have the ability to **easily** store such data on disk and then read back to structured format, csvobject library can be used for that purpose.

Data will be converted to CSV format, with a specially formatted header, where original types of data are preserved:

| StudentID&#124;string | FirstName&#124;string | LastName&#124;string | EntryYear&#124;int | Speciality&#124;string | AdditionalData.IsCovidVaccinated&#124;bool | AdditionalData.HasDriveLicense&#124;bool | AdditionalData.IsPaidStudent&#124;bool | AdditionalData.Stipend&#124;*big.Int |
|-----------------------|-----------------------|----------------------|--------------------|------------------------|--------------------------------------------|------------------------------------------|----------------------------------------|--------------------------------------|
| 1                     | John                  | Smith                | 2010               | Physics                | true                                       | true                                     | false                                  | 1000000000000000000                  |
| 2                     | Max                   | Mustermann           | 2011               | Electronics            | false                                      | true                                     | true                                   | 2000000000000000000                  |
| 3                     | John                  | Davis                | 2012               | Mathematics            | true                                       | false                                    | true                                   | 3000000000000000000                  |

When we read such data back from disk and want to parse to original structured format, column headers are used to get original structure field names **and their type**,
what allows to completely reconstruct objects in runtime.
