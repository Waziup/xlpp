package xlpp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

// A Reader decodes values from the underlying reader.
type Reader struct {
	r *bufio.Reader
}

// NewReader constructs a new XLPP reader to get XLPP values from a underlying reader.
func NewReader(r io.Reader) *Reader {
	if br, ok := r.(*bufio.Reader); ok {
		return &Reader{r: br}
	}
	return &Reader{
		r: bufio.NewReader(r),
	}
}

func toErr(err error) error {
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return err
}

func read(r io.Reader) (v Value, n int64, err error) {
	var t Type
	{
		// read Type byte
		var buf [1]byte
		var m int
		m, err = r.Read(buf[:])
		n += int64(m)
		if err != nil {
			err = toErr(err)
			return
		}
		t = Type(buf[0])
	}
	{
		// init zero Type
		c := Registry[t]
		if c == nil {
			err = fmt.Errorf("unregistered XLPP type 0x%02x", t)
			return
		}
		v = c()
		if v == nil {
			panic(fmt.Errorf("registered XLPP type 0x%02x returned nil value", t))
		}
	}
	{
		// read value
		var m int64
		m, err = v.ReadFrom(r)
		n += m
		if err != nil {
			err = fmt.Errorf("can not read XLPP type 0x%02x: %v", t, err)
			return
		}
	}
	return
}

// Next reads the next channel and value from the reader.
func (r *Reader) Next() (channel int, v Value, err error) {
	var c byte
	c, err = r.r.ReadByte()
	channel = int(c)
	if err != nil {
		if err == io.EOF {
			err = nil
		}
		return
	}
	switch channel {
	case ChanDelay:
		v = new(Delay)
		_, err = v.ReadFrom(r.r)
	case ChanActuators:
		v = new(Actuators)
		_, err = v.ReadFrom(r.r)
	case ChanActuatorsWithChannel:
		v = new(ActuatorsWithChannel)
		_, err = v.ReadFrom(r.r)
	default:
		v, _, err = read(r.r)
	}

	return
}

func (r *Reader) Print() error {
	log.Printf("chan | value")
	i := 0
	for {
		channel, value, err := r.Next()
		if err != nil {
			log.Printf("xlpp error: %v", err)
			return err
		}
		if value == nil {
			log.Printf("end (%d values)", i)
			return nil
		}
		i++
		log.Printf("%-4d  %+v", channel, value)
	}
}

func (r *Reader) Sprint() (string, error) {
	var s strings.Builder
	log.Printf("chan | value\n")
	i := 0
	for {
		channel, value, err := r.Next()
		if err != nil {
			s.WriteString(fmt.Sprintf("xlpp error: %v\n", err))
			return s.String(), err
		}
		if value == nil {
			s.WriteString(fmt.Sprintf("end (%d values)\n", i))
			return s.String(), err
		}
		i++
		s.WriteString(fmt.Sprintf("%-4d  %+v\n", channel, value))
	}
}
