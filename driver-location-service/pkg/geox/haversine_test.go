package geox

import (
	"math"
	"testing"
)

func TestHaversine(t *testing.T) {
	tests := []struct {
		name     string
		lon1     float64
		lat1     float64
		lon2     float64
		lat2     float64
		expected float64 // in meters
		tol      float64 // tolerance in meters
	}{
		{
			name:     "Istanbul to Ankara",
			lon1:     28.9784, // Istanbul
			lat1:     41.0082,
			lon2:     32.8597, // Ankara
			lat2:     39.9334,
			expected: 351000, // Approx 351km
			tol:      2000,   // 2km tolerance for rough estimation
		},
		{
			name:     "Same point",
			lon1:     29.0,
			lat1:     41.0,
			lon2:     29.0,
			lat2:     41.0,
			expected: 0,
			tol:      0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Haversine(tt.lon1, tt.lat1, tt.lon2, tt.lat2)
			if math.Abs(got-tt.expected) > tt.tol {
				t.Errorf("Haversine() = %v, want %v (tol %v)", got, tt.expected, tt.tol)
			}
		})
	}
}
