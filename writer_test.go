package ndjson

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriter(t *testing.T) {
	type Doc struct {
		ID int64 `json:"id"`
	}

	tests := []struct {
		Input  []Doc
		Output []byte
		Error  string
	}{
		// #0
		{
			Input:  nil,
			Output: nil,
		},
		// #1
		{
			Input: []Doc{
				{ID: 1},
				{ID: 2},
			},
			Output: []byte(`{"id":1}
{"id":2}`),
		},
	}

	for i, tt := range tests {
		var out bytes.Buffer
		w := NewWriter(&out)
		for _, doc := range tt.Input {
			if err := w.Encode(doc); err != nil {
				if tt.Error == "" {
					t.Fatalf("#%d. expected no error, got %v", i, err)
				}
				if want, have := tt.Error, err.Error(); !strings.Contains(have, want) {
					t.Fatalf("#%d. want Error=~%q, have %q", i, want, have)
				}
			}
		}
	}
}
