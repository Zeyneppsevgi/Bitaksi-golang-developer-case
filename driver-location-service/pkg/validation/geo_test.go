package validation

import "testing"

func TestValidateCoordinates(t *testing.T) {
	tests := []struct {
		name    string
		lon     float64
		lat     float64
		wantErr bool
	}{
		{name: "valid", lon: 29.0, lat: 41.0, wantErr: false},
		{name: "invalid lon", lon: 200, lat: 41.0, wantErr: true},
		{name: "invalid lat", lon: 29.0, lat: 100, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCoordinates(tt.lon, tt.lat)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error=%v wantErr=%v", err, tt.wantErr)
			}
		})
	}
}
