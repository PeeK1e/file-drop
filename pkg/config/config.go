package config

import (
	"fmt"

	"github.com/alecthomas/kingpin/v2"
)

var config Config = Config{
	HttpServer: &HttpServer{
		ListenAddress: kingpin.Flag("listen", "HTTP Listen Address").
			Short('l').
			Envar("FD_HTTP_ADDRESS").
			Default("0.0.0.0:8080").
			String(),
	},
	DbSettings: &DatabaseSettings{
		HOSTNAME: kingpin.
			Flag("db-host", "Database Host").
			Short('H').
			Envar("FD_DB_HOST").
			Default("postgres").
			String(),
		PORT: kingpin.
			Flag("db-port", "Database Host").
			Short('p').
			Envar("FD_DB_PORT").
			Default("5432").
			String(),
		USERNAME: kingpin.
			Flag("db-user", "Database User").
			Short('U').
			Envar("FD_DB_USER").
			Default("postgres").
			String(),
		PASSWORD: kingpin.
			Flag("db-password", "Database Password").
			Short('P').
			Envar("FD_DB_PASSWORD").
			Default("s3cr3t!").
			String(),
		DATABASENAME: kingpin.
			Flag("db-name", "Database Name").
			Short('N').
			Envar("FD_DB_NAME").
			Default("postgresdb").
			String(),
		SSL: kingpin.
			Flag("db-ssl-mode", "Database SSL Mode").
			Short('S').
			Envar("FD_DB_SSL_MODE").
			Default("disable").
			Enum("verify-full", "verify-ca", "disable"),
		ConnectionRetryCount: kingpin.
			Flag("db-connect-retry", "Database connection retry count").
			Short('R').
			Envar("FD_DB_RETRY").
			Default("5").
			Int(),
	},
}

func NewConfig() *Config {
	kingpin.Parse()
	return &config
}

func (db *DatabaseSettings) GetDatabaseString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		*db.HOSTNAME,
		*db.PORT,
		*db.USERNAME,
		*db.PASSWORD,
		*db.DATABASENAME,
		*db.SSL,
	)
}
