package ndjson

import (
	"encoding/json"
	"io"
)

var (
	_ io.Writer = (*Writer)(nil)
)

// Writer implements writing line-oriented JSON data following the
// ndjson spec at http://ndjson.org/.
type Writer struct {
	w   io.Writer
	enc *json.Encoder
}

// NewWriter returns a new Writer, using the underlying writer for output.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:   w,
		enc: json.NewEncoder(w),
	}
}

// Write writes the contents of p into the buffer. It returns the number of
// bytes written. if n < len(p), it also returns an error explaining why
// the write is short.
func (w *Writer) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

// Encode encodes v with a JSON encoder and writes the output to the
// underlying writer.
func (w *Writer) Encode(v interface{}) error {
	return w.enc.Encode(v)
}
