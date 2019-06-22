# Package nmea implements generic parsing of NMEA 0183 sentences.

[![Build Status](https://www.travis-ci.org/kortschak/nmea.svg?branch=master)](https://www.travis-ci.org/kortschak/nmea/branches) [![coveralls.io](https://coveralls.io/repos/kortschak/nmea/badge.svg?branch=master&service=github)](https://coveralls.io/github/kortschak/nmea?branch=master) [![GoDoc](https://godoc.org/github.com/kortschak/nmea?status.svg)](https://godoc.org/github.com/kortschak/nmea)

The package provides a parsing function that iterates over a NMEA sentences and fills fields of a struct with sequential values according to methods specified in field tags with the name "nmea".