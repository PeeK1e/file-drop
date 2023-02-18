package environment

import (
	"log"
	"os"
)

var (
	DatabaseSettings *databaseSettings
)

func init() {
	loadEnvironmentAndFlags()
}

// Loads all ENV vars for the programm startup
func loadEnvironmentAndFlags() {
	log.Print("Loading Environment")
	DatabaseSettings = &databaseSettings{
		HOSTNAME:     getEnvOrDefault("DATABASE_HOSTNAME", "db"),
		PORT:         getEnvOrDefault("DATABASE_PORT", "5432"),
		USERNAME:     getEnvOrDefault("DATABASE_USERNAME", "postgres"),
		PASSWORD:     getEnvOrDefault("DATABASE_PASSWORD", "password"),
		DATABASENAME: getEnvOrDefault("DATABASE_DATABASENAME", "uploads"),
		SSL:          getEnvOrDefault("DATABASE_SSL", "disable"),
	}
}

func getEnvOrDefault(envName string, defaultValue string) string {
	if value, ok := os.LookupEnv(envName); ok {
		return value
	}
	return defaultValue
}
