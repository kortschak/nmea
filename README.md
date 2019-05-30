# Package nmea implements generic parsing of NMEA 0183 sentences.

The package provides a parsing function that iterates over a NMEA sentences and fills fields of a struct with sequential values according to methods specified in field tags with the name "nmea".