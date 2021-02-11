package xlpp

import (
	"fmt"
	"io"
	"math"
	"strconv"
)

// The following types are supported by this library:
const (
	TypeDigitalInput       Type = 0   // 1 byte
	TypeDigitalOutput      Type = 1   // 1 byte
	TypeAnalogInput        Type = 2   // 2 bytes, 0.01 signed
	TypeAnalogOutput       Type = 3   // 2 bytes, 0.01 signed
	TypeLuminosity         Type = 101 // 2 bytes, 1 lux unsigned
	TypePresence           Type = 102 // 1 byte, 1
	TypeTemperature        Type = 103 // 2 bytes, 0.1°C signed
	TypeRelativeHumidity   Type = 104 // 1 byte, 0.5% unsigned
	TypeAccelerometer      Type = 113 // 2 bytes per axis, 0.001G
	TypeBarometricPressure Type = 115 // 2 bytes 0.1 hPa Unsigned
	TypeGyrometer          Type = 134 // 2 bytes per axis, 0.01 °/s
	TypeGPS                Type = 136 // 3 byte lon/lat 0.0001 °, 3 bytes alt 0.01m
)

////////////////////////////////////////////////////////////////////////////////

// DigitalInput is a one byte integer value (unsigned).
type DigitalInput uint8

// XLPPType for DigitalInput returns TypeDigitalInput.
func (v DigitalInput) XLPPType() Type {
	return TypeDigitalInput
}

func (v DigitalInput) String() string {
	return strconv.Itoa(int(v))
}

// ReadFrom reads the DigitalInput from the reader.
func (v *DigitalInput) ReadFrom(r io.Reader) (n int64, err error) {
	var b [1]byte
	n, err = readFrom(r, b[:])
	*v = DigitalInput(b[0])
	return
}

// WriteTo writes the DigitalInput to the writer.
func (v DigitalInput) WriteTo(w io.Writer) (n int64, err error) {
	m, err := w.Write([]byte{byte(v)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// DigitalOutput is a one byte integer value (unsigned).
type DigitalOutput uint8

// XLPPType for DigitalOutput returns TypeDigitalOutput.
func (v DigitalOutput) XLPPType() Type {
	return TypeDigitalOutput
}

func (v DigitalOutput) String() string {
	return strconv.Itoa(int(v))
}

// ReadFrom reads the DigitalOutput from the reader.
func (v *DigitalOutput) ReadFrom(r io.Reader) (n int64, err error) {
	var b [1]byte
	n, err = readFrom(r, b[:])
	*v = DigitalOutput(b[0])
	return
}

// WriteTo writes the DigitalOutput to the writer.
func (v DigitalOutput) WriteTo(w io.Writer) (n int64, err error) {
	m, err := w.Write([]byte{byte(v)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// AnalogInput is a floating point number with 0.01 data resolution (signed).
type AnalogInput float32

// XLPPType for AnalogInput returns TypeAnalogInput.
func (v AnalogInput) XLPPType() Type {
	return TypeAnalogInput
}

func (v AnalogInput) String() string {
	return fmt.Sprintf("%.2f", v)
}

// ReadFrom reads the AnalogInput from the reader.
func (v *AnalogInput) ReadFrom(r io.Reader) (n int64, err error) {
	var b [2]byte
	n, err = readFrom(r, b[:])
	d := int16(b[0])<<8 + int16(b[1])
	*v = AnalogInput(d) / 100
	return
}

// WriteTo writes the AnalogInput to the writer.
func (v AnalogInput) WriteTo(w io.Writer) (n int64, err error) {
	i := int16(v * 100)
	m, err := w.Write([]byte{byte(i >> 8), byte(i)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// AnalogOutput is a floating point number with 0.01 data resolution (signed).
type AnalogOutput float32

// XLPPType for AnalogOutput returns TypeAnalogOutput.
func (v AnalogOutput) XLPPType() Type {
	return TypeAnalogOutput
}

func (v AnalogOutput) String() string {
	return fmt.Sprintf("%.2f", v)
}

// ReadFrom reads the AnalogOutput from the reader.
func (v *AnalogOutput) ReadFrom(r io.Reader) (n int64, err error) {
	var b [2]byte
	n, err = readFrom(r, b[:])
	d := int16(b[0])<<8 + int16(b[1])
	*v = AnalogOutput(d) / 100
	return
}

// WriteTo writes the AnalogOutput to the writer.
func (v AnalogOutput) WriteTo(w io.Writer) (n int64, err error) {
	d := int16(v * 100)
	m, err := w.Write([]byte{byte(d >> 8), byte(d)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// Luminosity is a two byte luminusity value [lux].
type Luminosity uint16

// XLPPType for Luminosity returns TypeLuminosity.
func (v Luminosity) XLPPType() Type {
	return TypeLuminosity
}

func (v Luminosity) String() string {
	return fmt.Sprintf("%d lux", v)
}

// ReadFrom reads the Luminosity from the reader.
func (v *Luminosity) ReadFrom(r io.Reader) (n int64, err error) {
	var b [2]byte
	n, err = readFrom(r, b[:])
	*v = Luminosity(b[0])<<8 + Luminosity(b[1])
	return
}

// WriteTo writes the Luminosity to the writer.
func (v Luminosity) WriteTo(w io.Writer) (n int64, err error) {
	m, err := w.Write([]byte{byte(v >> 8), byte(v)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// Presence is a one byte integer value (unsigned).
type Presence uint8

// XLPPType for Presence returns TypePresence.
func (v Presence) XLPPType() Type {
	return TypePresence
}

func (v Presence) String() string {
	if v == 0 {
		return "no"
	}
	return "yes"
}

// ReadFrom reads the Presence from the reader.
func (v *Presence) ReadFrom(r io.Reader) (n int64, err error) {
	var b [1]byte
	n, err = readFrom(r, b[:])
	*v = Presence(b[0])
	return
}

// WriteTo writes the Presence to the writer.
func (v Presence) WriteTo(w io.Writer) (n int64, err error) {
	m, err := w.Write([]byte{byte(v)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// Temperature is a floating point number temperature [°C] with 0.1 data resolution (signed).
// E.g. a value of 27.3456°C is written as 27.3.
type Temperature float32

// XLPPType for Temperature returns TypeTemperature.
func (v Temperature) XLPPType() Type {
	return TypeTemperature
}

func (v Temperature) String() string {
	return fmt.Sprintf("%.2f °C", v)
}

// ReadFrom reads the Temperature from the reader.
func (v *Temperature) ReadFrom(r io.Reader) (n int64, err error) {
	var b [2]byte
	n, err = readFrom(r, b[:])
	d := int16(b[0])<<8 + int16(b[1])
	*v = Temperature(d) / 10
	return
}

// WriteTo writes the Temperature to the writer.
func (v Temperature) WriteTo(w io.Writer) (n int64, err error) {
	i := int16(v * 10)
	m, err := w.Write([]byte{byte(i >> 8), byte(i)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// RelativeHumidity is a floating point number humidity [%] with 0.5 data resolution (unsigned).
// E.g. a value of 12.64% is written as 12.5, or 51.434% is written as 51.0.
type RelativeHumidity float32

// XLPPType for RelativeHumidity returns TypeRelativeHumidity.
func (v RelativeHumidity) XLPPType() Type {
	return TypeRelativeHumidity
}

func (v RelativeHumidity) String() string {
	return fmt.Sprintf("%.1f %%", v)
}

// ReadFrom reads the RelativeHumidity from the reader.
func (v *RelativeHumidity) ReadFrom(r io.Reader) (n int64, err error) {
	var b [1]byte
	n, err = readFrom(r, b[:])
	*v = RelativeHumidity(b[0]) / 2
	return
}

// WriteTo writes the RelativeHumidity to the writer.
func (v RelativeHumidity) WriteTo(w io.Writer) (n int64, err error) {
	m, err := w.Write([]byte{byte(v * 2)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// Accelerometer is a struct of {x, y, z} floating point numbers [G] with 0.001 data resolution (signed) per axis.
type Accelerometer struct {
	X, Y, Z float32
}

func (v Accelerometer) String() string {
	return fmt.Sprintf("X: %.3f G, Y: %.3f G, Z: %.3f G", v.X, v.Y, v.Z)
}

// XLPPType for Accelerometer returns TypeAccelerometer.
func (v Accelerometer) XLPPType() Type {
	return TypeAccelerometer
}

// ReadFrom reads the Accelerometer from the reader.
func (v *Accelerometer) ReadFrom(r io.Reader) (n int64, err error) {
	var b [6]byte
	n, err = readFrom(r, b[:])
	vx := int16(b[0])<<8 + int16(b[1])
	vy := int16(b[2])<<8 + int16(b[3])
	vz := int16(b[4])<<8 + int16(b[5])
	v.X = float32(vx) / 1000
	v.Y = float32(vy) / 1000
	v.Z = float32(vz) / 1000
	return
}

// WriteTo writes the Accelerometer to the writer.
func (v Accelerometer) WriteTo(w io.Writer) (n int64, err error) {
	vx := int16(v.X * 1000)
	vy := int16(v.Y * 1000)
	vz := int16(v.Z * 1000)
	m, err := w.Write([]byte{byte(vx >> 8), byte(vx), byte(vy >> 8), byte(vy), byte(vz >> 8), byte(vz)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// BarometricPressure is a floating point number barometric pressure value [hPa] with 0.1 data resolution (unsigned).
type BarometricPressure float32

func (v BarometricPressure) String() string {
	return fmt.Sprintf("%.1f hPa", v)
}

// XLPPType for BarometricPressure returns TypeBarometricPressure.
func (v BarometricPressure) XLPPType() Type {
	return TypeBarometricPressure
}

// ReadFrom reads the BarometricPressure from the reader.
func (v *BarometricPressure) ReadFrom(r io.Reader) (n int64, err error) {
	var b [2]byte
	n, err = readFrom(r, b[:])
	d := int16(b[0])<<8 + int16(b[1])
	*v = BarometricPressure(d) / 10
	return
}

// WriteTo writes the BarometricPressure to the writer.
func (v BarometricPressure) WriteTo(w io.Writer) (n int64, err error) {
	i := int16(v * 10)
	m, err := w.Write([]byte{byte(i >> 8), byte(i)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// Gyrometer is a struct of {x, y, z} floating point numbers [°/s] with 0.01 data resolution (signed) per axis.
type Gyrometer struct {
	X, Y, Z float32
}

func (v Gyrometer) String() string {
	return fmt.Sprintf("X: %.3f °/s, Y: %.3f °/s, Z: %.3f °/s", v.X, v.Y, v.Z)
}

// XLPPType for Gyrometer returns TypeGyrometer.
func (v Gyrometer) XLPPType() Type {
	return TypeGyrometer
}

// ReadFrom reads the Gyrometer from the reader.
func (v *Gyrometer) ReadFrom(r io.Reader) (n int64, err error) {
	var b [6]byte
	n, err = readFrom(r, b[:])
	vx := int16(b[0])<<8 + int16(b[1])
	vy := int16(b[2])<<8 + int16(b[3])
	vz := int16(b[4])<<8 + int16(b[5])
	v.X = float32(vx) / 100
	v.Y = float32(vy) / 100
	v.Z = float32(vz) / 100
	return
}

// WriteTo writes the Gyrometer to the writer.
func (v Gyrometer) WriteTo(w io.Writer) (n int64, err error) {
	vx := int16(v.X * 100)
	vy := int16(v.Y * 100)
	vz := int16(v.Z * 100)
	m, err := w.Write([]byte{byte(vx >> 8), byte(vx), byte(vy >> 8), byte(vy), byte(vz >> 8), byte(vz)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

// GPS is a {latitude [°], longitude [°], altitude [m]} GPS location with 0.0001 data resolution (signed) for latitude and longitude, and 0.01 data resolution (signed) for altitude.
type GPS struct {
	Latitude, Longitude, Meters float32
}

// XLPPType for GPS returns TypeGPS.
func (v GPS) XLPPType() Type {
	return TypeGPS
}

func (v GPS) String() string {

	return fmt.Sprintf("%s, %s, %.2fm", dms(v.Latitude, "N", "S"), dms(v.Longitude, "E", "W"), v.Meters)
}

func dms(f float32, n string, s string) string {
	abs := abs32(f)
	deg := floor32(abs)
	_min := (abs - deg) * 60
	min := floor32(_min)
	sec := (_min - min) * 60
	dir := n
	if f < 0 {
		dir = s
	}
	return fmt.Sprintf("%.0f°%.0f'%.2f\"%s", deg, min, sec, dir)
}

func floor32(f float32) float32 {
	return float32(math.Floor(float64(f)))
}

func abs32(f float32) float32 {
	return float32(math.Abs(float64(f)))
}

// ReadFrom reads the GPS from the reader.
func (v *GPS) ReadFrom(r io.Reader) (n int64, err error) {
	var b [9]byte
	n, err = readFrom(r, b[:])
	lat := int32(b[0])<<16 + int32(b[1])<<8 + int32(b[2])
	lon := int32(b[3])<<16 + int32(b[4])<<8 + int32(b[5])
	alt := int32(b[6])<<16 + int32(b[7])<<8 + int32(b[8])
	v.Latitude = float32(lat) / 10000
	v.Longitude = float32(lon) / 10000
	v.Meters = float32(alt) / 100
	return
}

// WriteTo writes the GPS to the writer.
func (v GPS) WriteTo(w io.Writer) (n int64, err error) {
	lat := int32(v.Latitude * 10000)
	lon := int32(v.Longitude * 10000)
	alt := int32(v.Meters * 100)
	m, err := w.Write([]byte{byte(lat >> 16), byte(lat >> 8), byte(lat), byte(lon >> 16), byte(lon >> 8), byte(lon), byte(alt >> 16), byte(alt >> 8), byte(alt)})
	return int64(m), err
}

////////////////////////////////////////////////////////////////////////////////

func readFrom(r io.Reader, b []byte) (n int64, err error) {
	var m int
	m, err = io.ReadFull(r, b[:])
	n += int64(m)
	return
}
