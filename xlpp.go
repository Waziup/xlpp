package xlpp

import "io"

// Type is the XLPP data type id.
type Type uint8

// A Value is a XLPP item with type and value.
type Value interface {
	XLPPType() Type
	io.ReaderFrom
	io.WriterTo
}

// A ExtendedValue handles more complex XLPP types.
// type ExtendedValue interface {
// 	XLPPWriteHeadTo(w io.Writer) (n int64, err error)
// 	XLPPReadHeadFrom(r io.Reader) (n int64, err error)
// }
