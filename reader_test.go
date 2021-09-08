package ndjson

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func text(size int) []byte {
	letters := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	buf := make([]byte, size)
	for i := range buf {
		buf[i] = letters[i%len(letters)]
	}
	return buf
}

func docs(sizes ...int) []byte {
	var buf []byte
	for i := 0; i < len(sizes); i++ {
		if i > 0 {
			buf = append(buf, '\n')
		}
		buf = append(buf, fmt.Sprintf(`{"id":%d,"text":"`, i+1)...)
		buf = append(buf, text(sizes[i])...)
		buf = append(buf, `"}`...)
	}
	return buf
}

func TestReader(t *testing.T) {
	type Doc struct {
		ID   int64  `json:"id"`
		Text string `json:"text,omitempty"`
	}

	tests := []struct {
		Input  []byte
		Output []Doc
		Error  string
	}{
		// #0
		{
			Input:  nil,
			Output: nil,
		},
		// #1
		{
			Input: []byte(`{"id":1}
{"id":2}`),
			Output: []Doc{
				{ID: 1},
				{ID: 2},
			},
		},
		// #2
		{
			Input:  []byte(`{"id":"abc"}`),
			Output: nil,
			Error:  "json: cannot unmarshal string into Go struct field",
		},
		// #3
		{
			Input: []byte("{\"id\":1,\"text\":\"A room\\nwith\\na\\nnewline\\n\"}\n{\"id\":2,\"text\":\"No\\tsuch\\ntext\\r\\n\\r\\n\"}\n"),
			Output: []Doc{
				{ID: 1, Text: `A room
with
a
newline
`},
				{ID: 2, Text: "No\tsuch\ntext\r\n\r\n"},
			},
		},
		// #4
		{
			Input: docs(128*1024),
			Output: nil,
			Error: "bufio.Scanner: token too long",
		},
	}

	for i, tt := range tests {
		r := NewReader(bytes.NewReader(tt.Input))
		var n int
		for r.Next() {
			var doc Doc
			if err := r.Decode(&doc); err != nil {
				if tt.Error == "" {
					t.Fatalf("#%d. expected no error, got %v", i, err)
				}
				if want, have := tt.Error, err.Error(); !strings.Contains(have, want) {
					t.Fatalf("#%d. want Error=~%q, have %q", i, want, have)
				}
			} else {
				if want, have := tt.Output[n], doc; !cmp.Equal(want, have) {
					t.Fatalf("#%d. want Doc=%v, have %v", i, want, cmp.Diff(want, have))
				}
			}
			n++
		}
		if err := r.Err(); err != nil {
			if tt.Error == "" {
				t.Fatalf("#%d. expected no error, got %v", i, err)
			}
			if want, have := tt.Error, err.Error(); !strings.Contains(have, want) {
				t.Fatalf("#%d. want Error=~%q, have %q", i, want, have)
			}
		}
	}
}

func TestReaderSize(t *testing.T) {
	type Doc struct {
		ID   int64  `json:"id"`
		Text string `json:"text,omitempty"`
	}

	input := docs(128*1024, 10)
	output := []Doc{
		{ID: 1, Text: string(text(128*1024))},
		{ID: 2, Text: string(text(10))},
	}

	r := NewReaderSize(bytes.NewReader(input), 256*1024)
	var n int
	for r.Next() {
		var doc Doc
		if err := r.Decode(&doc); err != nil {
			t.Fatalf("expected no error, got %v", err)
		} else {
			if want, have := output[n], doc; !cmp.Equal(want, have) {
				t.Fatalf("want Doc=%v, have %v", want, cmp.Diff(want, have))
			}
		}
		n++
	}
	if err := r.Err(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
