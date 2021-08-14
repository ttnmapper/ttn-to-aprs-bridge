package aprs

import (
	"fmt"
	"math"
)

/*
	!4903.50N/07201.75W-Test /A=001234 = no timestamp, no APRS messaging, altitude = 1234 ft.
	/092345z4903.50N/07201.75W>Test1234 = with timestamp, no APRS messaging, zulu time, with comment.
	@092345/4903.50N/07201.75W>088/036 = with timestamp, with APRS messaging, local time, course/speed.

	Altitude in Comment Text — The comment may contain an altitude value,
	in the form /A=aaaaaa, where aaaaaa is the altitude in feet. For example:
	/A=001234. The altitude may appear anywhere in the comment.
*/

type Position struct {
	Latitude  float64 // degrees
	Longitude float64 // degrees
	Altitude  float64 // meter

	Speed   float64 // use SI units, ie meter per seconds
	Heading float64 // degrees

	SymbolTable string
	Symbol      string

	Message string
}

// Zero zeroes all measurements in the position payload.
func (p *Position) Zero() {
	p.Latitude = 0
	p.Longitude = 0
	p.Altitude = 0

	p.SymbolTable = "/"
	p.Symbol = ">"

	p.Message = ""
}

func MeterToFeet(meter float64) int {
	return int(math.Round(meter * 3.28084))
}

// String returns an APRS packet for the provided measurements.
func (p Position) String() (s string) {
	// Base prefix
	latDeg, latMin, latHem := decToDMS(p.Latitude, [2]string{"N", "S"})
	lonDeg, lonMin, lonHem := decToDMS(p.Longitude, [2]string{"E", "W"})

	// no timestamp, no APRS messaging, altitude = 1234 ft
	// !4903.50N/07201.75W-Test /A=001234
	s = "!"
	s += fmt.Sprintf("%02.0f%05.2f%s", latDeg, latMin, latHem)
	if p.Latitude == 0 && p.Longitude == 0 {
		s += "\\"
	} else {
		s += p.SymbolTable
	}
	s += fmt.Sprintf("%03.0f%05.2f%s", lonDeg, lonMin, lonHem)
	if p.Latitude == 0 && p.Longitude == 0 {
		s += "."
	} else {
		s += p.Symbol
	}
	if p.Altitude != 0 {
		s += fmt.Sprintf("/A=%06d", MeterToFeet(p.Altitude))
	}
	s += p.Message

	return s
}
