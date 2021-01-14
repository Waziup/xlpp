package xlpp

import (
	"encoding/binary"
	"fmt"
	"io"
	"sort"
)

// The following types are supported by this library:
const (
	// extended LPP types
	TypeInteger     Type = 51
	TypeString      Type = 52
	TypeBool        Type = 53
	TypeBoolTrue    Type = 54
	TypeBoolFalse   Type = 55
	TypeObject      Type = 123 // '{'
	TypeEndOfObject Type = 0   // '}'
	TypeArray       Type = 91  // '['
	// TypeArrayOf     Type = 92  // '['
	TypeEndOfArray Type = 93 // '['
	TypeFlags      Type = 56
	TypeBinary     Type = 57
	TypeNull       Type = 58
)

// Null is a empty type. It holds no data.
type Null struct{}

// XLPPType for Null returns TypeNull
func (v Null) XLPPType() Type {
	return TypeNull
}

// ReadFrom reads the Null from the reader.
func (v *Null) ReadFrom(r io.Reader) (n int64, err error) {
	return 0, nil
}

// WriteTo writes the Null to the writer.
func (v Null) WriteTo(w io.Writer) (n int64, err error) {
	return 0, nil
}

////////////////////////////////////////////////////////////////////////////////

// Binary is a simple array of bytes.
type Binary []byte

// XLPPType for Binary returns TypeBinary.
func (v Binary) XLPPType() Type {
	return TypeBinary
}

func (v Binary) String() string {
	return fmt.Sprintf("%X", []byte(v))
}

// ReadFrom reads the Binary from the reader.
func (v *Binary) ReadFrom(r io.Reader) (n int64, err error) {
	var brc byteReaderCounter
	brc.ByteReader = newByteReader(r)
	l, err := binary.ReadUvarint(&brc)
	if err != nil {
		return int64(brc.Count), err
	}
	*v = make(Binary, l)
	var m int
	m, err = io.ReadFull(r, *v)
	return int64(brc.Count + m), err
}

// WriteTo writes the Binary to the writer.
func (v Binary) WriteTo(w io.Writer) (n int64, err error) {
	var buf [9]byte
	var m int
	m = binary.PutUvarint(buf[:], uint64(len(v)))
	n += int64(m)
	m, err = w.Write(buf[:m])
	n += int64(m)
	if err == nil {
		m, err = w.Write(v)
		n += int64(m)
	}
	return
}

////////////////////////////////////////////////////////////////////////////////

// Bool is a boolean true/false.
type Bool bool

// XLPPType for Bool returns TypeBool.
func (v Bool) XLPPType() Type {
	if v {
		return TypeBoolTrue
	}
	return TypeBoolFalse
}

// ReadFrom reads the Bool from the reader.
func (v *Bool) ReadFrom(r io.Reader) (n int64, err error) {
	return 0, nil
}

// WriteTo writes the Bool to the writer.
func (v Bool) WriteTo(w io.Writer) (n int64, err error) {
	return 0, nil
}

////////////////////////////////////////////////////////////////////////////////

// Integer is a simple integer value.
type Integer int

// XLPPType for Integer returns TypeInteger.
func (v Integer) XLPPType() Type {
	return TypeInteger
}

// ReadFrom reads the Integer from the reader.
func (v *Integer) ReadFrom(r io.Reader) (n int64, err error) {
	var brc byteReaderCounter
	brc.ByteReader = newByteReader(r)
	i64, err := binary.ReadVarint(&brc)
	*v = Integer(i64)
	return int64(brc.Count), err
}

// WriteTo writes the Integer to the writer.
func (v Integer) WriteTo(w io.Writer) (n int64, err error) {
	var buf [9]byte
	m := binary.PutVarint(buf[:], int64(v))
	m, err = w.Write(buf[:m])
	n = int64(m)
	return
}

////////////////////////////////////////////////////////////////////////////////

// String is a simple string value.
type String string

// XLPPType for String returns TypeString.
func (v String) XLPPType() Type {
	return TypeString
}

// ReadFrom reads the String from the reader.
func (v *String) ReadFrom(r io.Reader) (n int64, err error) {
	buf := make([]byte, 0, 32)
	var brc byteReaderCounter
	brc.ByteReader = newByteReader(r)
	for {
		b, err := brc.ReadByte()
		if err != nil {
			return int64(brc.Count), err
		}
		if b == 0 {
			*v = String(buf)
			return int64(brc.Count), nil
		}
		buf = append(buf, b)
	}
}

// WriteTo writes the String to the writer.
func (v String) WriteTo(w io.Writer) (n int64, err error) {
	var m int
	m, err = w.Write([]byte(v))
	n += int64(m)
	if err == nil {
		m, err = w.Write([]byte{0})
		n += int64(m)
	}
	return
}

////////////////////////////////////////////////////////////////////////////////

// Object is a simple key-value map.
type Object map[string]Value

// XLPPType for Object returns TypeObject.
func (v Object) XLPPType() Type {
	return TypeObject
}

func (v Object) keys() []string {
	keys := make([]string, len(v))
	i := 0
	for key := range v {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// func (v Object) XLPPWriteHeadTo(w io.Writer) (n int64, err error) {
// 	keys := v.keys()
// 	var m int
// 	m, err = w.Write([]byte{byte(len(v))})
// 	n += int64(m)
// 	if err != nil {
// 		return
// 	}
// 	for _, key := range keys {
// 		var m int64
// 		m, err = String(key).WriteTo(w)
// 		n += m
// 		if err != nil {
// 			return
// 		}
// 	}
// 	return
// }

// func (v *Object) XLPPReadHeadFrom(r io.Reader) (n int64, err error) {
// 	var buf [1]byte
// 	var m int
// 	m, err = r.Read(buf[:])
// 	n += int64(m)
// 	if err != nil {
// 		return
// 	}
// 	l := int(buf[0])
// 	*v = make(Object, l)
// 	for i := 0; i < l; i++ {
// 		var str String
// 		var m int64
// 		m, err = str.ReadFrom(r)
// 		n += m
// 		if err != nil {
// 			return
// 		}
// 		(*v)[string(str)] = nil
// 	}
// 	return
// }

// ReadFrom reads the Object from the reader.
func (v *Object) ReadFrom(r io.Reader) (n int64, err error) {
	*v = make(Object)

	buf := make([]byte, 32)
	var brc byteReaderCounter
	brc.ByteReader = newByteReader(r)

	for {
		var key string
		{
			var b byte
			b, err = brc.ReadByte()
			if b == byte(TypeEndOfObject) {
				return
			}
			buf = buf[:0]
			for {
				if err != nil {
					return int64(brc.Count), err
				}
				if b == 0 {
					key = string(buf)
					break
				}
				buf = append(buf, b)
				b, err = brc.ReadByte()
			}
		}
		{
			var m int64
			(*v)[key], m, err = read(r)
			n += m
			if err != nil {
				return
			}
		}
	}
}

// WriteTo writes the Object to the writer.
func (v Object) WriteTo(w io.Writer) (n int64, err error) {
	keys := v.keys()
	for _, key := range keys {
		{

			var m int64
			m, err = String(key).WriteTo(w)
			n += m
			if err != nil {
				return
			}
		}
		{
			var m int
			m, err = write(w, v[key])
			n += int64(m)
			if err != nil {
				return
			}
		}
	}
	{
		var m int
		m, err = w.Write([]byte{byte(TypeEndOfObject)})
		n += int64(m)
		if err != nil {
			return
		}
	}
	return
}

////////////////////////////////////////////////////////////////////////////////

// Array is a simple list of values.
type Array []Value

// XLPPType for Array returns TypeArray.
func (v Array) XLPPType() Type {
	// if t := v.getItemType(); t != 0 {
	// 	return TypeArrayOf
	// }
	return TypeArray
}

// func (v Array) getItemType() (t Type) {
// 	if len(v) == 0 {
// 		return 0
// 	}
// 	for i, value := range v {
// 		if i == 0 {
// 			t = value.XLPPType()
// 		} else {
// 			if t != value.XLPPType() {
// 				return 0
// 			}
// 		}
// 	}
// 	return
// }

// ReadFrom reads the Array from the reader.
func (v *Array) ReadFrom(r io.Reader) (n int64, err error) {
	*v = make(Array, 0, 8)
	for {
		var m int64
		var i Value
		i, m, err = read(r)
		n += m
		if err != nil {
			return
		}
		if _, ok := i.(endOfArray); ok {
			return
		}
		*v = append(*v, i)
	}
}

// WriteTo writes the Array to the writer.
func (v Array) WriteTo(w io.Writer) (n int64, err error) {
	{
		for _, value := range v {
			var m int
			m, err = write(w, value)
			n += int64(m)
			if err != nil {
				return
			}
		}
	}
	{
		var m int
		m, err = w.Write([]byte{byte(TypeEndOfArray)})
		n += int64(m)
		if err != nil {
			return
		}
	}
	return
}

type endOfArray struct{}

func (endOfArray) XLPPType() Type {
	return TypeEndOfArray
}

func (endOfArray) ReadFrom(r io.Reader) (n int64, err error) {
	return 0, nil
}

func (endOfArray) WriteTo(w io.Writer) (n int64, err error) {
	return 0, nil
}

////////////////////////////////////////////////////////////////////////////////

type byteReader struct {
	io.Reader
}

func newByteReader(r io.Reader) io.ByteReader {
	if br, ok := r.(io.ByteReader); ok {
		return br
	}
	return byteReader{Reader: r}
}

func (br byteReader) ReadByte() (byte, error) {
	var buf [1]byte
	_, err := br.Reader.Read(buf[:])
	return buf[0], err
}

////////////////////

type byteReaderCounter struct {
	io.ByteReader
	Count int
}

func (br byteReaderCounter) ReadByte() (byte, error) {
	b, err := br.ByteReader.ReadByte()
	br.Count++
	return b, err
}
