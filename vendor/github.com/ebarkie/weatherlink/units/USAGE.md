# units
```
import "github.com/ebarkie/weatherlink/units"
```

Package units implements simple unit conversion functions.

## Usage

```go
const (
	Inches = 1.0
	Feet   = 0.0833333333
	Meters = 39.37008
)
```
From units.

```go
const (
	Centibars = 1.0
)
```
From units.

```go
const (
	MPH = 1.0
)
```
From units.

#### type Length

```go
type Length float64
```

Length is a length stored in Inches.

#### func (Length) Feet

```go
func (l Length) Feet() float64
```
Feet returns the length in feet.

#### func (Length) Inches

```go
func (l Length) Inches() float64
```
Inches returns the length in inches.

#### func (Length) Meters

```go
func (l Length) Meters() float64
```
Meters returns the length in meters.

#### func (Length) Millimeters

```go
func (l Length) Millimeters() float64
```
Millimeters returns the length in Millimeters.

#### type Pressure

```go
type Pressure float64
```

Pressure is a barometric pressure stored in Inches.

#### func (Pressure) Hectopascals

```go
func (p Pressure) Hectopascals() float64
```
Hectopascals returns the pressure in Hectopascals.

#### func (Pressure) Inches

```go
func (p Pressure) Inches() float64
```
Inches returns the pressure in Inches.

#### func (Pressure) Millibars

```go
func (p Pressure) Millibars() float64
```
Millibars returns the pressure in Millibars.

#### type SoilMoisture

```go
type SoilMoisture int
```

SoilMoisture is moisture in centibars of tension.

#### func (SoilMoisture) Percent

```go
func (m SoilMoisture) Percent(t SoilType) int
```
Percent returns the soil moisure as an approximated percentage.

#### type SoilType

```go
type SoilType uint
```

SoilType is the soil type used for calculating suction.

```go
const (
	Sand SoilType = iota
	SandyLoam
	Loam
	Clay
)
```
Soil types ranging from sand to clay.

#### type Speed

```go
type Speed float64
```

Speed is a speed stored in MPH.

#### func (Speed) Knots

```go
func (s Speed) Knots() float64
```
Knots returns the speed in Knots.

#### func (Speed) MPH

```go
func (s Speed) MPH() float64
```
MPH returns the speed in Miles per Hour.

#### func (Speed) MPS

```go
func (s Speed) MPS() float64
```
MPS returns the speed in Meters per Second.

#### type Temperature

```go
type Temperature float64
```

Temperature is a temperature stored in Fahrenheit.

#### func  Celsius

```go
func Celsius(c float64) Temperature
```
Celsius returns a temperature from a value in Celsius.

#### func  Fahrenheit

```go
func Fahrenheit(f float64) Temperature
```
Fahrenheit returns a temperature from a value in Fahrenheit.

#### func (Temperature) Celsius

```go
func (t Temperature) Celsius() float64
```
Celsius returns the temperature in Celsius.

#### func (Temperature) Fahrenheit

```go
func (t Temperature) Fahrenheit() float64
```
Fahrenheit returns the temperature in Fahrenheit.
