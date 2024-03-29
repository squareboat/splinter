package config

type config struct {
	userConfig userConfiguration
}

var configuration = &config{}

func Load() {
	configuration.userConfig.load()
}

func GetDbHost() string {
	return configuration.userConfig.host
}

func GetDbPort() int {
	return configuration.userConfig.port
}

func GetDbUser() string {
	return configuration.userConfig.user
}

func GetDbPassword() string {
	return configuration.userConfig.password
}

func GetDbName() string {
	return configuration.userConfig.dbname
}

func GetDbDriver() string {
	return configuration.userConfig.driver
}

func GetDbUri() string {
	return configuration.userConfig.uri
}

func GetMigrationsPath() string {
	return configuration.userConfig.migrations_path
}
