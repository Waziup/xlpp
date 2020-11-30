package xlpp

import (
	"errors"
	"io"
)

var errObjectKeyNoDepth = errors.New("xlpp: AddObjectKey requires AddObject first")
var errEndObjectNoDepth = errors.New("xlpp: EndObject requires AddObject first")
var errEndArrayNoDepth = errors.New("xlpp: EndArray requires AddArray first")

// Writer wrapps an [io.Writer](https://golang.org/pkg/io/#Writer) with simple LPP methods for known data types.
type Writer struct {
	io.Writer
}

// NewWriter creates a Writer that wrapps an [io.Writer](https://golang.org/pkg/io/#Writer).
func NewWriter(w io.Writer) *Writer {
	return &Writer{Writer: w}
}

// Add writes a new Value to the Writer.
func (w *Writer) Add(channel uint8, v Value) (n int, err error) {
	n, err = w.Write([]byte{byte(channel)})
	if err == nil {
		var m int
		m, err = write(w.Writer, v)
		n += m
	}
	return
}

func write(w io.Writer, v Value) (n int, err error) {
	{
		var m int
		t := v.XLPPType()
		m, err = w.Write([]byte{byte(t)})
		n += m
		if err != nil {
			return
		}
	}
	{
		var m int64
		m, err = v.WriteTo(w)
		n += int(m)
		if err != nil {
			return
		}
	}
	return
}
