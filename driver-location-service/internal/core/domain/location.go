package domain

import (
	"fmt"
	"time"

	"github.com/driver-location-service/pkg/validation"
)

type Point struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type DriverLocation struct {
	DriverID  string    `json:"driverId"`
	Location  Point     `json:"location"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type NearestDriver struct {
	DriverID  string  `json:"driverId"`
	DistanceM float64 `json:"distanceM"`
	Location  Point   `json:"location"`
}

func ValidateLocation(in DriverLocation) error {
	if in.DriverID == "" {
		return fmt.Errorf("%w: driverId is required", ErrValidation)
	}
	if in.Location.Type != "Point" {
		return fmt.Errorf("%w: location.type must be Point", ErrValidation)
	}
	if len(in.Location.Coordinates) != 2 {
		return fmt.Errorf("%w: coordinates must be [lon,lat]", ErrValidation)
	}
	if err := validation.ValidateCoordinates(in.Location.Coordinates[0], in.Location.Coordinates[1]); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}
