# Extended Low Power Payload (XLPP)

This go library allows you to encode and decode light-weight payloads.
It is an extended version of [Cayenne Low Power Payload](https://www.thethingsnetwork.org/docs/devices/arduino/api/cayennelpp.html).
It can be used in Chirpstack LoRaWAN servers.

## Example

```go

package main

import (
	"log"

	"github.com/waziup/xlpp"
)

// LPP Types
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

// XLPP Types
var integer = xlpp.Integer(5182)
var str = xlpp.String("test :)")
var boolean = xlpp.Bool(true)
var bin = xlpp.Binary([]byte{1, 2, 3, 7, 8, 9})
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
var null = xlpp.Null{}


var values = []xlpp.Value{
	// LPP types
	&digitalInput, &digitalOutput, &analogInput, &analogOutput,
	&luminosity, &presence, &temperature, &relativeHumidity,
	&accelerometer, &barometricPressure, &gyromter, &gps,
	// XLPP types
	&integer, &str, &boolean, &bin, &object, &array, &null,
}


func main() {
	var buf bytes.Buffer

	// write types using xlpp.Writer
	w := xlpp.NewWriter(&buf)
	for i, value := range values {
		var channel = uint8(i)
		w.Add(channel, value)
    }
    
    log.Printf("buffer size: %d B", buf.Len())
    // > buffer size 130 B

	// read types using xlpp.Reader
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

```


# LPP Types
Those types are inherited from Cayenne LPP (https://developers.mydevices.com/cayenne/docs/lora/#lora-cayenne-low-power-payload)

Type | IPSO | LPP | Hex | Data Size | Data Resolution per bit
-- | -- | -- | -- | -- | --
Digital Input | 3200 | 0 | 0 | 1 | 1
Digital Output | 3201 | 1 | 1 | 1 | 1
Analog Input | 3202 | 2 | 2 | 2 | 0.01 Signed
Analog Output | 3203 | 3 | 3 | 2 | 0.01 Signed
Illuminance Sensor | 3301 | 101 | 65 | 2 | 1 Lux Unsigned MSB
Presence Sensor | 3302 | 102 | 66 | 1 | 1
Temperature Sensor | 3303 | 103 | 67 | 2 | 0.1 °C Signed MSB
Humidity Sensor | 3304 | 104 | 68 | 1 | 0.5 % Unsigned
Accelerometer | 3313 | 113 | 71 | 6 | 0.001 G Signed MSB per axis
Barometer | 3315 | 115 | 73 | 2 | 0.1 hPa Unsigned MSB
Gyrometer | 3334 | 134 | 86 | 6 | 0.01 °/s Signed MSB per axis
GPS Location | 3336 | 136 | 88 | 9 | Latitude : 0.0001 ° Signed MSB Longitude : 0.0001 ° Signed MSB Altitude : 0.01 meter Signed MSB

# XLPP
Additionnal types with physical dimension:

Type | LPP | Data Size | Data Resolution per bit
-- | -- | -- | --
Voltage | 116 | 2 | 0.01V Unsigned
Current | 117 | 2 | 0.001A Unsigned
Frequency | 118 | 4 | 1Hz Unsigned
Percentage | 120 | 1 | 1-100% Unsigned
Altitude | 121 | 2 | 1m Signed
Concentration | 125 | 2 | 1 ppm Unsigned
Power | 128 | 2 | 1W Unsigned
Distance | 130 | 4 | 0.001m Unsigned
Energy | 131 | 4 | 0.001kWh Unsigned
Direction | 132 | 2 | 1deg Unsigned
UnixTime | 133 | 4 | Unsigned
Colour | 135 | 1 | RGB Color
Switch | 142 | 1 | 0/1 (OFF/ON)

Additionnal types without physical dimension:

Type | XLPP | Data Size | Data Resolution per bit
-- | -- | -- | --
Integer | 51 | variant | 1
String | 52 | len(string)+1 | null terminated C string
Object | 123 | len(keys)+values+1 | keys are null terminated C strings, followed by the values
Array | 91 | len(values)+1 | list of values
Bool | 54 (true), 55 (false) | 0 | true of false
Null | 58 | 0 | (no value)
Binary | 57 | len+1 | raw binary data

# XLPP Marker Types

Markers break the normal flow of types to add more information to the stream. Markers are no sensor values and do not follow the [channel, type, data] structure!
They are identified by a reserved (fixed) channel byte, followed by theire respective content.

## Delay Marker

A Delay Marker is used to put values in a historical context. It uses the reserved channel 253 always.

Marker (Channel) | Data Size | Usage
-- | -- | --
253 | 3 | Time duration that puts all following values in a historical context.

```
 XLPP Message
------------------------------------------------------------------
 | Channel 1 |     | Channel 2 | Channel 253  | Channel 1 |     |
 | Sensor  1 | ... | Sensor  2 | Delay Marker | Sensor 1  | ... |
 | Value   1 |     | Value   2 | e.g. 1h30m   | Value 1   |     |
----------------------------------------------------------------- -
```

In the above example, a Delay Marker with `1h30m` (1 hour, 30 minutes) was added to the data, meaning that all subsequent values have been measured at this time in the past. Sensor 1 and Channel 1 can be used twice in the message, as the second occurrence reflects a measurement from the past.

A message can container multiple Delay Markers. The delays will be accumulated to a total delay. 

## Actuator Marker

An Actuator Marker is used to declare the existance of actuators to the receiver. This holds no value or state for the actuator, but the XLPP Type that this actuator consumes.

The Marker uses the reserved channels 252 (Actuators with Channel) and 251 (Actuators without channel).

Marker (Channel) | Data Size | Usage
-- | -- | --
252 | 1 + 1 x num actuators | A list of actuators (Type) that the sender of this message can consume.
251 | 1 + 2 x num actuator | A list of actuators (Channel+Type) that the sender of this message can consume.



# Binary format

```
XLPP =
	Field
	Field, Field
Field = 
	Value
	Marker
Value =
	Channel, Type, Data
	# using ony free channels
Data =
	# depends on Type and Channel
Channel =
	0   .. 249 # free channels
	250 .. 255 # reserved channels
Marker =
	Channel, Data 
```


# XLPP Binary

Install the xlpp binary from source using the [go programming language](https://golang.org/dl/) or [download a prebuild binary file](https://github.com/Waziup/xlpp/tree/main/bin).

```cmd
go install github.com/waziup/xlpp/cmd/xlpp
```

## Usage:

```bash
# Encoding: JSON -> XLPP Base64:
xlpp -e '{"temperature0":23.5}'
# AGcA6w==

# Decoding: XLPP Base64 -> JSON
xlpp -d AGcA6w==
# {"temperature0":23.5}

# Encoding Binary
xlpp -e -f bin '{"string1":"hello:)"}' > pl1.xlpp
xlpp -e -f bin '{"temperature0":23.5}' > pl1.xlpp
# Decoding Binary
xlpp -d -f bin < pl1.xlpp
# {"string1":"hello:)","temperature0":23.5}
```

## Commandline flags:

Flag | Help
-- | --
-d | decode from XLPP
-e | encode to XLPP
-f | format: `base64` (default) or `bin`


## Windows:

Keep in mind Windows CMD and Windows Powershell might need `"` double quotes escaped like this:

```bat
xlpp -e "{"""colour0""":"""#ffaa00"""}"
xlpp -e "{"""analogoutput0""":1,"""switch1""":true}"
```


# References

Based on [Cayenne Low Power Payload](https://www.thethingsnetwork.org/docs/devices/arduino/api/cayennelpp.html). See [developers.mydevices.com/cayenne](https://developers.mydevices.com/cayenne/docs/lora/#lora-cayenne-low-power-payload)
