package config

import (
	"github.com/the-e3n/splinter/constants"
	"github.com/the-e3n/splinter/utils"
)

type UserConfig struct {
	host            string
	port            int
	user            string
	password        string
	dbname          string
	driver          string
	uri             string
	migrations_path string
}

var userConfig = UserConfig{}

// Loads user config
func LoadUserConfig() {
	userConfig.host = utils.GetStringFromFlagOrConfig(constants.HOST, constants.HOST_FLAG)
	userConfig.port = utils.GetIntFromFlagOrConfig(constants.PORT, constants.PORT_FLAG)
	userConfig.user = utils.GetStringFromFlagOrConfig(constants.USER, constants.USERNAME_FLAG)
	userConfig.password = utils.GetStringFromFlagOrConfig(constants.PASSWORD, constants.PASSWORD_FLAG)
	userConfig.dbname = utils.GetStringFromFlagOrConfig(constants.DB_NAME, constants.DB_NAME_FLAG)
	userConfig.driver = utils.GetStringFromFlagOrConfig(constants.DRIVER, constants.DRIVER_FLAG)
	userConfig.uri = utils.GetStringFromFlagOrConfig(constants.URI, constants.URI_FLAG)
	userConfig.migrations_path = utils.GetStringFromFlagOrConfig(constants.MIGRATION_PATH, constants.MIGRATION_PATH_FLAG)
}

func GetDbHost() string {
	return userConfig.host
}

func GetDbPort() int {
	return userConfig.port
}

func GetDbUser() string {
	return userConfig.user
}

func GetDbPassword() string {
	return userConfig.password
}

func GetDbName() string {
	return userConfig.dbname
}

func GetDbDriver() string {
	return userConfig.driver
}

func GetDbUri() string {
	return userConfig.uri
}

func GetMigrationsPath() string {
	return userConfig.migrations_path
}
