// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

// Temperature is a temperature stored in Fahrenheit.
type Temperature float64

// Celsius returns a temperature from a value in Celsius.
func Celsius(c float64) Temperature {
	return Temperature(c*1.8 + 32.0)
}

// Fahrenheit returns a temperature from a value in Fahrenheit.
func Fahrenheit(f float64) Temperature {
	return Temperature(f)
}

// Celsius returns the temperature in Celsius.
func (t Temperature) Celsius() float64 {
	return (t.Fahrenheit() - 32.0) * 5.0 / 9.0
}

// Fahrenheit returns the temperature in Fahrenheit.
func (t Temperature) Fahrenheit() float64 {
	return float64(t)
}
