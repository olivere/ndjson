# What is it?

[![Build Status](https://github.com/olivere/ndjson/workflows/Test/badge.svg)](https://github.com/olivere/ndjson/actions)

The `ndjson` package implements reading and writing files according
to the [ndjson specification](http://ndjson.org/).

Example:

```go
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
```

# License

MIT. See [LICENSE](https://github.com/olivere/ndjson/blob/master/LICENSE) file.
