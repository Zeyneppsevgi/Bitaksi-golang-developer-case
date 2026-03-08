package config

import "os"

type Config struct {
	Port            string
	MongoURI        string
	MongoDB         string
	MongoCollection string
	InternalAPIKey  string
	SeedOnStart     bool
	SeedFile        string
}

func Load() Config {
	return Config{
		Port:            getenv("DRIVER_PORT", getenv("PORT", "8080")),
		MongoURI:        getenv("MONGO_URI", "mongodb://mongo:27017"),
		MongoDB:         getenv("MONGO_DB", "driver_location"),
		MongoCollection: getenv("MONGO_COLLECTION", "driver_locations"),
		InternalAPIKey:  getenv("INTERNAL_API_KEY", "some-secret"),
		SeedOnStart:     getenv("SEED_ON_START", "false") == "true",
		SeedFile:        getenv("SEED_FILE", "data/drivers.csv"),
	}
}

func getenv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
