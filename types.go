// Copyright ©2019 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nmea

import "time"

// http://aprs.gids.nl/nmea/#bod
type BOD struct {
	Type string `nmea:"GPBOD"`

	True        float64 `nmea:"number"`
	_           [0]byte
	Magnetic    float64 `nmea:"number"`
	_           [0]byte
	Destination string `nmea:"string"`
	Start       string `nmea:"string"`

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#bwc
type BWC struct {
	Type string `nmea:"GPBWC"`

	Timestamp  time.Time `nmea:"time"`
	Latitude   float64   `nmea:"latlon"`
	NorthSouth string    `nmea:"string"`
	Longitude  float64   `nmea:"latlon"`
	EastWest   string    `nmea:"string"`
	True       float64   `nmea:"number"`
	_          [0]byte
	Magnetic   float64 `nmea:"number"`
	_          [0]byte
	Range      float64 `nmea:"number"`
	RangeUnit  string  `nmea:"string"`
	Waypoint   string  `nmea:"string"`

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#gll
type GLL struct {
	Type string `nmea:"GPGLL"`

	Latitude   float64   `nmea:"latlon"`
	NorthSouth string    `nmea:"string"`
	Longitude  float64   `nmea:"latlon"`
	EastWest   string    `nmea:"string"`
	Timestamp  time.Time `nmea:"time"`

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#gga
type GGA struct {
	Type string `nmea:"GPGGA"`

	Timestamp  time.Time `nmea:"time"`
	Latitude   float64   `nmea:"latlon"`
	NorthSouth string    `nmea:"string"`
	Longitude  float64   `nmea:"latlon"`
	EastWest   string    `nmea:"string"`

	Quality    int `nmea:"number"`
	Satellites int `nmea:"number"`

	HDOP float64 `nmea:"number"`

	Altitude     float64 `nmea:"number"`
	AltitudeUnit string  `nmea:"string"`

	Separation     float64 `nmea:"number"`
	SeparationUnit string  `nmea:"string"`

	Age float64 `nmea:"number"`

	DiffReferenceStationID string `nmea:"string"`

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#gsa
type GSA struct {
	Type string `nmea:"GPGSA"`

	Mode string `nmea:"string"`
	Fix  int    `nmea:"number"`

	SV0  string `nmea:"string"`
	SV1  string `nmea:"string"`
	SV2  string `nmea:"string"`
	SV3  string `nmea:"string"`
	SV4  string `nmea:"string"`
	SV5  string `nmea:"string"`
	SV6  string `nmea:"string"`
	SV7  string `nmea:"string"`
	SV8  string `nmea:"string"`
	SV9  string `nmea:"string"`
	SV10 string `nmea:"string"`
	SV11 string `nmea:"string"`

	PDOP float64 `nmea:"number"`
	HDOP float64 `nmea:"number"`
	VDOP float64 `nmea:"number"`

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#gsv
type GSV struct {
	Type string `nmea:"GPGSV"`

	Messages      int `nmea:"number"`
	MessageNumber int `nmea:"number"`

	SatellitesInView int `nmea:"number"`

	Satellite0PRN int `nmea:"number"`
	Elevation0    int `nmea:"number"`
	Azimuth0      int `nmea:"number"`
	SNR0          int `nmea:"number"`

	Satellite1PRN int `nmea:"number"`
	Elevation1    int `nmea:"number"`
	Azimuth1      int `nmea:"number"`
	SNR1          int `nmea:"number"`

	Satellite2PRN int `nmea:"number"`
	Elevation2    int `nmea:"number"`
	Azimuth2      int `nmea:"number"`
	SNR2          int `nmea:"number"`

	Satellite3PRN int `nmea:"number"`
	Elevation3    int `nmea:"number"`
	Azimuth3      int `nmea:"number"`
	SNR3          int `nmea:"number"`

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#hdt
type HDT struct {
	Type string `nmea:"GPHDT"`

	Heading float64 `nmea:"number"`
	_       [0]byte

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#r00
type R00 struct {
	Type string `nmea:"GPR00"`

	WP0  string `nmea:"string"`
	WP1  string `nmea:"string"`
	WP2  string `nmea:"string"`
	WP3  string `nmea:"string"`
	WP4  string `nmea:"string"`
	WP5  string `nmea:"string"`
	WP6  string `nmea:"string"`
	WP7  string `nmea:"string"`
	WP8  string `nmea:"string"`
	WP9  string `nmea:"string"`
	WP10 string `nmea:"string"`
	WP11 string `nmea:"string"`
	WP13 string `nmea:"string"`

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#rma
type RMA struct {
	Type string `nmea:"GPRMA"`

	Status string `nmea:"string"`

	Latitude   float64 `nmea:"latlon"`
	NorthSouth string  `nmea:"string"`
	Longitude  float64 `nmea:"latlon"`
	EastWest   string  `nmea:"string"`

	_, _ [0]byte

	Speed            float64 `nmea:"number"`
	CourseOverGround int     `nmea:"number"`
	Variation        float64 `nmea:"number"`
	VarDirection     string  `nmea:"string"`

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#rmb
type RMB struct {
	Type string `nmea:"GPRMB"`

	Status string `nmea:"string"`

	CrosstrackError  float64 `nmea:"number"`
	CorrectDirection string  `nmea:"string"`

	Origin string `nmea:"string"`

	Destination string  `nmea:"string"`
	Latitude    float64 `nmea:"latlon"`
	NorthSouth  string  `nmea:"string"`
	Longitude   float64 `nmea:"latlon"`
	EastWest    string  `nmea:"string"`

	RangeToDestination   float64 `nmea:"number"`
	BearingToDestination float64 `nmea:"number"`
	ClosingVelocity      float64 `nmea:"number"`

	ArrivalStatus string `nmea:"string"`

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#rmc
type RMC struct {
	Type string `nmea:"GPRMC"`

	Time time.Time `nmea:"time"`

	Status string `nmea:"string"`

	Latitude   float64 `nmea:"latlon"`
	NorthSouth string  `nmea:"string"`
	Longitude  float64 `nmea:"latlon"`
	EastWest   string  `nmea:"string"`

	Speed float64 `nmea:"number"`
	Track float64 `nmea:"number"`

	Date time.Time `nmea:"date"`

	MagneticVariation float64 `nmea:"number"`
	VarDirection      string  `nmea:"string"`

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#rte
// TODO(kortschak): $GPRTE requires multiple field handling.
//
// Routes
//
// eg. $GPRTE,2,1,c,0,PBRCPK,PBRTO,PTELGR,PPLAND,PYAMBU,PPFAIR,PWARRN,PMORTL,PLISMR*73
//     $GPRTE,2,2,c,0,PCRESY,GRYRIE,GCORIO,GWERR,GWESTG,7FED*34
//            1 2 3 4 5 ..
//
//     1. Number of sentences in sequence
//     2. Sentence number
//     3. 'c' = Current active route, 'w' = waypoint list starts with destination waypoint
//     4. Name or number of the active route
//     5. onwards, Names of waypoints in Route

// http://aprs.gids.nl/nmea/#trf
type TRF struct {
	Type string `nmea:"GPTRF"`

	Time time.Time `nmea:"time"`
	Date time.Time `nmea:"date"`

	Latitude   float64 `nmea:"latlon"`
	NorthSouth string  `nmea:"string"`
	Longitude  float64 `nmea:"latlon"`
	EastWest   string  `nmea:"string"`

	Elevation       float64 `nmea:"number"`
	Iterations      float64 `nmea:"number"`
	DoplerIntervals float64 `nmea:"number"`
	UpdateDistance  float64 `nmea:"number"`

	Satellite string `nmea:"string"`
}

// http://aprs.gids.nl/nmea/#stn
type STN struct {
	Type string `nmea:"GPSTN"`

	Talker byte `nmea:"number"`
}

// http://aprs.gids.nl/nmea/#vbw
type VBW struct {
	Type string `nmea:"GPVBW"`

	LongitudinalWaterSpeed float64 `nmea:"number"`
	TransverseWaterSpeed   float64 `nmea:"number"`
	WaterSpeedStatus       string  `nmea:"string"`

	LongitudinalGroundSpeed float64 `nmea:"number"`
	TransverseGroundSpeed   float64 `nmea:"number"`
	GroundSpeedStatus       string  `nmea:"string"`
}

// http://aprs.gids.nl/nmea/#vtg
type VTG struct {
	Type string `nmea:"GPVTG"`

	TrackTrue     float64 `nmea:"number"`
	_             [0]byte
	TrackMagnetic float64 `nmea:"number"`
	_             [0]byte
	SpeedKnots    float64 `nmea:"number"`
	_             [0]byte
	SpeedKph      float64 `nmea:"number"`
	_             [0]byte

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#wpl
type WPL struct {
	Type string `nmea:"GPWPL"`

	Latitude   float64 `nmea:"latlon"`
	NorthSouth string  `nmea:"string"`
	Longitude  float64 `nmea:"latlon"`
	EastWest   string  `nmea:"string"`

	Waypoint string `nmea:"string"`

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#xte
type XTE struct {
	Type string `nmea:"GPXTE"`

	GeneralWarning string `nmea:"string"`
	LockFlag       string `nmea:"string"`

	CrossTrackError float64 `nmea:"number"`
	Steer           string  `nmea:"string"`
	Units           string  `nmea:"string"`

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#zda
type ZDA struct {
	Type string `nmea:"GPZDA"`

	Time            time.Time `nmea:"time"`
	Day             byte      `nmea:"number"`
	Month           byte      `nmea:"number"`
	Year            int       `nmea:"number"`
	TimeZone        int8      `nmea:"number"`
	TimeZoneMinutes int8      `nmea:"number"`
}

// http://aprs.gids.nl/nmea/#rme
type RME struct {
	Type string `nmea:"PGRME"`

	HPE   float64 `nmea:"number"`
	_     [0]byte
	VPE   float64 `nmea:"number"`
	_     [0]byte
	OSEPE float64 `nmea:"number"`
	_     [0]byte

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#rmm
type RMM struct {
	Type string `nmea:"PGRMM"`

	MapDatum string `nmea:"string"`

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#rmz
type RMZ struct {
	Type string `nmea:"PGRMZ"`

	Altitude              float64 `nmea:"number"`
	_                     [0]byte
	PositionFixDimensions int8 `nmea:"number"`

	Checksum byte `nmea:"checksum"`
}

// http://aprs.gids.nl/nmea/#lib
type LIB struct {
	Type string `nmea:"PSLIB"`

	Frequency   float64 `nmea:"number"`
	BitRate     float64 `nmea:"number"`
	RequestType string  `nmea:"string"`

	Checksum byte `nmea:"checksum"`
}

// https://www.trimble.com/oem_receiverhelp/v4.44/en/NMEA-0183messages_GNS.html
type GNS struct {
	Type string `nmea:"/G[NPL]GNS/"`

	Timestamp  time.Time `nmea:"time"`
	Latitude   float64   `nmea:"latlon"`
	NorthSouth string    `nmea:"string"`
	Longitude  float64   `nmea:"latlon"`
	EastWest   string    `nmea:"string"`

	Mode       string `nmea:"string"`
	Satellites int    `nmea:"number"`

	HDOP       float64 `nmea:"number"`
	Altitude   float64 `nmea:"number"`
	Separation float64 `nmea:"number"`

	Age float64 `nmea:"number"`

	ReferenceStation uint16 `nmea:"number"`

	Checksum byte `nmea:"checksum"`
}

// http://www.nuovamarea.net/pytheas_9.html
type THS struct {
	Type string `nmea:"/..THS/"`

	Heading float64 `nmea:"number"`
	Status  string  `nmea:"string"`

	Checksum byte `nmea:"checksum"`
}

// https://gpsd.gitlab.io/gpsd/AIVDM.html
type VDMVDO struct {
	Type string `nmea:"/..VD[MO]/"`

	Fragments      int `nmea:"number"`
	FragmentNumber int `nmea:"number"`

	MessageID string `nmea:"string"`

	ChannelCode string `nmea:"string"`

	// Data is AIS ASCII-armored.
	Data    string `nmea:"string"`
	Padding byte   `nmea:"number"`

	Checksum byte `nmea:"checksum"`
}
