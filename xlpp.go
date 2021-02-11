package xlpp

import "io"

// Type is the XLPP data type id.
type Type uint8

// A Value is a XLPP item with type and value.
type Value interface {
	XLPPType() Type
	String() string
	io.ReaderFrom
	io.WriterTo
}

type Marker interface {
	Value
	XLPPChannel() int
}
