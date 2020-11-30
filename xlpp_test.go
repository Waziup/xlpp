package xlpp_test

import (
	"bytes"
	"log"
	"reflect"
	"testing"

	"github.com/waziup/xlpp"
)

var digitalInput = xlpp.DigitalInput(12)
var digitalOutput = xlpp.DigitalOutput(12)
var analogInput = xlpp.AnalogInput(3.75)
var analogOutput = xlpp.AnalogOutput(4.25)
var luminosity = xlpp.Luminosity(45)
var presence = xlpp.Presence(5)
var temperature = xlpp.Temperature(31.6)
var relativeHumidity = xlpp.RelativeHumidity(22.5)
var accelerometer = xlpp.Accelerometer{X: 3.245, Y: -0.171, Z: 0.909}
var barometricPressure = xlpp.BarometricPressure(4.1)
var gyromter = xlpp.Gyrometer{X: 4.25, Y: 5.1, Z: 0.2}
var gps = xlpp.GPS{Latitude: 51.0493, Longitude: 13.7381, Meters: 122}

var null = xlpp.Null{}
var bin = xlpp.Binary([]byte{1, 2, 3, 7, 8, 9})
var integer = xlpp.Integer(5182)
var str = xlpp.String("test :)")
var boolean = xlpp.Bool(true)
var object = xlpp.Object{
	"count": &integer,
	"pos":   &gps,
	"val":   &digitalInput,
}
var array = xlpp.Array{
	&presence,
	&luminosity,
	&temperature,
}

var values = []xlpp.Value{
	// LPP types
	&digitalInput,
	&digitalOutput,
	&analogInput,
	&analogOutput,
	&luminosity,
	&presence,
	&temperature,
	&relativeHumidity,
	&accelerometer,
	&barometricPressure,
	&gyromter,
	&gps,
	// XLPP types
	&null,
	&bin,
	&integer,
	&str,
	&boolean,
	&object,
	&array,
}

func TestSimple(t *testing.T) {
	var buf bytes.Buffer

	w := xlpp.NewWriter(&buf)

	for i, value := range values {
		var channel = uint8(i)
		w.Add(channel, value)
	}

	log.Printf("buffer size: %d", buf.Len())

	r := xlpp.NewReader(&buf)
	for {
		channel, value, err := r.Next()
		if err != nil {
			log.Fatal("reading error:", err)
		}
		if value == nil {
			log.Fatal("end")
			break
		}
		log.Printf("%2d: %v", channel, value)
	}
}

func TestWriter(t *testing.T) {
	var buf bytes.Buffer
	w := xlpp.NewWriter(&buf)
	r := xlpp.NewReader(&buf)

	for i, vIn := range values {
		_, err := w.Add(uint8(i), vIn)
		if err != nil {
			t.Fatalf("can not write %T (%+v): %v", deref(vIn), deref(vIn), err)
		}
		data := buf.Bytes()
		channel, vOut, err := r.Next()
		if err != nil {
			t.Logf("data: %v", data)
			t.Fatalf("can not read %T (%+v): %v", deref(vIn), deref(vIn), err)
		}
		if reflect.TypeOf(vIn) != reflect.TypeOf(vOut) {
			t.Logf("data: %v", data)
			t.Fatalf("write <> read: %T <> %T", deref(vIn), deref(vOut))
		}
		if !reflect.DeepEqual(vIn, vOut) {
			t.Logf("data: %v", data)
			t.Fatalf("write <> read: %T (%+v) <> (%+v)", deref(vIn), deref(vIn), deref(vOut))
		}
		if channel != i {
			t.Logf("data: %v", data)
			t.Fatalf("write chan <> read chan: %T %d <> %d", deref(vIn), deref(vIn), deref(vOut))
		}
		if buf.Len() != 0 {
			t.Logf("data: %v", data)
			t.Fatalf("buffer has %d pending bytes after write: %T (%#v)", buf.Len(), deref(vIn), deref(vIn))
		}
	}
}

func deref(i interface{}) interface{} {
	v := reflect.ValueOf(i)
	if v.Type().Kind() == reflect.Ptr {
		return v.Elem().Interface()
	}
	return i
}
