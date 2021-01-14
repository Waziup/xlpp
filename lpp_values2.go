package xlpp

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// The following types are supported by this library:
const (
	TypeVoltage       Type = 116 // 2 bytes 0.01V unsigned
	TypeCurrent       Type = 117 // 2 bytes 0.001A unsigned
	TypeFrequency     Type = 118 // 4 bytes 1Hz unsigned
	TypePercentage    Type = 120 // 1 byte 1-100% unsigned
	TypeAltitude      Type = 121 // 2 byte 1m signed
	TypeConcentration Type = 125 // 2 bytes, 1 ppm unsigned
	TypePower         Type = 128 // 2 byte, 1W, unsigned
	TypeDistance      Type = 130 // 4 byte, 0.001m, unsigned
	TypeEnergy        Type = 131 // 4 byte, 0.001kWh, unsigned
	TypeDirection     Type = 132 // 2 bytes, 1deg, unsigned
	TypeUnixTime      Type = 133 // 4 bytes, unsigned
	TypeColour        Type = 135 // 1 byte per RGB Color
	TypeSwitch        Type = 142 // 1 byte, 0/1
)

////////////////////////////////////////////////////////////////////////////////

// Voltage is a floating point number electrical voltage [V] with 0.01V data resolution (unsigned).
// E.g. a value of 2.3456V is written as 2.34.
type Voltage float32

// XLPPType for Voltage returns TypeVoltage.
func (v Voltage) XLPPType() Type {
	return TypeVoltage
}

func (v Voltage) String() string {
	return fmt.Sprintf("%.2f V", v)
}

// ReadFrom reads the Voltage from the reader.
func (v *Voltage) ReadFrom(r io.Reader) (n int64, err error) {
	var b [2]byte
	n, err = readFrom(r, b[:])
	d := int16(b[0])<<8 + int16(b[1])
	*v = Voltage(d) / 100
	return
}

// WriteTo writes the Voltage to the writer.
func (v Voltage) WriteTo(w io.Writer) (n int64, err error) {
	i := int16(v * 100)
	m, err := w.Write([]byte{byte(i >> 8), byte(i)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// Current is a floating point number electrical Current [A] with 0.001A data resolution (unsigned).
// E.g. a value of 2.3456A is written as 2.345.
type Current float32

// XLPPType for Current returns TypeCurrent.
func (v Current) XLPPType() Type {
	return TypeCurrent
}

func (v Current) String() string {
	return fmt.Sprintf("%.3f A", v)
}

// ReadFrom reads the Current from the reader.
func (v *Current) ReadFrom(r io.Reader) (n int64, err error) {
	var b [2]byte
	n, err = readFrom(r, b[:])
	d := int16(b[0])<<8 + int16(b[1])
	*v = Current(d) / 1000
	return
}

// WriteTo writes the Current to the writer.
func (v Current) WriteTo(w io.Writer) (n int64, err error) {
	i := int16(v * 1000)
	m, err := w.Write([]byte{byte(i >> 8), byte(i)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// Frequency is a four byte integer value (unsigned).
type Frequency uint32

// XLPPType for Frequency returns TypeFrequency.
func (v Frequency) XLPPType() Type {
	return TypeFrequency
}

func (v Frequency) String() string {
	return fmt.Sprintf("%d Hz", v)
}

// ReadFrom reads the Frequency from the reader.
func (v *Frequency) ReadFrom(r io.Reader) (n int64, err error) {
	var b [4]byte
	n, err = readFrom(r, b[:])
	*v = Frequency((uint32(b[0]) << 24) + (uint32(b[1]) << 16) + (uint32(b[2]) << 8) + uint32(b[3]))
	return
}

// WriteTo writes the Frequency to the writer.
func (v Frequency) WriteTo(w io.Writer) (n int64, err error) {
	m, err := w.Write([]byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// Percentage is a one byte integer value (unsigned).
type Percentage int8

// XLPPType for Percentage returns TypePercentage.
func (v Percentage) XLPPType() Type {
	return TypePercentage
}

func (v Percentage) String() string {
	return fmt.Sprintf("%d", v)
}

// ReadFrom reads the Percentage from the reader.
func (v *Percentage) ReadFrom(r io.Reader) (n int64, err error) {
	var b [1]byte
	n, err = readFrom(r, b[:])
	*v = Percentage(b[0])
	return
}

// WriteTo writes the Percentage to the writer.
func (v Percentage) WriteTo(w io.Writer) (n int64, err error) {
	m, err := w.Write([]byte{byte(v)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// Altitude is a floating point with 1 m data resolution (unsigned).
// E.g. a value of 3145.82m is written as 3145.
type Altitude float32

// XLPPType for Altitude returns TypeAltitude.
func (v Altitude) XLPPType() Type {
	return TypeAltitude
}

func (v Altitude) String() string {
	return fmt.Sprintf("%.0f m", v)
}

// ReadFrom reads the Altitude from the reader.
func (v *Altitude) ReadFrom(r io.Reader) (n int64, err error) {
	var b [2]byte
	n, err = readFrom(r, b[:])
	d := int16(b[0])<<8 + int16(b[1])
	*v = Altitude(d)
	return
}

// WriteTo writes the Altitude to the writer.
func (v Altitude) WriteTo(w io.Writer) (n int64, err error) {
	i := int16(v)
	m, err := w.Write([]byte{byte(i >> 8), byte(i)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// Concentration is a two byte integer chemical concentration value [ppm] (unsigned).
type Concentration uint16

// XLPPType for Concentration returns TypeConcentration.
func (v Concentration) XLPPType() Type {
	return TypeConcentration
}

func (v Concentration) String() string {
	return fmt.Sprintf("%d ppm", v)
}

// ReadFrom reads the Concentration from the reader.
func (v *Concentration) ReadFrom(r io.Reader) (n int64, err error) {
	var b [2]byte
	n, err = readFrom(r, b[:])
	*v = Concentration(b[0])<<8 + Concentration(b[1])
	return
}

// WriteTo writes the Concentration to the writer.
func (v Concentration) WriteTo(w io.Writer) (n int64, err error) {
	m, err := w.Write([]byte{byte(v >> 8), byte(v)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// Power is a two byte integer value (unsigned).
type Power uint16

// XLPPType for Power returns TypePower.
func (v Power) XLPPType() Type {
	return TypePower
}

func (v Power) String() string {
	return fmt.Sprintf("%d W", v)
}

// ReadFrom reads the Power from the reader.
func (v *Power) ReadFrom(r io.Reader) (n int64, err error) {
	var b [2]byte
	n, err = readFrom(r, b[:])
	*v = Power(b[0])<<8 + Power(b[1])
	return
}

// WriteTo writes the Power to the writer.
func (v Power) WriteTo(w io.Writer) (n int64, err error) {
	m, err := w.Write([]byte{byte(v >> 8), byte(v)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// Distance is a 4-byte floating point number [m] with 0.001 data resolution (unsigned).
type Distance float32

// XLPPType for Distance returns TypeDistance.
func (v Distance) XLPPType() Type {
	return TypeDistance
}

func (v Distance) String() string {
	return fmt.Sprintf("%.4f m", v)
}

// ReadFrom reads the Distance from the reader.
func (v *Distance) ReadFrom(r io.Reader) (n int64, err error) {
	var b [4]byte
	n, err = readFrom(r, b[:])
	d := int32(b[0])<<24 + int32(b[0])<<16 + int32(b[0])<<8 + int32(b[0])
	*v = Distance(d) / 1000
	return
}

// WriteTo writes the Distance to the writer.
func (v Distance) WriteTo(w io.Writer) (n int64, err error) {
	i := int32(v * 1000)
	m, err := w.Write([]byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// Energy is a 4-byte floating point number [kWh] with 0.001 data resolution (unsigned).
type Energy float32

// XLPPType for Energy returns TypeEnergy.
func (v Energy) XLPPType() Type {
	return TypeEnergy
}

func (v Energy) String() string {
	return fmt.Sprintf("%.4f m", v)
}

// ReadFrom reads the Energy from the reader.
func (v *Energy) ReadFrom(r io.Reader) (n int64, err error) {
	var b [4]byte
	n, err = readFrom(r, b[:])
	d := int32(b[0])<<24 + int32(b[0])<<16 + int32(b[0])<<8 + int32(b[0])
	*v = Energy(d) / 1000
	return
}

// WriteTo writes the Energy to the writer.
func (v Energy) WriteTo(w io.Writer) (n int64, err error) {
	i := int32(v * 1000)
	m, err := w.Write([]byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// Direction is a floating point with 1 deg data resolution (unsigned).
type Direction float32

// XLPPType for Direction returns TypeDirection.
func (v Direction) XLPPType() Type {
	return TypeDirection
}

func (v Direction) String() string {
	return fmt.Sprintf("%.0f deg", v)
}

// ReadFrom reads the Direction from the reader.
func (v *Direction) ReadFrom(r io.Reader) (n int64, err error) {
	var b [2]byte
	n, err = readFrom(r, b[:])
	d := uint16(b[0])<<8 + uint16(b[1])
	*v = Direction(d)
	return
}

// WriteTo writes the Direction to the writer.
func (v Direction) WriteTo(w io.Writer) (n int64, err error) {
	i := uint16(v)
	m, err := w.Write([]byte{byte(i >> 8), byte(i)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// UnixTime is a 4-byte integer value (unsigned).
type UnixTime time.Time

// XLPPType for UnixTime returns TypeUnixTime.
func (v UnixTime) XLPPType() Type {
	return TypeUnixTime
}

// ReadFrom reads the UnixTime from the reader.
func (v *UnixTime) ReadFrom(r io.Reader) (n int64, err error) {
	var b [4]byte
	n, err = readFrom(r, b[:])
	u := uint32(b[0])<<24 + uint32(b[1])<<16 + uint32(b[2])<<8 + uint32(b[0])
	*v = UnixTime(time.Unix(int64(u), 0))
	return
}

// WriteTo writes the UnixTime to the writer.
func (v UnixTime) WriteTo(w io.Writer) (n int64, err error) {
	u := uint32(time.Time(v).Unix())
	m, err := w.Write([]byte{byte(u >> 24), byte(u >> 16), byte(u >> 8), byte(u)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// Colour is a struct of {R, G, B} integer numbers with 1 byte each.
type Colour struct {
	R, G, B uint8
}

func (v Colour) String() string {
	return fmt.Sprintf("R:%d G:%d B:%d (#%02x%02x%02x)", v.R, v.G, v.B, v.R, v.G, v.B)
}

// XLPPType for Colour returns TypeColour.
func (v Colour) XLPPType() Type {
	return TypeColour
}

// ReadFrom reads the Colour from the reader.
func (v *Colour) ReadFrom(r io.Reader) (n int64, err error) {
	var b [3]byte
	n, err = readFrom(r, b[:])
	v.R = uint8(b[0])
	v.G = uint8(b[1])
	v.B = uint8(b[2])
	return
}

// WriteTo writes the Colour to the writer.
func (v Colour) WriteTo(w io.Writer) (n int64, err error) {
	m, err := w.Write([]byte{byte(v.R), byte(v.G), byte(v.B)})
	return int64(m), err
}

func (v Colour) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("#%02x%02x%02x", v.R, v.G, v.B)
	return json.Marshal(str)
}

func (v *Colour) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	_, err := fmt.Sscanf(str, "#%02x%02x%02x", &v.R, &v.G, &v.B)
	return err
}

////////////////////////////////////////////////////////////////////////////////

// Switch is a simple ON / OFF switch.
type Switch bool

// XLPPType for GPS returns TypeGPS.
func (v Switch) XLPPType() Type {
	return TypeSwitch
}

func (v Switch) String() string {
	if v {
		return "ON"
	}
	return "OFF"
}

// ReadFrom reads the Switch from the reader.
func (v *Switch) ReadFrom(r io.Reader) (n int64, err error) {
	var b [1]byte
	n, err = readFrom(r, b[:])
	*v = b[0] != 0
	return
}

// WriteTo writes the Switch to the writer.
func (v Switch) WriteTo(w io.Writer) (n int64, err error) {
	var m int
	if v {
		m, err = w.Write([]byte{byte(1)})
	} else {
		m, err = w.Write([]byte{byte(0)})
	}
	return int64(m), err
}
