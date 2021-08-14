// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

// Length is a length stored in Inches.
type Length float64

// From units.
const (
	Inches = 1.0
	Feet   = 0.0833333333
	Meters = 39.37008
)

// Feet returns the length in feet.
func (l Length) Feet() float64 {
	return l.Inches() / 12.0
}

// Inches returns the length in inches.
func (l Length) Inches() float64 {
	return float64(l)
}

// Meters returns the length in meters.
func (l Length) Meters() float64 {
	return l.Inches() * 0.025399999187200026
}

// Millimeters returns the length in Millimeters.
func (l Length) Millimeters() float64 {
	return l.Meters() * 1000.0
}
