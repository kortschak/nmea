// Copyright ©2019 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nmea

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

var parseTests = []struct {
	sentence string
	dst      interface{}
	want     interface{}
}{
	{
		sentence: "$GPBOD,099.3,T,105.6,M,POINTB,*48",
		dst:      &BOD{},
		want: &BOD{
			Type:        "GPBOD",
			True:        99.3,
			Magnetic:    105.6,
			Destination: "POINTB",
			Checksum:    0x48,
		},
	},
	{
		sentence: "$GPBOD,097.0,T,103.2,M,POINTB,POINTA*4A",
		dst:      &BOD{},
		want: &BOD{
			Type:        "GPBOD",
			True:        97,
			Magnetic:    103.2,
			Destination: "POINTB",
			Start:       "POINTA",
			Checksum:    0x4a,
		},
	},
	{
		sentence: "$GPBWC,081837,,,,,,T,,M,,N,*13",
		dst:      &BWC{},
		want: &BWC{
			Type:      "GPBWC",
			Timestamp: time.Date(0, 1, 1, 8, 18, 37, 0, time.UTC),
			RangeUnit: "N",
			Checksum:  0x13,
		},
	},
	{
		sentence: "$GPBWC,225444,4917.24,N,12309.57,W,051.9,T,031.6,M,001.3,N,004*29",
		dst:      &BWC{},
		want: &BWC{
			Type:      "GPBWC",
			Timestamp: time.Date(0, 1, 1, 22, 54, 44, 0, time.UTC),
			Latitude:  49.28733333333333, NorthSouth: "N",
			Longitude: 123.1595, EastWest: "W",
			True:     51.9,
			Magnetic: 31.6,
			Range:    1.3, RangeUnit: "N",
			Waypoint: "004",
			Checksum: 0x29,
		},
	},
	{
		sentence: "$GPBWC,220516,5130.02,N,00046.34,W,213.8,T,218.0,M,0004.6,N,EGLM*21",
		dst:      &BWC{},
		want: &BWC{
			Type:      "GPBWC",
			Timestamp: time.Date(0, 1, 1, 22, 05, 16, 0, time.UTC),
			Latitude:  51.50033333333334, NorthSouth: "N",
			Longitude: 0.7723333333333334, EastWest: "W",
			True:     213.8,
			Magnetic: 218,
			Range:    4.6, RangeUnit: "N",
			Waypoint: "EGLM",
			Checksum: 0x21,
		},
	},
	{
		sentence: "$GPGGA,123456,3455.083,S,13836.285,E,1,2,3,4,M,5,M,,*4A",
		dst:      &GGA{},
		want: &GGA{
			Type:      "GPGGA",
			Timestamp: time.Date(0, 1, 1, 12, 34, 56, 0, time.UTC),
			Latitude:  34.918049999999994, NorthSouth: "S",
			Longitude: 138.60475000000002, EastWest: "E",
			Quality:    1,
			Satellites: 2,
			HDOP:       3,
			Altitude:   4, AltitudeUnit: "M",
			Separation: 5, SeparationUnit: "M",
			Age:                    0,
			DiffReferenceStationID: "",
			Checksum:               0x4a,
		},
	},
	{
		sentence: "$GPGGA,123519,4807.038,N,01131.000,W,1,2,3,4,M,5,M,,*41",
		dst:      &GGA{},
		want: &GGA{
			Type:      "GPGGA",
			Timestamp: time.Date(0, 1, 1, 12, 35, 19, 0, time.UTC),
			Latitude:  48.117299999999986, NorthSouth: "N",
			Longitude: 11.516666666666667, EastWest: "W",
			Quality:    1,
			Satellites: 2,
			HDOP:       3,
			Altitude:   4, AltitudeUnit: "M",
			Separation: 5, SeparationUnit: "M",
			Age:                    0,
			DiffReferenceStationID: "",
			Checksum:               0x41,
		},
	},
	{
		sentence: "$GPGGA,170834,4124.8963,N,08151.6838,W,1,05,1.5,280.2,M,-34.0,M,,*75",
		dst:      &GGA{},
		want: &GGA{
			Type:      "GPGGA",
			Timestamp: time.Date(0, 1, 1, 17, 8, 34, 0, time.UTC),
			Latitude:  41.41493833333334, NorthSouth: "N",
			Longitude: 81.86139666666665, EastWest: "W",
			Quality:    1,
			Satellites: 5,
			HDOP:       1.5,
			Altitude:   280.2, AltitudeUnit: "M",
			Separation: -34, SeparationUnit: "M",
			Age:                    0,
			DiffReferenceStationID: "",
			Checksum:               0x75,
		},
	},
	{
		sentence: "$GPGLL,5300.97914,N,00259.98174,E,125926,A*28",
		dst:      &GLL{},
		want: &GLL{
			Type:     "GPGLL",
			Latitude: 53.01631900000001, NorthSouth: "N",
			Longitude: 2.9996956666666668, EastWest: "E",
			Timestamp: time.Date(0, 1, 1, 12, 59, 26, 0, time.UTC),
			Checksum:  0x28,
		},
	},
	{
		sentence: "$GPGLL,3751.65,S,14507.36,E*77",
		dst:      &GLL{},
		want: &GLL{
			Type:     "GPGLL",
			Latitude: 37.86083333333333, NorthSouth: "S",
			Longitude: 145.12266666666667, EastWest: "E",
			Checksum: 0x77,
		},
	},
	{
		sentence: "$GPGLL,4916.45,N,12311.12,W,225444,A",
		dst:      &GLL{},
		want: &GLL{
			Type:     "GPGLL",
			Latitude: 49.27416666666666, NorthSouth: "N",
			Longitude: 123.18533333333335, EastWest: "W",
			Timestamp: time.Date(0, 1, 1, 22, 54, 44, 0, time.UTC),
		},
	},
	{
		sentence: "$GPGSA,A,3,,,,,,16,18,,22,24,,,3.6,2.1,2.2*3C",
		dst:      &GSA{},
		want: &GSA{
			Type:     "GPGSA",
			Mode:     "A",
			Fix:      3,
			SV5:      "16",
			SV6:      "18",
			SV8:      "22",
			SV9:      "24",
			PDOP:     3.6,
			HDOP:     2.1,
			VDOP:     2.2,
			Checksum: 0x3c,
		},
	},
	{
		sentence: "$GPGSA,A,3,19,28,14,18,27,22,31,39,,,,,1.7,1.0,1.3*34",
		dst:      &GSA{},
		want: &GSA{
			Type:     "GPGSA",
			Mode:     "A",
			Fix:      3,
			SV0:      "19",
			SV1:      "28",
			SV2:      "14",
			SV3:      "18",
			SV4:      "27",
			SV5:      "22",
			SV6:      "31",
			SV7:      "39",
			PDOP:     1.7,
			HDOP:     1,
			VDOP:     1.3,
			Checksum: 0x34,
		},
	},
	{
		sentence: "$GPGSV,3,1,11,03,03,111,00,04,15,270,00,06,01,010,00,13,06,292,00*74",
		dst:      &GSV{},
		want: &GSV{
			Type:             "GPGSV",
			Messages:         3,
			MessageNumber:    1,
			SatellitesInView: 11,
			Satellite0PRN:    3,
			Elevation0:       3,
			Azimuth0:         111,
			SNR0:             0,
			Satellite1PRN:    4,
			Elevation1:       15,
			Azimuth1:         270,
			SNR1:             0,
			Satellite2PRN:    6,
			Elevation2:       1,
			Azimuth2:         10,
			SNR2:             0,
			Satellite3PRN:    13,
			Elevation3:       6,
			Azimuth3:         292,
			SNR3:             0,
			Checksum:         0x74,
		},
	},
	{
		sentence: "$GPGSV,3,2,11,14,25,170,00,16,57,208,39,18,67,296,40,19,40,246,00*74",
		dst:      &GSV{},
		want: &GSV{
			Type:             "GPGSV",
			Messages:         3,
			MessageNumber:    2,
			SatellitesInView: 11,
			Satellite0PRN:    14,
			Elevation0:       25,
			Azimuth0:         170,
			SNR0:             0,
			Satellite1PRN:    16,
			Elevation1:       57,
			Azimuth1:         208,
			SNR1:             39,
			Satellite2PRN:    18,
			Elevation2:       67,
			Azimuth2:         296,
			SNR2:             40,
			Satellite3PRN:    19,
			Elevation3:       40,
			Azimuth3:         246,
			SNR3:             0,
			Checksum:         0x74,
		},
	},
	{
		sentence: "$GPGSV,3,3,11,22,42,067,42,24,14,311,43,27,05,244,00,,,,*4D",
		dst:      &GSV{},
		want: &GSV{
			Type:             "GPGSV",
			Messages:         3,
			MessageNumber:    3,
			SatellitesInView: 11,
			Satellite0PRN:    22,
			Elevation0:       42,
			Azimuth0:         67,
			SNR0:             42,
			Satellite1PRN:    24,
			Elevation1:       14,
			Azimuth1:         311,
			SNR1:             43,
			Satellite2PRN:    27,
			Elevation2:       5,
			Azimuth2:         244,
			SNR2:             0,
			Checksum:         0x4d,
		},
	},
	{
		sentence: "$GPGSV,1,1,13,02,02,213,,03,-3,000,,11,00,121,,14,13,172,05*62",
		dst:      &GSV{},
		want: &GSV{
			Type:             "GPGSV",
			Messages:         1,
			MessageNumber:    1,
			SatellitesInView: 13,
			Satellite0PRN:    2,
			Elevation0:       2,
			Azimuth0:         213,
			SNR0:             0,
			Satellite1PRN:    3,
			Elevation1:       -3,
			Azimuth1:         0,
			SNR1:             0,
			Satellite2PRN:    11,
			Elevation2:       0,
			Azimuth2:         121,
			SNR2:             0,
			Satellite3PRN:    14,
			Elevation3:       13,
			Azimuth3:         172,
			SNR3:             5,
			Checksum:         0x62,
		},
	},
	{
		sentence: "$GPHDT,1,T*2A",
		dst:      &HDT{},
		want: &HDT{
			Type:     "GPHDT",
			Heading:  1,
			Checksum: 0x2a,
		},
	},
	{
		sentence: "$GPR00,EGLL,EGLM,EGTB,EGUB,EGTK,MBOT,EGTB,,,,,,,*58",
		dst:      &R00{},
		want: &R00{
			Type:     "GPR00",
			WP0:      "EGLL",
			WP1:      "EGLM",
			WP2:      "EGTB",
			WP3:      "EGUB",
			WP4:      "EGTK",
			WP5:      "MBOT",
			WP6:      "EGTB",
			Checksum: 0x58,
		},
	},
	{
		sentence: "$GPR00,MINST,CHATN,CHAT1,CHATW,CHATM,CHATE,003,004,005,006,007,,,*05",
		dst:      &R00{},
		want: &R00{
			Type:     "GPR00",
			WP0:      "MINST",
			WP1:      "CHATN",
			WP2:      "CHAT1",
			WP3:      "CHATW",
			WP4:      "CHATM",
			WP5:      "CHATE",
			WP6:      "003",
			WP7:      "004",
			WP8:      "005",
			WP9:      "006",
			WP10:     "007",
			Checksum: 0x5,
		},
	},
	{
		sentence: "$GPRMA,A,1234.56,N,12345.67,W,,,12.3,123,12.3,W*6D",
		dst:      &RMA{},
		want: &RMA{
			Type:     "GPRMA",
			Status:   "A",
			Latitude: 12.575999999999999, NorthSouth: "N",
			Longitude: 123.76116666666667, EastWest: "W",
			Speed:            12.3,
			CourseOverGround: 123,
			Variation:        12.3, VarDirection: "W",
			Checksum: 0x6d,
		},
	},
	{
		sentence: "$GPRMB,A,0.66,L,003,004,4917.24,N,12309.57,W,001.3,052.5,000.5,V*20",
		dst:      &RMB{},
		want: &RMB{
			Type:             "GPRMB",
			Status:           "A",
			CrosstrackError:  0.66,
			CorrectDirection: "L",
			Origin:           "003",
			Destination:      "004",
			Latitude:         49.28733333333333, NorthSouth: "N",
			Longitude: 123.1595, EastWest: "W",
			RangeToDestination:   1.3,
			BearingToDestination: 52.5,
			ClosingVelocity:      0.5,
			ArrivalStatus:        "V",
			Checksum:             0x20,
		},
	},
	{
		sentence: "$GPRMB,A,4.08,L,EGLL,EGLM,5130.02,N,00046.34,W,004.6,213.9,122.9,A*3D",
		dst:      &RMB{},
		want: &RMB{
			Type:             "GPRMB",
			Status:           "A",
			CrosstrackError:  4.08,
			CorrectDirection: "L",
			Origin:           "EGLL",
			Destination:      "EGLM",
			Latitude:         51.50033333333334, NorthSouth: "N",
			Longitude: 0.7723333333333334, EastWest: "W",
			RangeToDestination:   4.6,
			BearingToDestination: 213.9,
			ClosingVelocity:      122.9,
			ArrivalStatus:        "A",
			Checksum:             0x3d,
		},
	},
	{
		sentence: "$GPRMC,081836,A,3751.65,S,14507.36,E,000.0,360.0,130998,011.3,E*62",
		dst:      &RMC{},
		want: &RMC{
			Type:     "GPRMC",
			Time:     time.Date(0, 1, 1, 8, 18, 36, 0, time.UTC),
			Status:   "A",
			Latitude: 37.86083333333333, NorthSouth: "S",
			Longitude: 145.12266666666667, EastWest: "E",
			Speed:             0,
			Track:             360,
			Date:              time.Date(1998, 9, 13, 0, 0, 0, 0, time.UTC),
			MagneticVariation: 11.3, VarDirection: "E",
			Checksum: 0x62,
		},
	},
	{
		sentence: "$GPRMC,225446,A,4916.45,N,12311.12,W,000.5,054.7,191194,020.3,E*68",
		dst:      &RMC{},
		want: &RMC{
			Type:     "GPRMC",
			Time:     time.Date(0, 1, 1, 22, 54, 46, 0, time.UTC),
			Status:   "A",
			Latitude: 49.27416666666666, NorthSouth: "N",
			Longitude: 123.18533333333335, EastWest: "W",
			Speed:             0.5,
			Track:             54.7,
			Date:              time.Date(1994, 11, 19, 0, 0, 0, 0, time.UTC),
			MagneticVariation: 20.3, VarDirection: "E",
			Checksum: 0x68,
		},
	},
	{
		sentence: "$GPRMC,220516,A,5133.82,N,00042.24,W,173.8,231.8,130694,004.2,W*70",
		dst:      &RMC{},
		want: &RMC{
			Type:     "GPRMC",
			Time:     time.Date(0, 1, 1, 22, 05, 16, 0, time.UTC),
			Status:   "A",
			Latitude: 51.56366666666667, NorthSouth: "N",
			Longitude: 0.7040000000000001, EastWest: "W",
			Speed:             173.8,
			Track:             231.8,
			Date:              time.Date(1994, 06, 13, 0, 0, 0, 0, time.UTC),
			MagneticVariation: 4.2, VarDirection: "W",
			Checksum: 0x70,
		},
	},
	{
		sentence: "$GPTRF,053220.03,051197,4916.45,N,12311.12,W,1.2,3.4,5.6,7.8,SAT",
		dst:      &TRF{},
		want: &TRF{
			Type:     "GPTRF",
			Time:     time.Date(0, 1, 1, 05, 32, 20, 30e6, time.UTC),
			Date:     time.Date(1997, 11, 05, 0, 0, 0, 0, time.UTC),
			Latitude: 49.27416666666666, NorthSouth: "N",
			Longitude: 123.18533333333335, EastWest: "W",
			Elevation:       1.2,
			Iterations:      3.4,
			DoplerIntervals: 5.6,
			UpdateDistance:  7.8,
			Satellite:       "SAT",
		},
	},
	{
		sentence: "$GPSTN,3",
		dst:      &STN{},
		want: &STN{
			Type:   "GPSTN",
			Talker: 3,
		},
	},
	{
		sentence: "$GPVBW,1.2,3.4,A,5.6,7.8,A",
		dst:      &VBW{},
		want: &VBW{
			Type:                    "GPVBW",
			LongitudinalWaterSpeed:  1.2,
			TransverseWaterSpeed:    3.4,
			WaterSpeedStatus:        "A",
			LongitudinalGroundSpeed: 5.6,
			TransverseGroundSpeed:   7.8,
			GroundSpeedStatus:       "A",
		},
	},
	{
		sentence: "$GPVTG,360.0,T,348.7,M,000.0,N,000.0,K*43",
		dst:      &VTG{},
		want: &VTG{
			Type:          "GPVTG",
			TrackTrue:     360,
			TrackMagnetic: 348.7,
			SpeedKnots:    0,
			SpeedKph:      0,
			Checksum:      0x43,
		},
	},
	{
		sentence: "$GPVTG,054.7,T,034.4,M,005.5,N,010.2,K",
		dst:      &VTG{},
		want: &VTG{
			Type:          "GPVTG",
			TrackTrue:     54.7,
			TrackMagnetic: 34.4,
			SpeedKnots:    5.5,
			SpeedKph:      10.2,
		},
	},
	{
		sentence: "$GPVTG,78.9,T,,,1.23,N,4.56,K*1C",
		dst:      &VTG{},
		want: &VTG{
			Type:          "GPVTG",
			TrackTrue:     78.9,
			TrackMagnetic: 0,
			SpeedKnots:    1.23,
			SpeedKph:      4.56,
			Checksum:      0x1c,
		},
	},
	{
		sentence: "$GPWPL,4917.16,N,12310.64,W,003*65",
		dst:      &WPL{},
		want: &WPL{
			Type:     "GPWPL",
			Latitude: 49.285999999999994, NorthSouth: "N",
			Longitude: 123.17733333333332, EastWest: "W",
			Waypoint: "003",
			Checksum: 0x65,
		},
	},
	{
		sentence: "$GPWPL,5128.62,N,00027.58,W,EGLL*59",
		dst:      &WPL{},
		want: &WPL{
			Type:     "GPWPL",
			Latitude: 51.477000000000004, NorthSouth: "N",
			Longitude: 0.4596666666666666, EastWest: "W",
			Waypoint: "EGLL",
			Checksum: 0x59,
		},
	},
	{
		sentence: "$GPXTE,A,A,0.67,L,N",
		dst:      &XTE{},
		want: &XTE{
			Type:            "GPXTE",
			GeneralWarning:  "A",
			LockFlag:        "A",
			CrossTrackError: 0.67,
			Steer:           "L",
			Units:           "N",
		},
	},
	{
		sentence: "$GPXTE,A,A,4.07,L,N*6D",
		dst:      &XTE{},
		want: &XTE{
			Type:            "GPXTE",
			GeneralWarning:  "A",
			LockFlag:        "A",
			CrossTrackError: 4.07,
			Steer:           "L",
			Units:           "N",
			Checksum:        0x6d,
		},
	},
	{
		sentence: "$GPZDA,173958.45,01,05,1970,10,30",
		dst:      &ZDA{},
		want: &ZDA{
			Type:            "GPZDA",
			Time:            time.Date(0, 1, 1, 17, 39, 58, 450e6, time.UTC),
			Day:             1,
			Month:           5,
			Year:            1970,
			TimeZone:        10,
			TimeZoneMinutes: 30,
		},
	},
	{
		sentence: "$PGRME,15.0,M,45.0,M,25.0,M*1C",
		dst:      &RME{},
		want: &RME{
			Type:     "PGRME",
			HPE:      15,
			VPE:      45,
			OSEPE:    25,
			Checksum: 0x1c,
		},
	},
	{
		sentence: "$PGRMM,Astrln Geod '66*51",
		dst:      &RMM{},
		want: &RMM{
			Type:     "PGRMM",
			MapDatum: "Astrln Geod '66",
			Checksum: 0x51,
		},
	},
	{
		sentence: "$PGRMM,NAD27 Canada*2F",
		dst:      &RMM{},
		want: &RMM{
			Type:     "PGRMM",
			MapDatum: "NAD27 Canada",
			Checksum: 0x2f,
		},
	},
	{
		sentence: "$PGRMZ,246,f,3*1B",
		dst:      &RMZ{},
		want: &RMZ{
			Type:                  "PGRMZ",
			Altitude:              246,
			PositionFixDimensions: 3,
			Checksum:              0x1b},
	},
	{
		sentence: "$PGRMZ,93,f,3*21",
		dst:      &RMZ{},
		want: &RMZ{
			Type:                  "PGRMZ",
			Altitude:              93,
			PositionFixDimensions: 3,
			Checksum:              0x21},
	},
	{
		sentence: "$PGRMZ,201,f,3*18",
		dst:      &RMZ{},
		want: &RMZ{
			Type:                  "PGRMZ",
			Altitude:              201,
			PositionFixDimensions: 3,
			Checksum:              0x18},
	},
	{
		sentence: "$PSLIB,,,J*22",
		dst:      &LIB{},
		want: &LIB{
			Type:        "PSLIB",
			RequestType: "J",
			Checksum:    0x22},
	},
	{
		sentence: "$PSLIB,,,K*23",
		dst:      &LIB{},
		want: &LIB{
			Type:        "PSLIB",
			RequestType: "K",
			Checksum:    0x23},
	},
	{
		sentence: "$PSLIB,320.0,200*59",
		dst:      &LIB{},
		want: &LIB{
			Type:      "PSLIB",
			Frequency: 320,
			BitRate:   200,
			Checksum:  0x59},
	},
	{
		sentence: "$GNGNS,014035.00,4332.69262,S,17235.48549,E,RR,13,0.9,25.63,11.24,,*70",
		dst:      &GNS{},
		want: &GNS{
			Type:      "GNGNS",
			Timestamp: time.Date(0, 1, 1, 1, 40, 35, 0, time.UTC),
			Latitude:  43.54487699999999, NorthSouth: "S",
			Longitude: 172.59142483333332, EastWest: "E",
			Mode:             "RR",
			Satellites:       13,
			HDOP:             0.9,
			Altitude:         25.63,
			Separation:       11.24,
			Age:              0,
			ReferenceStation: 0x0,
			Checksum:         0x70,
		},
	},
	{
		sentence: "$GPGNS,014035.00,,,,,,8,,,,1.0,23*76",
		dst:      &GNS{},
		want: &GNS{
			Type:             "GPGNS",
			Timestamp:        time.Date(0, 1, 1, 1, 40, 35, 0, time.UTC),
			Satellites:       8,
			Age:              1,
			ReferenceStation: 23,
			Checksum:         0x76,
		},
	},
	{
		sentence: "$GLGNS,014035.00,,,,,,5,,,,1.0,23*67",
		dst:      &GNS{},
		want: &GNS{
			Type:             "GLGNS",
			Timestamp:        time.Date(0, 1, 1, 1, 40, 35, 0, time.UTC),
			Satellites:       5,
			Age:              1,
			ReferenceStation: 23,
			Checksum:         0x67,
		},
	},
	{
		sentence: "$GPTHS,1.2,A*34",
		dst:      &THS{},
		want: &THS{
			Type:     "GPTHS",
			Heading:  1.2,
			Status:   "A",
			Checksum: 0x34,
		},
	},
	{
		sentence: "!AIVDM,1,1,,B,177KQJ5000G?tO`K>RA1wUbN0TKH,0*5C",
		dst:      &VDMVDO{},
		want: &VDMVDO{
			Type:           "AIVDM",
			Fragments:      1,
			FragmentNumber: 1,
			ChannelCode:    "B",
			Data:           "177KQJ5000G?tO`K>RA1wUbN0TKH",
			Padding:        0x0,
			Checksum:       0x5c,
		},
	},
	{
		sentence: "!AIVDM,2,1,3,B,55P5TL01VIaAL@7WKO@mBplU@<PDhh000000001S;AJ::4A80?4i@E53,0*3E",
		dst:      &VDMVDO{},
		want: &VDMVDO{
			Type:           "AIVDM",
			Fragments:      2,
			FragmentNumber: 1,
			MessageID:      "3",
			ChannelCode:    "B",
			Data:           "55P5TL01VIaAL@7WKO@mBplU@<PDhh000000001S;AJ::4A80?4i@E53",
			Padding:        0x0,
			Checksum:       0x3e,
		},
	},
	{
		sentence: "!AIVDM,2,2,3,B,1@0000000000000,2*55",
		dst:      &VDMVDO{},
		want: &VDMVDO{
			Type:           "AIVDM",
			Fragments:      2,
			FragmentNumber: 2,
			MessageID:      "3",
			ChannelCode:    "B",
			Data:           "1@0000000000000",
			Padding:        0x2,
			Checksum:       0x55,
		},
	},
}

func TestParseTo(t *testing.T) {
	for _, test := range parseTests {
		err := ParseTo(test.dst, test.sentence)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(test.dst, test.want) {
			t.Errorf("unexpected result:\ngot: %#v\nwant:%#v", test.dst, test.want)
		}
	}
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		got, err := Parse(test.sentence)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		want := reflect.ValueOf(test.want).Elem().Interface()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("unexpected result:\ngot: %#v\nwant:%#v", got, want)
		}
	}
}

var aisArmorTests = []struct {
	payload  string
	padding  int
	want     []byte
	wantBits []byte
}{
	{
		payload:  "",
		want:     nil,
		wantBits: nil,
	},
	{
		payload: "177KQJ5000G?tO`K>RA1wUbN0TKH", padding: 0,
		want: []byte{
			0x01, 0x07, 0x07, 0x1b, 0x21, 0x1a, 0x05, 0x00,
			0x00, 0x00, 0x17, 0x0f, 0x3c, 0x1f, 0x28, 0x1b,
			0x0e, 0x22, 0x11, 0x01, 0x3f, 0x25, 0x2a, 0x1e,
			0x00, 0x24, 0x1b, 0x18,
		},
		wantBits: []byte{
			0x04, 0x71, 0xdb, 0x85, 0xa1, 0x40, 0x00, 0x05,
			0xcf, 0xf1, 0xfa, 0x1b, 0x3a, 0x24, 0x41, 0xfe,
			0x5a, 0x9e, 0x02, 0x46, 0xd8,
		},
	},
	{
		payload: "55P5TL01VIaAL@7WKO@mBplU@<PDhh000000001S;AJ::4A80?4i@E53", padding: 0,
		want: []byte{
			0x05, 0x05, 0x20, 0x05, 0x24, 0x1c, 0x00, 0x01,
			0x26, 0x19, 0x29, 0x11, 0x1c, 0x10, 0x07, 0x27,
			0x1b, 0x1f, 0x10, 0x35, 0x12, 0x38, 0x34, 0x25,
			0x10, 0x0c, 0x20, 0x14, 0x30, 0x30, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x23,
			0x0b, 0x11, 0x1a, 0x0a, 0x0a, 0x04, 0x11, 0x08,
			0x00, 0x0f, 0x04, 0x31, 0x10, 0x15, 0x05, 0x03,
		},
		wantBits: []byte{
			0x14, 0x58, 0x05, 0x91, 0xc0, 0x01, 0x99, 0x9a,
			0x51, 0x71, 0x01, 0xe7, 0x6d, 0xf4, 0x35, 0x4b,
			0x8d, 0x25, 0x40, 0xc8, 0x14, 0xc3, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x63, 0x2d, 0x16,
			0x8a, 0x28, 0x44, 0x48, 0x00, 0xf1, 0x31, 0x41,
			0x51, 0x43,
		},
	},
	{
		payload: "1@0000000000000", padding: 2,
		want: []byte{
			0x01, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
		wantBits: []byte{
			0x01, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00,
		},
	},
}

func TestDeArmorAIS(t *testing.T) {
	for _, test := range aisArmorTests {
		got, err := DeArmorAIS(test.payload)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !bytes.Equal(got, test.want) {
			t.Errorf("unexpected result:\ngot: %#v\nwant:%#v", got, test.want)
			continue
		}
		end := len(test.payload)*6 - test.padding
		gotBits := AISBitField(got, 0, end)
		if !bytes.Equal(gotBits, test.wantBits) {
			t.Errorf("unexpected result:\ngot: %08b\nwant:%08b", gotBits, test.wantBits)
		}
	}
}
