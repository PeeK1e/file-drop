package environment

import (
	"fmt"
	"log"
	"os"
)

var (
	DatabaseSettings  *databaseSettings
	HttpServerOptions *httpServerOptions
)

func init() {
	load()
}

type httpServerOptions struct {
	listenAddress string
	port          string
}

type databaseSettings struct {
	HOSTNAME     string `json:"hostname"`
	PORT         string `json:"port"`
	USERNAME     string `json:"userName"`
	PASSWORD     string `json:"userPasswd"`
	DATABASENAME string `json:"DBName"`
	SSL          string `json:"SSLMode"`
}

func (s httpServerOptions) String() string {
	return s.listenAddress + ":" + s.port
}

// Formats and print the database connetion string for the postgres library.
//
// All values are set via the `DATABASE_*` environment variables.
func (db databaseSettings) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.HOSTNAME,
		db.PORT,
		db.USERNAME,
		db.PASSWORD,
		db.DATABASENAME,
		db.SSL)
}

func getEnvOrDefault(envName string, defaultValue string) string {
	if value, ok := os.LookupEnv(envName); ok {
		return value
	}
	return defaultValue
}

// Loads all ENV vars for the programm startup
func load() {
	log.Print("Loading Environment")
	DatabaseSettings = &databaseSettings{
		HOSTNAME:     getEnvOrDefault("DATABASE_HOSTNAME", "db"),
		PORT:         getEnvOrDefault("DATABASE_PORT", "5432"),
		USERNAME:     getEnvOrDefault("DATABASE_USERNAME", "postgres"),
		PASSWORD:     getEnvOrDefault("DATABASE_PASSWORD", "password"),
		DATABASENAME: getEnvOrDefault("DATABASE_DATABASENAME", "uploads"),
		SSL:          getEnvOrDefault("DATABASE_SSL", "disable"),
	}

	HttpServerOptions = &httpServerOptions{
		listenAddress: getEnvOrDefault("HTTP_LISTEN_ADDRESS", "0.0.0.0"),
		port:          getEnvOrDefault("HTTP_PORT", "8080"),
	}
}
