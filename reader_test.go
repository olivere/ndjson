package ndjson

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReader(t *testing.T) {
	type Doc struct {
		ID int64 `json:"id"`
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
