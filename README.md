# What is it?

[![Build Status](https://github.com/olivere/ndjson/workflows/Test/badge.svg)](https://github.com/olivere/ndjson/actions)

## Ndjson
The `ndjson` package implements reading and writing files according
to the [ndjson specification](http://ndjson.org/).

## Documentation

This Go library allows reading and writing files in the newline delimited JSON (NDJSON) format, which is a simple way to represent structured data that can be read and written efficiently.

The package provides two main types, ndjson.Reader and ndjson.Writer, that allow reading and writing NDJSON data, respectively. It also provides a method for decoding NDJSON data into Go structs and encoding Go structs into NDJSON data.

The example code provided demonstrates how to use the ndjson package to read and write NDJSON files, as well as encode and decode Go structs to and from NDJSON data.

## Installation
To install this package, you can use the go get command:
```bash
go get github.com/olivere/ndjson
```

## Usage
```go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/olivere/ndjson"
)

type Location struct {
	City string `json:"city"`
}

func main() {
	// Writing to an ndjson file
	locations := []Location{
		{City: "Munich"},
		{City: "Berlin"},
		{City: "London"},
	}

	file, err := os.Create("locations.ndjson")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := ndjson.NewWriter(file)
	for _, loc := range locations {
		if err := writer.Encode(loc); err != nil {
			panic(err)
		}
	}
	writer.Flush()

	// Reading from an ndjson file
	file, err = os.Open("locations.ndjson")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := ndjson.NewReader(file)
	for reader.Next() {
		var loc Location
		if err := reader.Decode(&loc); err != nil {
			panic(err)
		}
		fmt.Println(loc.City)
	}
	if err := reader.Err(); err != nil {
		panic(err)
	}
}

```
This example writes three `Location` objects to an ndjson file called `locations.ndjson`, and then reads them back from the file, printing the `City` field of each object to the console.

Example:

```go
package ndjson_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/olivere/ndjson"
)

type Location struct {
	City string `json:"city"`
}


type collection struct {
	Index map[string]interface{} `json:"index"`
}

// ExampleReader demonstrates how to use ndjson.Reader to read NDJSON data from a string and decode it into a Go struct.
func ExampleReader() {
	r := ndjson.NewReader(strings.NewReader(`{"city":"Munich"}
{"city":"Berlin"}
{"city":"London"}`))
	for r.Next() {
		var loc Location
		if err := r.Decode(&loc); err != nil {
			fmt.Fprintf(os.Stderr, "Decode failed: %v", err)
			return
		}
		fmt.Println(loc.City)
	}
	if err := r.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Reader failed: %v", err)
		return
	}
	// Output:
	// Munich
	// Berlin
	// London
}

// ExampleWriter demonstrates how to use ndjson.Writer to encode Go structs into NDJSON data and write it to an output stream.
func ExampleWriter() {
	locations := []Location{
		{City: "Munich"},
		{City: "Berlin"},
		{City: "London"},
	}
	r := ndjson.NewWriter(os.Stdout)
	for _, loc := range locations {
		if err := r.Encode(loc); err != nil {
			fmt.Fprintf(os.Stderr, "Encode failed: %v", err)
			return
			// If an error occurs while encoding, the function stops writing and returns an error message.
		}
	}
	// Output:
	// {"city":"Munich"}
	// {"city":"Berlin"}
	// {"city":"London"}
}

// ExampleWriter_zincSearch demonstrates how to use ndjson.Writer to encode Go structs and write NDJSON data to an output stream
// for use with the ZincSearch engine.
func ExampleWriter_zincSearch() {
	locations := []Location{
		{City: "Munich"},
	}

	index := map[string]interface{}{"_index": "States"}
	collections := []collection{{Index: index}}

	var buf bytes.Buffer
	r := ndjson.NewWriter(&buf)

	for _, c := range collections {
		if err := r.Encode(c); err != nil {
			fmt.Fprintf(os.Stderr, "Encode failed: %v", err)
			return
			// If an error occurs while encoding, the function stops writing and returns an error message.
		}
	}

	for _, loc := range locations {
		if err := r.Encode(loc); err != nil {
			fmt.Fprintf(os.Stderr, "Encode failed: %v", err)
			return
			// If an error occurs while encoding, the function stops writing and returns an error message.
		}
	}
}

// Output:
// { "index" : { "_index" : "States" } }
// {"city":"Munich"}

```

The ExampleWriter_zincSearch demonstrates how to use the ndjson package to write multiple JSON documents into a single NDJSON stream. In this example, we have a slice of Location structs and a slice of collection structs.

Before writing the Location documents, we first write the collection documents that define Elasticsearch indices. Each collection struct includes metadata for the index, such as the name of the index ("_index") and its value ("States").

Once the collection documents are written, we iterate over the Location slice and encode each struct as an NDJSON document using the Encode method of the ndjson.Writer. The resulting NDJSON stream can be used to index the Location documents into an Elasticsearch index.

The expected output of the ExampleWriter_zincSearch function is the NDJSON stream containing the collection and Location documents, formatted for Elasticsearch bulk indexing.

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

## License

MIT. See [LICENSE](https://github.com/olivere/ndjson/blob/master/LICENSE) file.
