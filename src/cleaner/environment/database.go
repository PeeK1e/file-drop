package environment

import "fmt"

type databaseSettings struct {
	HOSTNAME     string `json:"hostname"`
	PORT         string `json:"port"`
	USERNAME     string `json:"userName"`
	PASSWORD     string `json:"userPasswd"`
	DATABASENAME string `json:"DBName"`
	SSL          string `json:"SSLMode"`
}

// Formats and print the database connetion string for the postgres library.
//
// All values are set via the `DATABASE_*` environment variables.
func (db *databaseSettings) GetDatabaseString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.HOSTNAME,
		db.PORT,
		db.USERNAME,
		db.PASSWORD,
		db.DATABASENAME,
		db.SSL)
}
