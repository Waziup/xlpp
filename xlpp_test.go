package xlpp_test

import (
	"bytes"
	"log"
	"reflect"
	"testing"
	"time"

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
var gyromter = xlpp.Gyrometer{X: 4.25, Y: 5.10, Z: 0.21}
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

var exampleTime, _ = time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 2 15:04:05 -0700 MST 2006")

var voltage = xlpp.Voltage(1.45)
var current = xlpp.Current(4.41)
var frequency = xlpp.Frequency(8100)
var percentage = xlpp.Percentage(17)
var altitude = xlpp.Altitude(8849)
var concentration = xlpp.Concentration(2512)
var power = xlpp.Power(1142)
var distance = xlpp.Distance(2.411)
var energy = xlpp.Energy(2.876)
var direcion = xlpp.Direction(90)
var unixtime = xlpp.UnixTime(exampleTime.Round(0))
var color = xlpp.Colour{R: 123, G: 54, B: 89}
var swithc = xlpp.Switch(true)

var delay = xlpp.Delay(time.Second * 4235)
var actuators = xlpp.Actuators{xlpp.TypeColour, xlpp.TypeAnalogOutput, xlpp.TypeSwitch}
var actuatorsWithChannel = xlpp.ActuatorsWithChannel{
	xlpp.Actuator{
		Channel: 3,
		Type:    xlpp.TypeVoltage,
	},
	xlpp.Actuator{
		Channel: 17,
		Type:    xlpp.TypeColour,
	},
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
	// more LPP types
	&voltage,
	&current,
	&frequency,
	&percentage,
	&altitude,
	&concentration,
	&power,
	&distance,
	&energy,
	&direcion,
	&unixtime,
	&color,
	&swithc,
	// XLPP types
	&null,
	&bin,
	&integer,
	&str,
	&boolean,
	&object,
	&array,
	// special XLPP types
	&delay,
	&actuators,
	&actuatorsWithChannel,
}

func TestSimple(t *testing.T) {
	var buf bytes.Buffer

	w := xlpp.NewWriter(&buf)

	for i, value := range values {
		w.Add(i, value)
	}

	log.Printf("buffer size: %d", buf.Len())

	r := xlpp.NewReader(&buf)
	for {
		channel, value, err := r.Next()
		if err != nil {
			log.Fatal("reading error:", err)
		}
		if value == nil {
			log.Print("end")
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
		chanIn := i
		if marker, ok := vIn.(xlpp.Marker); ok {
			chanIn = marker.XLPPChannel()
		}
		_, err := w.Add(chanIn, vIn)
		if err != nil {
			t.Fatalf("can not write %T (%+v): %v", deref(vIn), deref(vIn), err)
		}
		data := buf.Bytes()
		chanOut, vOut, err := r.Next()
		if err != nil {
			t.Logf("data: %v", data)
			t.Fatalf("can not read %T (%+v): %v", deref(vIn), deref(vIn), err)
		}
		if reflect.TypeOf(vIn) != reflect.TypeOf(vOut) {
			t.Logf("data: %v", data)
			t.Fatalf("write <> read: %T <> %T", deref(vIn), deref(vOut))
		}
		if tIn, ok := vIn.(*xlpp.UnixTime); ok {
			tOut := vOut.(*xlpp.UnixTime)
			if !time.Time(*tIn).Equal(time.Time(*tOut)) {
				t.Logf("data: %v", data)
				t.Fatalf("write <> read: %T (%+v) <> (%+v)", deref(vIn), deref(vIn), deref(vOut))
			}
		} else {
			if !reflect.DeepEqual(vIn, vOut) {
				t.Logf("data: %v", data)
				t.Fatalf("write <> read: %T (%+v) <> (%+v)", deref(vIn), deref(vIn), deref(vOut))
			}
		}
		if chanIn != chanOut {
			t.Logf("data: %v", data)
			t.Fatalf("write chan <> read chan: %T %d <> %d", deref(vIn), chanIn, chanOut)
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
