// Copyright ©2019 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nmea

import (
	"errors"
	"math"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrTooShort      = errors.New("nmea: sentence is too short")
	ErrNoSigil       = errors.New("nmea: no initial sentence sigil")
	ErrChecksum      = errors.New("nmea: checksum mismatch")
	ErrNotPointer    = errors.New("nmea: destination not a pointer")
	ErrNotStruct     = errors.New("nmea: destination is not a struct")
	ErrNMEAType      = errors.New("nmea: wrong nmea type for sentence")
	ErrType          = errors.New("nmea: wrong type for method")
	ErrLateType      = errors.New("nmea: late type field")
	ErrMissingType   = errors.New("nmea: missing type field")
	ErrTypeSyntax    = errors.New("nmea: bad syntax for type match")
	ErrNotRegistered = errors.New("nmea: sentence type not registered")
	ErrBadBinary     = errors.New("nmea: invalid binary data encoding")
)

// ParseTo parses a raw NMEA 0183 sentence and fills the fields of dst with the
// data contained within the sentence. If the sentence has a checksum it is
// compared with the checksum of the sentence's bytes.
//
// The concrete value of dst must be a pointer to a struct.
func ParseTo(dst interface{}, sentence string) error {
	switch {
	case len(sentence) < 6: // [!$].{5}
		return ErrTooShort
	case sentence[0] != '$' && sentence[0] != '!':
		return ErrNoSigil
	}
	sentence = sentence[1:]

	var sum, wantSum int64
	if sumMarkIdx := strings.Index(sentence, "*"); sumMarkIdx != -1 {
		var err error
		wantSum, err = strconv.ParseInt(sentence[sumMarkIdx+1:], 16, 8)
		if err != nil {
			return err
		}
		sentence = sentence[:sumMarkIdx]
		sum = checksum(sentence)
	}

	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr {
		return ErrNotPointer
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	err := parseTo(rv, strings.Split(sentence, ","), wantSum)
	if sum != wantSum {
		return ErrChecksum
	}
	return err
}

// Register registers the NMEA 0183 type to be parsed into the given
// destination type, dst. The kind of dst must be a struct, otherwise
// Register will panic. Calling Register with an already registered
// type will overwrite the existing registration. If dst is nil, the
// type will be deregistered.
//
// The following types are registered by default:
//
//  - "AIVDM", "AIVDO": VDMVDO{}
//  - "GLBOD", "GNBOD", "GPBOD": BOD{}
//  - "GLBWC", "GNBWC", "GPBWC": BWC{}
//  - "GLGGA", "GNGGA", "GPGGA": GGA{}
//  - "GLGLL", "GNGLL", "GPGLL": GLL{}
//  - "GLGNS", "GNGNS", "GPGNS": GNS{}
//  - "GLGSA", "GNGSA", "GPGSA": GSA{}
//  - "GLGSV", "GNGSV", "GPGSV": GSV{}
//  - "GLHDT", "GNHDT", "GPHDT": HDT{}
//  - "GLR00", "GNR00", "GPR00": R00{}
//  - "GLRMA", "GNRMA", "GPRMA": RMA{}
//  - "GLRMB", "GNRMB", "GPRMB": RMB{}
//  - "GLRMC", "GNRMC", "GPRMC": RMC{}
//  - "GLSTN", "GNSTN", "GPSTN": STN{}
//  - "GLTHS", "GNTHS", "GPTHS": THS{}
//  - "GLTRF", "GNTRF", "GPTRF": TRF{}
//  - "GLVBW", "GNVBW", "GPVBW": VBW{}
//  - "GLVTG", "GNVTG", "GPVTG": VTG{}
//  - "GLWPL", "GNWPL", "GPWPL": WPL{}
//  - "GLXTE", "GNXTE", "GPXTE": XTE{}
//  - "GLZDA", "GNZDA", "GPZDA": ZDA{}
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
		panic(ErrNotStruct)
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
		"GLGNS": GNS{}, "GNGNS": GNS{}, "GPGNS": GNS{},
		"GLBOD": BOD{}, "GNBOD": BOD{}, "GPBOD": BOD{},
		"GLBWC": BWC{}, "GNBWC": BWC{}, "GPBWC": BWC{},
		"GLGGA": GGA{}, "GNGGA": GGA{}, "GPGGA": GGA{},
		"GLGLL": GLL{}, "GNGLL": GLL{}, "GPGLL": GLL{},
		"GLGSA": GSA{}, "GNGSA": GSA{}, "GPGSA": GSA{},
		"GLGSV": GSV{}, "GNGSV": GSV{}, "GPGSV": GSV{},
		"GLHDT": HDT{}, "GNHDT": HDT{}, "GPHDT": HDT{},
		"GLR00": R00{}, "GNR00": R00{}, "GPR00": R00{},
		"GLRMA": RMA{}, "GNRMA": RMA{}, "GPRMA": RMA{},
		"GLRMB": RMB{}, "GNRMB": RMB{}, "GPRMB": RMB{},
		"GLRMC": RMC{}, "GNRMC": RMC{}, "GPRMC": RMC{},
		"GLSTN": STN{}, "GNSTN": STN{}, "GPSTN": STN{},
		"GLTHS": THS{}, "GNTHS": THS{}, "GPTHS": THS{},
		"GLTRF": TRF{}, "GNTRF": TRF{}, "GPTRF": TRF{},
		"GLVBW": VBW{}, "GNVBW": VBW{}, "GPVBW": VBW{},
		"GLVTG": VTG{}, "GNVTG": VTG{}, "GPVTG": VTG{},
		"GLWPL": WPL{}, "GNWPL": WPL{}, "GPWPL": WPL{},
		"GLXTE": XTE{}, "GNXTE": XTE{}, "GPXTE": XTE{},
		"GLZDA": ZDA{}, "GNZDA": ZDA{}, "GPZDA": ZDA{},
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
		return nil, ErrTooShort
	case sentence[0] != '$' && sentence[0] != '!':
		return nil, ErrNoSigil
	}
	sentence = sentence[1:]

	var sum, wantSum int64
	if sumMarkIdx := strings.Index(sentence, "*"); sumMarkIdx != -1 {
		var err error
		wantSum, err = strconv.ParseInt(sentence[sumMarkIdx+1:], 16, 8)
		if err != nil {
			return nil, err
		}
		sentence = sentence[:sumMarkIdx]
		sum = checksum(sentence)
	}

	fields := strings.Split(sentence, ",")

	registryLock.RLock()
	dst, ok := registry[fields[0]]
	registryLock.RUnlock()
	if !ok {
		return nil, ErrNotRegistered
	}

	typ := reflect.TypeOf(dst)
	if typ.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}
	rv := reflect.New(typ).Elem()
	err := parseTo(rv, fields, wantSum)
	if sum != wantSum {
		err = ErrChecksum
	}
	return rv.Interface(), err
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
				return ErrLateType
			}
			if tag[0] == '/' {
				if tag[len(tag)-1] != '/' {
					return ErrTypeSyntax
				}
				re, err := regexp.Compile(tag[1 : len(tag)-1])
				if err != nil {
					return ErrTypeSyntax
				}
				if !re.MatchString(fields[i]) {
					f.SetString(fields[i])
					return ErrNMEAType
				}
			} else if tag != fields[i] {
				f.SetString(fields[i])
				return ErrNMEAType
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
				return ErrType
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				f.SetInt(sum)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				f.SetUint(uint64(sum))
			}
		}
	}

	if !hasType {
		return ErrMissingType
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
		return ErrType
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
		return ErrType
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
		return ErrType
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
		return ErrType
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
		return ErrType
	case reflect.String:
		dst.SetString(field)
	case reflect.Slice:
		if dst.Type().Elem().Kind() != reflect.Uint8 {
			return ErrType
		}
		dst.SetBytes([]byte(field))
	}
	return nil
}

var timeType = reflect.TypeOf(time.Time{})

func setDate(dst reflect.Value, field string) error {
	if dst.Type() != timeType {
		return ErrType
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
		return ErrType
	}
	t, err := time.ParseInLocation("150405", field, time.UTC)
	if err != nil {
		return err
	}
	dst.Set(reflect.ValueOf(t))
	return nil
}

// DeArmorAIS returns 6-bit-nibble payload data extracted from AIS
// ASCII armoring. Each byte of the returned byte slice as a single
// 6-bit value.
//
// See https://gpsd.gitlab.io/gpsd/AIVDM.html#_aivdm_aivdo_payload_armoring
func DeArmorAIS(data string) ([]byte, error) {
	if data == "" {
		return nil, nil
	}
	dst := make([]byte, len(data))
	for i, b := range []byte(data) {
		if b < '0' || 'w' < b || ('X' <= b && b <= '_') {
			return dst, ErrBadBinary
		}

		v := b - '0'
		if v > 40 { // We are in ['X', '_'].
			v -= 8
		}

		dst[i] = v
	}

	return dst, nil
}

// SixBitToASCII returns the ASCII value corresponding to an AIS Sixbit
// ASCII-encoded character. If b6 is greater than 63, SixBitASCII will
// panic.
//
// See https://gpsd.gitlab.io/gpsd/AIVDM.html#_ais_payload_data_types
func SixBitToASCII(b6 byte) byte {
	if b6 > 63 {
		panic("nmea: six bit ascii overflow")
	}
	return asciiFor[b6]
}

var asciiFor = [64]byte{
	'@', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '[', '\\', ']', '^', '_',
	' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?',
}

// AISBitField returns an 8-bit packed byte slice holding the bits
// of an AIS 6-bit nibble slice, starting from bit s and extending to
// the bit before e. The resulting byte slice will be shifted such that
// the last bit of the 6-bit nibble slice will be the lowest significance
// bit of the returned value. If s or e are outside the length of the 6-bit
// bit slice, AISBitField will panic.
func AISBitField(b6 []byte, s, e int) []byte {
	if s < 0 || e < 0 || e < s {
		panic("nmea: bitfield index out of bounds")
	}
	if s == e {
		return nil
	}
	if ew, _ := bitAddr(s); ew > len(b6) {
		panic("nmea: bitfield index out of bounds")
	}

	var bits big.Int
	for i := s; i < e; i++ {
		w, b := bitAddr(i)
		bits.SetBit(&bits, e-(i-s)-1, uint(b6[w]&(1<<b))>>b)
	}
	padding := uint((6 - ((e - s) % 6)) % 6)
	bits.Rsh(&bits, padding)
	return bits.Bytes()
}

func bitAddr(i int) (w int, b uint) {
	w = i / 6
	b = 5 - uint(i)%6
	return w, b
}
