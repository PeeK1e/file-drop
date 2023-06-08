package config

type DatabaseSettings struct {
	HOSTNAME             *string `json:"hostname"`
	PORT                 *string `json:"port"`
	USERNAME             *string `json:"userName"`
	PASSWORD             *string `json:"userPasswd"`
	DATABASENAME         *string `json:"DBName"`
	SSL                  *string `json:"SSLMode"`
	ConnectionRetryCount *int    `json:"connectionRetryCount"`
}

type HttpServer struct {
	ListenAddress *string `json:"listenAddress"`
}

type Config struct {
	HttpServer *HttpServer       `json:"httpServer"`
	DbSettings *DatabaseSettings `json:"databaseSettings"`
}
