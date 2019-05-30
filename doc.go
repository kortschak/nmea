// Copyright ©2019 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package nmea implements generic parsing of NMEA 0183 sentences.
//
// The package provides a parsing function that iterates over a NMEA sentences
// and fills fields of a struct with sequential values according to methods
// specified in field tags with the name "nmea". A destination struct must have
// an initial field with the name of the NMEA sentence and a struct tag `nmea:"type"`.
//
// Parsing methods that are available are:
//
//  - "number": set the field to a number parsed from the NMEA value
//  - "string": set the field to the literal NMEA value
//  - "latlon": set the field to a latitude or longitude parsed from the NMEA value
//  - "date":   set the field to a data parsed from the NMEA value in the form ddmmyy.
//  - "time":   set the field to a time parsed from the NMEA value in the form hhmmss.ss.
//
// A special case method is "checksum" which will write the value of the sentence
// checksum if it is available.
package nmea
