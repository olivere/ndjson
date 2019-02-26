package ndjson

import (
	"bufio"
	"encoding/json"
	"io"
)

var (
	_ io.Reader = (*Reader)(nil)
)

// Reader allows reading line-oriented JSON data following the
// ndjson spec at http://ndjson.org/.
type Reader struct {
	r   io.Reader
	s   *bufio.Scanner
	row interface{}
}

// NewReader returns a new reader, using the underlying io.Reader
// as input.
func NewReader(r io.Reader) *Reader {
	s := bufio.NewScanner(r)
	return &Reader{r: r, s: s}
}

// Read reads data into p. It returns the number of bytes read into p.
// The bytes are taken from the underlying reader. Read follows the
// protocol defined by io.Reader.
func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	return n, err
}

// Next advances the Reader to the next line, which will then be available
// through the Decode method. It returns false when the reader stops,
// either by reaching the end of the input or an error. After Next returns
// false, the Err method will return any error that occured while reading,
// except if it was io.EOF, Err will return nil.
//
// Next might panic if the underlying split function returns too many tokens
// without advancing the input.
func (r *Reader) Next() bool {
	return r.s.Scan()
}

// Err returns the first non-EOF error that was encountered by the Reader.
func (r *Reader) Err() error {
	return r.s.Err()
}

// Decode decodes the bytes read after the last call to Next into
// the specified value.
func (r *Reader) Decode(v interface{}) error {
	return json.Unmarshal(r.s.Bytes(), v)
}
