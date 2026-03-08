package geox

import (
	"math"
)

const earthRadiusMeter = 6371000

// Haversine calculates the great-circle distance between two points (lon1, lat1) and (lon2, lat2)
// on a sphere (Earth) using the Haversine formula.
// The result is returned in meters.
func Haversine(lon1, lat1, lon2, lat2 float64) float64 {
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	lat1Rad := lat1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMeter * c
}
