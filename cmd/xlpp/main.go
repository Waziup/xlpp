package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/waziup/xlpp"
)

var registry = make(map[string]func() xlpp.Value, len(xlpp.Registry))

func main() {
	var err error
	log.SetFlags(0)

	decode := flag.Bool("d", false, "decode")
	encode := flag.Bool("e", false, "encode")
	format := flag.String("f", "", "format, json or bin")
	help := flag.Bool("h", false, "help")

	flag.Parse()

	if *help {
		log.Print("Usage:")
		log.Print(`  xlpp -e '{"temperature5":23.5}'`)
		log.Print(`  xlpp -d 'AGcA6w=='`)
		log.Print(``)
		log.Print(`JSON Format: { type channel : value, ...}`)
		log.Print("XLPP types and example zero value:")
		for _, f := range xlpp.Registry {
			if v := f(); v != nil {
				data, err := json.Marshal(v)
				if err == nil {
					log.Printf("%19s: %s", typeName(v), data)
				}
			}
		}
		return
	}

	for _, f := range xlpp.Registry {
		if v := f(); v != nil {
			registry[typeName(v)] = f
		}
	}

	var data []byte

	if *decode {
		if flag.Arg(0) != "" {
			data = []byte(flag.Arg(0))
		} else {
			data, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatal(err)
			}
		}
		switch *format {
		case "b64", "base64", "":
			data = base642xlpp(data)
		case "bin":
		default:
			log.Fatal("unknown format")
		}
		data = xlpp2json(data)
		os.Stdout.Write(data)
		return

	} else if *encode {
		if flag.Arg(0) != "" {
			data = []byte(flag.Arg(0))
		} else {
			data, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatal(err)
			}
		}
		data = json2xlpp(data)
		switch *format {
		case "bin":
		case "b64", "base64", "":
			data = xlpp2base64(data)
		default:
			log.Fatal("unknown format")
		}
		os.Stdout.Write(data)
		return
	}
}

var jsonKeyRegexp = regexp.MustCompile(`^([a-zA-Z]+)([0-9]+)$`)

func xlpp2base64(data []byte) []byte {
	str := base64.StdEncoding.EncodeToString(data)
	return []byte(str)
}

func base642xlpp(data []byte) []byte {
	data, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func json2xlpp(data []byte) []byte {
	var buf bytes.Buffer
	w := xlpp.NewWriter(&buf)

	values := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &values); err != nil {
		log.Fatal(err)
	}

	for key, m := range values {
		match := jsonKeyRegexp.FindStringSubmatch(key)
		if match == nil {
			log.Fatal("bad json entry: ", key)
		}
		name := match[1]
		channel, _ := strconv.Atoi(match[2])
		f, ok := registry[name]
		if !ok {
			log.Fatal("unknown type: ", name)
		}
		v := f()
		if err := json.Unmarshal(m, &v); err != nil {
			log.Fatalf("can not unmarshal %q: %v", name, err)
		}
		if _, err := w.Add(uint8(channel), v); err != nil {
			log.Fatalf("can not write %q: %v", name, err)
		}
	}

	return buf.Bytes()
}

func xlpp2json(data []byte) []byte {
	buf := bytes.NewBuffer(data)
	r := xlpp.NewReader(buf)
	values := make(map[string]interface{})

	for {
		channel, value, err := r.Next()
		if err != nil {
			log.Fatal("can not read xlpp: ", err)
		}
		if value == nil {
			break
		}
		name := typeName(value) + strconv.Itoa(channel)
		values[name] = value
	}
	data, err := json.Marshal(values)
	if err != nil {
		log.Fatal("can not marshal json: ", err)
	}
	return data
}

func typeName(v interface{}) (name string) {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}
	return strings.ToLower(name)
}
