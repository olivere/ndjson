package ndjson_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/olivere/ndjson"
)

type Location struct {
	City string `json:"city"`
}

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

func ExampleReader_Bytes() {
	r := ndjson.NewReader(strings.NewReader(`{"city":"Munich"}
{"city":"Invalid"
{"city":"London"}`))
	for r.Next() {
		var loc Location
		if err := r.Decode(&loc); err != nil {
			fmt.Printf("Decode failed: %v. Last read: %s\n", err, string(r.Bytes()))
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
	// Decode failed: unexpected end of JSON input. Last read: {"city":"Invalid"
}

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
		}
	}
	// Output:
	// {"city":"Munich"}
	// {"city":"Berlin"}
	// {"city":"London"}
}
