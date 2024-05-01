# CsvObject
CsvObject package allows to dump structured data to CSV ~~file~~ string, 
and parses CSV data from string to slice of structures.

In order to minimize dependencies, this package does not do direct write to file (filesystem). 
It only deals with (CSV) strings.

Limitations: CSV data can be represented as an arbitrary structure, 
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