// Copyright ©2019 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nmea

import (
	"errors"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	errTooShort      = errors.New("nmea: sentence is too short")
	errNoSigil       = errors.New("nmea: no initial sentence sigil")
	errChecksum      = errors.New("nmea: checksum mismatch")
	errNotPointer    = errors.New("nmea: destination not a pointer")
	errNotStruct     = errors.New("nmea: destination is not a struct")
	errType          = errors.New("nmea: wrong type for method")
	errLateType      = errors.New("nmea: late type field")
	errMissingType   = errors.New("nmea: missing type field")
	errTypeSyntax    = errors.New("nmea: bad syntax for type match")
	errNotRegistered = errors.New("nmea: sentence type not registered")
)

// ParseTo parses a raw NMEA 0183 sentence and fills the fields of dst with the
// data contained within the sentence. If the sentence has a checksum it is
// compared with the checksum of the sentence's bytes.
//
// The concrete value of dst must be a pointer to a struct.
func ParseTo(dst interface{}, sentence string) error {
	switch {
	case len(sentence) < 6: // [!$].{5}
		return errTooShort
	case sentence[0] != '$' && sentence[0] != '!':
		return errNoSigil
	}
	sentence = sentence[1:]

	var sum int64
	if sumMarkIdx := strings.Index(sentence, "*"); sumMarkIdx != -1 {
		wantSum, err := strconv.ParseInt(sentence[sumMarkIdx+1:], 16, 8)
		if err != nil {
			return err
		}
		sentence = sentence[:sumMarkIdx]
		if checksum(sentence) != wantSum {
			return errChecksum
		}
		sum = wantSum
	}

	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr {
		return errNotPointer
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errNotStruct
	}

	return parseTo(rv, strings.Split(sentence, ","), sum)
}

// Register registers the NMEA 0183 type to be parsed into the given
// destination type, dst. The kind of dst must be a struct, otherwise
// Register will panic. Calling Register with an already registered
// type will overwrite the existing registration. If dst is nil, the
// type will be deregistered.
//
// The following types are registered by default:
//
//  - "AIVDM": VDMVDO{}
//  - "AIVDO": VDMVDO{}
//  - "GLGNS": GNS{}
//  - "GNGNS": GNS{}
//  - "GPBOD": BOD{}
//  - "GPBWC": BWC{}
//  - "GPGGA": GGA{}
//  - "GPGLL": GLL{}
//  - "GPGNS": GNS{}
//  - "GPGSA": GSA{}
//  - "GPGSV": GSV{}
//  - "GPHDT": HDT{}
//  - "GPR00": R00{}
//  - "GPRMA": RMA{}
//  - "GPRMB": RMB{}
//  - "GPRMC": RMC{}
//  - "GPSTN": STN{}
//  - "GPTHS": THS{}
//  - "GPTRF": TRF{}
//  - "GPVBW": VBW{}
//  - "GPVTG": VTG{}
//  - "GPWPL": WPL{}
//  - "GPXTE": XTE{}
//  - "GPZDA": ZDA{}
//  - "PGRME": RME{}
//  - "PGRMM": RMM{}
//  - "PGRMZ": RMZ{}
//  - "PSLIB": LIB{}
//
func Register(typ string, dst interface{}) {
	if dst == nil {
		registryLock.Lock()
		delete(registry, typ)
		registryLock.Unlock()
		return
	}
	if reflect.TypeOf(dst).Kind() != reflect.Struct {
		panic(errNotStruct)
	}
	registryLock.Lock()
	registry[typ] = dst
	registryLock.Unlock()
}

var (
	registryLock sync.RWMutex
	registry     = map[string]interface{}{
		"AIVDM": VDMVDO{},
		"AIVDO": VDMVDO{},
		"GLGNS": GNS{},
		"GNGNS": GNS{},
		"GPBOD": BOD{},
		"GPBWC": BWC{},
		"GPGGA": GGA{},
		"GPGLL": GLL{},
		"GPGNS": GNS{},
		"GPGSA": GSA{},
		"GPGSV": GSV{},
		"GPHDT": HDT{},
		"GPR00": R00{},
		"GPRMA": RMA{},
		"GPRMB": RMB{},
		"GPRMC": RMC{},
		"GPSTN": STN{},
		"GPTHS": THS{},
		"GPTRF": TRF{},
		"GPVBW": VBW{},
		"GPVTG": VTG{},
		"GPWPL": WPL{},
		"GPXTE": XTE{},
		"GPZDA": ZDA{},
		"PGRME": RME{},
		"PGRMM": RMM{},
		"PGRMZ": RMZ{},
		"PSLIB": LIB{},
	}
)

// Parse parses a raw NMEA 0183 sentence and fills the fields of a destination
// registered struct with the data contained within the sentence and returns it.
// If the sentence has a checksum it is compared with the checksum of the
// sentence's bytes.
func Parse(sentence string) (interface{}, error) {
	switch {
	case len(sentence) < 6: // [!$].{5}
		return nil, errTooShort
	case sentence[0] != '$' && sentence[0] != '!':
		return nil, errNoSigil
	}
	sentence = sentence[1:]

	var sum int64
	if sumMarkIdx := strings.Index(sentence, "*"); sumMarkIdx != -1 {
		wantSum, err := strconv.ParseInt(sentence[sumMarkIdx+1:], 16, 8)
		if err != nil {
			return nil, err
		}
		sentence = sentence[:sumMarkIdx]
		if checksum(sentence) != wantSum {
			return nil, errChecksum
		}
		sum = wantSum
	}

	fields := strings.Split(sentence, ",")

	registryLock.RLock()
	dst, ok := registry[fields[0]]
	registryLock.RUnlock()
	if !ok {
		return nil, errNotRegistered
	}

	typ := reflect.TypeOf(dst)
	if typ.Kind() != reflect.Struct {
		return nil, errNotStruct
	}
	rv := reflect.New(typ).Elem()
	err := parseTo(rv, fields, sum)
	if err != nil {
		return nil, err
	}
	return rv.Interface(), nil
}

func parseTo(rv reflect.Value, fields []string, sum int64) error {
	rt := rv.Type()

	var hasType bool
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		tag := rt.Field(i).Tag.Get("nmea")
		if tag == "" {
			continue
		}

		if rt.Field(i).Name == "Type" {
			if i != 0 {
				return errLateType
			}
			if tag[0] == '/' {
				if tag[len(tag)-1] != '/' {
					return errTypeSyntax
				}
				re, err := regexp.Compile(tag[1 : len(tag)-1])
				if err != nil {
					return errTypeSyntax
				}
				if !re.MatchString(fields[i]) {
					return errType
				}
			} else if tag != fields[i] {
				return errType
			}
			hasType = true
			if f.Kind() == reflect.String {
				f.SetString(fields[i])
			}
			continue
		}

		switch tag {
		default:
			if i >= len(fields) {
				continue
			}
			err := methodFor[tag](f, fields[i])
			if err != nil {
				return err
			}
		case "checksum":
			switch f.Kind() {
			default:
				return errType
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				f.SetInt(sum)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				f.SetUint(uint64(sum))
			}
		}
	}

	if !hasType {
		return errMissingType
	}
	return nil
}

func checksum(s string) int64 {
	var sum byte
	for _, b := range []byte(s) {
		sum ^= b
	}
	return int64(sum)
}

// TODO(kortschak): Add helper method registration.
var methodFor = map[string]func(dst reflect.Value, field string) error{
	"number": setNumber,
	"string": setString,
	"latlon": setLatLon,
	"date":   setDate,
	"time":   setTime,
}

func setNumber(dst reflect.Value, field string) error {
	switch dst.Kind() {
	default:
		return errType
	case reflect.Float32, reflect.Float64:
		return setFloat(dst, field)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setInteger(dst, field)
	}
}

func setInteger(dst reflect.Value, field string) error {
	switch kind := dst.Kind(); kind {
	default:
		return errType
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if len(field) == 0 {
			dst.SetInt(0)
			break
		}
		val, err := strconv.ParseInt(field, 10, sizeOf[kind])
		if err != nil {
			return err
		}
		dst.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if len(field) == 0 {
			dst.SetUint(0)
			break
		}
		val, err := strconv.ParseUint(field, 10, sizeOf[kind])
		if err != nil {
			return err
		}
		dst.SetUint(val)
	}
	return nil
}

var sizeOf = [...]int{
	reflect.Int:    0,
	reflect.Int8:   8,
	reflect.Int16:  16,
	reflect.Int32:  32,
	reflect.Int64:  64,
	reflect.Uint:   0,
	reflect.Uint8:  8,
	reflect.Uint16: 16,
	reflect.Uint32: 32,
	reflect.Uint64: 64,
}

func setFloat(dst reflect.Value, field string) error {
	switch dst.Kind() {
	default:
		return errType
	case reflect.Float64:
		if len(field) == 0 {
			dst.SetFloat(0)
			break
		}
		val, err := strconv.ParseFloat(field, 64)
		if err != nil {
			return err
		}
		dst.SetFloat(val)
	case reflect.Float32:
		if len(field) == 0 {
			dst.SetFloat(0)
			break
		}
		val, err := strconv.ParseFloat(field, 32)
		if err != nil {
			return err
		}
		dst.SetFloat(val)
	}
	return nil
}

func setLatLon(dst reflect.Value, field string) error {
	switch dst.Kind() {
	default:
		return errType
	case reflect.Float64, reflect.Float32:
		if len(field) == 0 {
			dst.SetFloat(0)
			break
		}
		val, err := strconv.ParseFloat(field, 64)
		if err != nil {
			return err
		}
		deg, min := math.Modf(val / 100)
		dst.SetFloat(deg + min*100.0/60.0)
	}
	return nil
}

func setString(dst reflect.Value, field string) error {
	switch dst.Kind() {
	default:
		return errType
	case reflect.String:
		dst.SetString(field)
	case reflect.Slice:
		if dst.Type().Elem().Kind() != reflect.Uint8 {
			return errType
		}
		dst.SetBytes([]byte(field))
	}
	return nil
}

var timeType = reflect.TypeOf(time.Time{})

func setDate(dst reflect.Value, field string) error {
	if dst.Type() != timeType {
		return errType
	}
	t, err := time.ParseInLocation("020106", field, time.UTC)
	if err != nil {
		return err
	}
	dst.Set(reflect.ValueOf(t))
	return nil
}

func setTime(dst reflect.Value, field string) error {
	if dst.Type() != timeType {
		return errType
	}
	t, err := time.ParseInLocation("150405", field, time.UTC)
	if err != nil {
		return err
	}
	dst.Set(reflect.ValueOf(t))
	return nil
}
