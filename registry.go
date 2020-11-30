package xlpp

var Registry = map[Type]func() Value{
	// LPP Types
	TypeDigitalInput:       func() Value { return new(DigitalInput) },
	TypeDigitalOutput:      func() Value { return new(DigitalOutput) },
	TypeAnalogInput:        func() Value { return new(AnalogInput) },
	TypeAnalogOutput:       func() Value { return new(AnalogOutput) },
	TypeLuminosity:         func() Value { return new(Luminosity) },
	TypePresence:           func() Value { return new(Presence) },
	TypeTemperature:        func() Value { return new(Temperature) },
	TypeRelativeHumidity:   func() Value { return new(RelativeHumidity) },
	TypeAccelerometer:      func() Value { return new(Accelerometer) },
	TypeBarometricPressure: func() Value { return new(BarometricPressure) },
	TypeGyrometer:          func() Value { return new(Gyrometer) },
	TypeGPS:                func() Value { return new(GPS) },
	// XLPP Types
	TypeInteger:    func() Value { return new(Integer) },
	TypeNull:       func() Value { return new(Null) },
	TypeString:     func() Value { return new(String) },
	TypeBool:       func() Value { return new(Bool) },
	TypeObject:     func() Value { return new(Object) },
	TypeArray:      func() Value { return new(Array) },
	TypeEndOfArray: func() Value { return endOfArray{} },
	// TypeArrayOf: func() Value { return new(Array) },
	// TypeFlags: func() Value { return new(Flags) },
	TypeBinary: func() Value { return new(Binary) },
}
