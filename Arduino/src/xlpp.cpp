#include "xlpp.h"

#pragma region Head

XLPP::XLPP(uint8_t size) : size(size)
{
    buf = (uint8_t *)malloc(size);
    o = 0;
}

XLPP::~XLPP(void)
{
    free(buf);
}

////////////////////

#define int24_t int32_t
#define uint24_t uint32_t

#define _WRITE(v) buf[o++] = (v)&0xff;
#define WRITE_uint8_t(v) _WRITE(v);
#define WRITE_uint16_t(v) \
    _WRITE(v >> 8);       \
    _WRITE(v);
#define WRITE_uint24_t(v) \
    _WRITE(v >> 16);      \
    _WRITE(v >> 8);       \
    _WRITE(v);
#define WRITE_uint32_t(v) \
    _WRITE(v >> 24);      \
    _WRITE(v >> 16);      \
    _WRITE(v >> 8);       \
    _WRITE(v);

#define WRITE_int8_t(v) WRITE_uint8_t(uint8_t(v));
#define WRITE_int16_t(v) WRITE_uint16_t(uint16_t(v));
#define WRITE_int24_t(v) WRITE_uint24_t(uint24_t(v));
#define WRITE_int32_t(v) WRITE_uint32_t(uint32_t(v));

#define _READ buf[o++]
#define READ_uint8_t _READ
#define READ_uint16_t uint16_t(_READ) << 8 + uint16_t(_READ)
#define READ_uint24_t uint24_t(_READ) << 16 + uint24_t(_READ) << 8 + uint24_t(_READ)
#define READ_uint32_t uint32_t(_READ) << 24 + uint32_t(_READ) << 16 + uint32_t(_READ) << 8 + uint32_t(_READ)

#define READ_int8_t int8_t(READ_uint8_t)
#define READ_int16_t int16_t(READ_uint16_t)
#define READ_int24_t int24_t(READ_uint24_t)
#define READ_int32_t int32_t(READ_uint32_t)

#define XLPP_(NAME, TYPE, VALUE_T, MULTI, WIRE_T)        \
    void XLPP::add##NAME(uint8_t channel, VALUE_T value) \
    {                                                    \
        buf[o++] = channel;                              \
        add##NAME(value);                                \
    }                                                    \
    void XLPP::add##NAME(VALUE_T value)                  \
    {                                                    \
        buf[o++] = TYPE;                                 \
        WIRE_T v = value * MULTI;                        \
        WRITE_##WIRE_T(v);                               \
    }                                                    \
    VALUE_T XLPP::get##NAME()                            \
    {                                                    \
        return VALUE_T(READ_##WIRE_T) / MULTI;           \
    }

////////////////////

void XLPP::reset(void)
{
    o = 0;
}

uint8_t XLPP::getSize(void)
{
    return o;
}

uint8_t *XLPP::getBuffer(void)
{
    return buf;
}

//

uint8_t XLPP::getChannel()
{
    return buf[o++];
}

uint8_t XLPP::getType()
{
    return buf[o++];
}

#pragma endregion

////////////////////////////////////////////////////////////////////////////////

XLPP_(DigitalInput, LPP_DIGITAL_INPUT, uint8_t, 1, uint8_t);
XLPP_(DigitalOutput, LPP_DIGITAL_INPUT, uint8_t, 1, uint8_t);
XLPP_(AnalogInput, LPP_ANALOG_INPUT, float, 100, int16_t);
XLPP_(AnalogOutput, LPP_ANALOG_OUTPUT, float, 100, int16_t);
XLPP_(Luminosity, LPP_LUMINOSITY, uint16_t, 1, uint16_t);
XLPP_(Presence, LPP_PRESENCE, uint8_t, 1, uint8_t);
XLPP_(Temperature, LPP_TEMPERATURE, float, 10, int16_t);
XLPP_(RelativeHumidity, LPP_RELATIVE_HUMIDITY, float, 2, int8_t);

//

void XLPP::addAccelerometer(uint8_t channel, float x, float y, float z)
{
    buf[o++] = channel;
    addAccelerometer(x, y, z);
}

void XLPP::addAccelerometer(float x, float y, float z)
{
    buf[o++] = LPP_ACCELEROMETER;
    int16_t vx = int16_t(x * 1000);
    WRITE_int16_t(vx);
    int16_t vy = int16_t(y * 1000);
    WRITE_int16_t(vy);
    int16_t vz = int16_t(z * 1000);
    WRITE_int16_t(vz);
}

Accelerometer XLPP::getAccelerometer()
{
    Accelerometer a;
    a.x = float(READ_int16_t) / 1000;
    a.y = float(READ_int16_t) / 1000;
    a.z = float(READ_int16_t) / 1000;
    return a;
}

//

XLPP_(BarometricPressure, LPP_BAROMETRIC_PRESSURE, float, 10, int16_t);

//

void XLPP::addGyrometer(uint8_t channel, float x, float y, float z)
{
    buf[o++] = channel;
    addGyrometer(x, y, z);
}

void XLPP::addGyrometer(float x, float y, float z)
{
    buf[o++] = LPP_GYROMETER;
    int16_t vx = int16_t(x * 100);
    WRITE_int16_t(vx);
    int16_t vy = int16_t(y * 100);
    WRITE_int16_t(vy);
    int16_t vz = int16_t(z * 100);
    WRITE_int16_t(vz);
}

Gyrometer XLPP::getGyrometer()
{
    Gyrometer v;
    v.x = float(READ_int16_t) / 100;
    v.y = float(READ_int16_t) / 100;
    v.z = float(READ_int16_t) / 100;
    return v;
}

//

void XLPP::addGPS(uint8_t channel, float latitude, float longitude, float altitude)
{
    buf[o++] = channel;
    addGPS(latitude, longitude, altitude);
}

void XLPP::addGPS(float latitude, float longitude, float altitude)
{
    buf[o++] = LPP_GPS;
    int32_t lat = int32_t(latitude * 10000);
    WRITE_int24_t(lat);
    int32_t lon = int32_t(longitude * 10000);
    WRITE_int24_t(lon);
    int32_t alt = int32_t(altitude * 100);
    WRITE_int24_t(alt);
}

GPS XLPP::getGPS()
{
    GPS v;
    v.latitude = float(READ_int24_t) / 10000;
    v.longitude = float(READ_int24_t) / 10000;
    v.altitude = float(READ_int24_t) / 100;
    return v;
}

////////////////////

XLPP_(Voltage, LPP_VOLTAGE, float, 100, uint16_t);
XLPP_(Current, LPP_CURRENT, float, 1000, uint16_t);
XLPP_(Frequency, LPP_FREQUENCY, uint32_t, 1, uint32_t);
XLPP_(Percentage, LPP_PERCENTAGE, uint8_t, 1, uint8_t);
XLPP_(Altitude, LPP_ALTITUDE, float, 1, uint16_t);
XLPP_(Power, LPP_POWER, uint16_t, 1, uint16_t);
XLPP_(Distance, LPP_DISTANCE, float, 1000, uint32_t);
XLPP_(Energy, LPP_ENERGY, float, 1000, uint32_t);
XLPP_(UnixTime, LPP_UNIXTIME, uint32_t, 1, uint32_t);
XLPP_(Direction, LPP_DIRECTION, float, 1, uint16_t);
XLPP_(Switch, LPP_SWITCH, uint8_t, 1, uint8_t);
XLPP_(Concentration, LPP_CONCENTRATION, uint16_t, 1, uint16_t);

//

void XLPP::addColour(uint8_t channel, uint8_t r, uint8_t g, uint8_t b)
{
    buf[o++] = channel;
    addColour(r, g, b);
}

void XLPP::addColour(uint8_t r, uint8_t g, uint8_t b)
{
    buf[o++] = LPP_COLOUR;
    WRITE_uint8_t(r);
    WRITE_uint8_t(g);
    WRITE_uint8_t(b);
}

Colour XLPP::getColour()
{
    Colour c;
    c.r = READ_uint8_t;
    c.g = READ_uint8_t;
    c.b = READ_uint8_t;
    return c;
}

////////////////////////////////////////////////////////////////////////////////

void XLPP::addInteger(uint8_t channel, int64_t i)
{
    buf[o++] = channel;
    addInteger(i);
}

void XLPP::addInteger(int64_t i)
{
    buf[o++] = XLPP_INTEGER;
    uint64_t ui = uint64_t(i) << 1;
    if (i<0) ui = ~ui;
    while (ui >= 0x80) {
        WRITE_uint8_t(uint8_t(ui) | 0x80);
        ui >>= 7;
    }
    WRITE_uint8_t(uint8_t(ui));
}

int64_t XLPP::getInteger()
{
    uint64_t ui = 0;
    uint8_t s = 0;
    for (int i=0;; i++) {
        if(i==10) {
            return 0; // overflow
        }
        uint8_t b = READ_uint8_t;
        if (b<0x80) {
            if (i==9 && b>1) {
                return 0; // overflow
            }
            ui |= uint64_t(b)<<s;
            break;
        }
        ui |= uint64_t(b&0x7f)<<s;
        s += 7;
    }

    int64_t i = int64_t(ui >> 1);
    if (ui&1) {
        i = ~i;
    }
    return i;
}

//

void XLPP::addString(uint8_t channel, const char* str)
{
    buf[o++] = channel;
    addString(str);
}

void XLPP::addString(const char* str)
{
    buf[o++] = XLPP_STRING;
    strcpy((char*) buf+o, str);
    o += strlen(str);
}

void XLPP::getString(char* str)
{
    strcpy(str, (const char*) buf+o);
    o += strlen(str);
}

size_t XLPP::getString(char* str, size_t limit)
{
    size_t n=0;
    for (; *buf != 0 && n<limit; n++)
        str[n] = buf[o++];
    if (n==limit) {
        while (buf[o++]); // skip remaining bytes
    } else {
        str[n] = 0;
        o++;
    }
    return n;
}

//

void XLPP::addBool(uint8_t channel, bool b)
{
    buf[o++] = channel;
    addBool(b);
}

void XLPP::addBool(bool b)
{
    buf[o++] = b ? XLPP_BOOL_TRUE : XLPP_BOOL_FALSE;
}

bool XLPP::getBool()
{
    return READ_uint8_t != 0;
}

//

void XLPP::beginObject(uint8_t channel)
{
    buf[o++] = channel;
    beginObject();
}

void XLPP::beginObject()
{
    buf[o++] = XLPP_OBJECT;
}

void XLPP::addObjectKey(const char* key)
{
    strcpy((char*) buf+o, key);
    o += strlen(key);
}

void XLPP::endObject()
{
    buf[o++] = XLPP_END_OF_OBJECT;
}

//

void XLPP::beginArray(uint8_t channel)
{
    buf[o++] = channel;
    beginArray();
}

void XLPP::beginArray()
{
    buf[o++] = XLPP_ARRAY;
}

void XLPP::endArray()
{
    buf[o++] = XLPP_END_OF_ARRAY;
}

//

void XLPP::addBinary(uint8_t channel, const void* data, size_t s)
{
    buf[o++] = channel;
    addBinary(data, s);
}

void XLPP::addBinary(const void* data, size_t s)
{
    buf[o++] = XLPP_BINARY;
    size_t i = s;
    while (s >= 0x80) {
        WRITE_uint8_t(uint8_t(s) | 0x80);
        s >>= 7;
    }
    WRITE_uint8_t(uint8_t(s));
    memcpy(buf+o, data, s);
    o += s;
}

size_t XLPP::getBinary(void* data)
{
    uint64_t s = 0;
    uint8_t d = 0;
    for (int i=0;; i++) {
        if(i==10) {
            return 0; // overflow
        }
        uint8_t b = READ_uint8_t;
        if (b<0x80) {
            if (i==9 && b>1) {
                return 0; // overflow
            }
            s |= uint64_t(b)<<s;
            break;
        }
        s |= uint64_t(b&0x7f)<<s;
        s += 7;
    }

    memcpy(data, buf+o, s);
    o += s;
}

//

void XLPP::addNull(uint8_t channel)
{
    buf[o++] = channel;
    addNull();
}

void XLPP::addNull()
{
    buf[o++] = XLPP_NULL;
}

void XLPP::getNull()
{
    // NOP
}