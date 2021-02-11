package xlpp

import (
	"encoding/binary"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
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

// Special (reserved) channels for "Marker" types:
const (
	// special XLPP channels
	ChanDelay                = 253
	ChanActuators            = 252
	ChanActuatorsWithChannel = 251
)

// Null is a empty type. It holds no data.
type Null struct{}

// XLPPType for Null returns TypeNull
func (v Null) XLPPType() Type {
	return TypeNull
}

func (v Null) String() string {
	return "null"
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

func (v Bool) String() string {
	if v {
		return "true"
	}
	return "false"
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

func (v Integer) String() string {
	return fmt.Sprintf("%d", int(v))
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

func (v String) String() string {
	return fmt.Sprintf("%q", string(v))
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

func (v Object) String() string {
	var b strings.Builder
	b.WriteByte('{')
	first := true
	for key, value := range v {
		if !first {
			b.WriteByte(',')
		}
		first = false
		b.WriteString(key)
		b.WriteByte(':')
		b.WriteByte(' ')
		b.WriteString(value.String())
	}
	b.WriteByte('}')
	return b.String()
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

func (v Array) String() string {
	var b strings.Builder
	b.WriteByte('[')
	first := true
	for _, t := range v {
		if !first {
			b.WriteByte(',')
			b.WriteByte(' ')
		}
		first = false
		b.WriteString(t.String())
	}
	b.WriteByte(']')
	return b.String()
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

func (endOfArray) String() string {
	return "<end of array>"
}

func (endOfArray) ReadFrom(r io.Reader) (n int64, err error) {
	return 0, nil
}

func (endOfArray) WriteTo(w io.Writer) (n int64, err error) {
	return 0, nil
}

////////////////////////////////////////////////////////////////////////////////

// A Delay is not a Value, but a marker in XLPP data that puts values in a historical context.
// All subsequent values have been measured at this Delay in the past.
// You can use multiple Delays in one XLPP message, in which they will increment the total Delay.
type Delay time.Duration

// XLPPType for Delay returns 255.
func (v Delay) XLPPType() Type {
	return 255
}

// XLPPChannel for Delay returns the constant ChanDelay 253.
func (v Delay) XLPPChannel() int {
	return ChanDelay
}

func (v Delay) String() string {
	return time.Duration(v).String()
}

func (v Delay) Hours() int {
	return int(time.Duration(v).Hours())
}

func (v Delay) Minutes() int {
	return int(time.Duration(v).Minutes()) % 60
}

func (v Delay) Seconds() int {
	return int(time.Duration(v).Seconds()) % 60
}

// ReadFrom reads the Delay from the reader.
func (v *Delay) ReadFrom(r io.Reader) (n int64, err error) {
	var b [3]byte
	n, err = readFrom(r, b[:])
	*v = Delay(time.Hour)*Delay(b[0]) + Delay(time.Minute)*Delay(b[1]) + Delay(time.Second)*Delay(b[2])
	return
}

// WriteTo writes the Delay to the writer.
func (v Delay) WriteTo(w io.Writer) (n int64, err error) {
	m, err := w.Write([]byte{byte(v.Hours()), byte(v.Minutes()), byte(v.Seconds())})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

type Actuators []Type

// XLPPType for Actuators returns 255.
func (v Actuators) XLPPType() Type {
	return 255
}

// XLPPChannel for Actuators returns the constant ChanActuators 252.
func (v Actuators) XLPPChannel() int {
	return ChanActuators
}

func (v Actuators) String() string {
	var b strings.Builder
	b.WriteByte('[')
	first := true
	for _, t := range v {
		if !first {
			b.WriteByte(',')
			b.WriteByte(' ')
		}
		first = false
		fmt.Fprintf(&b, "0x%02X", int(t))
	}
	b.WriteByte(']')
	return b.String()
}

// ReadFrom reads the Actuators from the reader.
func (v *Actuators) ReadFrom(r io.Reader) (n int64, err error) {
	var b [1]byte
	n, err = readFrom(r, b[:])
	if err != nil {
		return
	}
	var m int64
	l := int(b[0])
	*v = make(Actuators, l)
	for i := 0; i < l; i++ {
		m, err = readFrom(r, b[:])
		if err != nil {
			return
		}
		n += m
		(*v)[i] = Type(b[0])
	}
	return
}

// WriteTo writes the Actuators to the writer.
func (v Actuators) WriteTo(w io.Writer) (n int64, err error) {
	d := make([]byte, len(v)+1)
	d[0] = byte(len(v))
	for i, a := range v {
		d[i+1] = byte(a)
	}
	m, err := w.Write(d)
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

type Actuator struct {
	Channel int
	Type    Type
}

type ActuatorsWithChannel []Actuator

// XLPPType for ActuatorsWithChannel returns 255.
func (v ActuatorsWithChannel) XLPPType() Type {
	return 255
}

// ActuatorsWithChannel for Actuators returns the constant ChanActuators 251.
func (v ActuatorsWithChannel) XLPPChannel() int {
	return ChanActuatorsWithChannel
}

func (v ActuatorsWithChannel) String() string {
	var b strings.Builder
	b.WriteByte('[')
	first := true
	for _, t := range v {
		if !first {
			b.WriteByte(',')
			b.WriteByte(' ')
		}
		first = false
		fmt.Fprintf(&b, "Chan %d: 0x%02X", int(t.Channel), int(t.Type))
	}
	b.WriteByte(']')
	return b.String()
}

// ReadFrom reads the ActuatorsWithChannel from the reader.
func (v *ActuatorsWithChannel) ReadFrom(r io.Reader) (n int64, err error) {
	var b [2]byte
	n, err = readFrom(r, b[:1])
	if err != nil {
		return
	}
	var m int64
	l := int(b[0])
	*v = make(ActuatorsWithChannel, l)
	for i := 0; i < l; i++ {
		m, err = readFrom(r, b[:])
		if err != nil {
			return
		}
		n += m
		(*v)[i] = Actuator{
			Channel: int(b[0]),
			Type:    Type(b[1]),
		}
	}
	return
}

// WriteTo writes the Actuators to the writer.
func (v ActuatorsWithChannel) WriteTo(w io.Writer) (n int64, err error) {
	d := make([]byte, len(v)*2+1)
	d[0] = byte(len(v))
	for i, a := range v {
		d[i*2+1] = byte(a.Channel)
		d[i*2+2] = byte(a.Type)
	}
	m, err := w.Write(d)
	return int64(m), err
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
