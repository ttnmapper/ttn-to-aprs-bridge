// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

// Pressure is a barometric pressure stored in Inches.
type Pressure float64

// From units.
const (
//Inches = 1.0
)

// Hectopascals returns the pressure in Hectopascals.
func (p Pressure) Hectopascals() float64 {
	return p.Millibars()
}

// Inches returns the pressure in Inches.
func (p Pressure) Inches() float64 {
	return float64(p)
}

// Millibars returns the pressure in Millibars.
func (p Pressure) Millibars() float64 {
	return p.Inches() * 33.8637526
}
