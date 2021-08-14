// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

// Speed is a speed stored in MPH.
type Speed float64

// From units.
const (
	MPH = 1.0
)

// Knots returns the speed in Knots.
func (s Speed) Knots() float64 {
	return s.MPH() * 0.8688
}

// MPH returns the speed in Miles per Hour.
func (s Speed) MPH() float64 {
	return float64(s)
}

// MPS returns the speed in Meters per Second.
func (s Speed) MPS() float64 {
	return s.MPH() * 0.44704
}
