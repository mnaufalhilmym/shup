package config

type Configuration struct {
	logLevel string

	webAddress string

	dbDriverName    string
	dbConnectionURL string
}

func (p Configuration) LogLevel() string {
	return p.logLevel
}

func (p Configuration) WebAddress() string {
	return p.webAddress
}

func (p Configuration) DBDriverName() string {
	return p.dbDriverName
}

func (p Configuration) DBConnectionURL() string {
	return p.dbConnectionURL
}

func NewConfiguration() Configuration {
	return Configuration{
		logLevel: "debug",
		webAddress: ":8080",
		dbDriverName: "sqlite",
		dbConnectionURL: "./sqlite.db",
	}
}
