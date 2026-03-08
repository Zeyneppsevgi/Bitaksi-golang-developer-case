package config

import "os"

type Config struct {
	Port                  string
	DriverLocationBaseURL string
	InternalAPIKey        string
	UserJWTSecret         string
}

func Load() Config {
	return Config{
		Port:                  getenv("MATCHING_PORT", getenv("PORT", "8081")),
		DriverLocationBaseURL: getenv("DRIVER_LOCATION_BASE_URL", "http://driver-location-service:8080"),
		InternalAPIKey:        getenv("INTERNAL_API_KEY", "some-secret"),
		UserJWTSecret:         getenv("USER_JWT_SECRET", "user-secret"),
	}
}

func getenv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
