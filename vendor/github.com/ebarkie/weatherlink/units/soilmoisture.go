// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

// SoilMoisture is moisture in centibars of tension.
type SoilMoisture int

// From units.
const (
	Centibars = 1.0
)

// SoilType is the soil type used for calculating suction.
type SoilType uint

// Soil types ranging from sand to clay.
const (
	Sand SoilType = iota
	SandyLoam
	Loam
	Clay
)

// Percent returns the soil moisure as an approximated percentage.
func (m SoilMoisture) Percent(t SoilType) int {
	// Linear scale based on depletion of plant available water for each
	// soil type.

	// Scheduling Irrigations: When and How Much Water to Apply.
	// Division of Agriculture and Natural Resources Publication 3396.
	// University of California Irrigation Program.
	// University of California, Davis. pp. 106.
	var depletedCbs = []int{
		30,  // Sand/Loamy Sand
		50,  // Sandy Loam
		130, // Loam/Silt Loam
		170, // Clay Loam/Clay
	}

	dcb := depletedCbs[t]
	if int(m) > dcb {
		return 0
	}
	return 100 - int((float32(m)/float32(dcb))*100)
}
