package config

import (
	"github.com/the-e3n/splinter/constants"
	"github.com/the-e3n/splinter/logger"
	"github.com/the-e3n/splinter/utils"
)

type userConfiguration struct {
	host            string
	port            int
	user            string
	password        string
	dbname          string
	driver          string
	uri             string
	migrations_path string
}

// Loads user config
func (userConfig *userConfiguration) load() {
	userConfig.host = utils.GetStringFromFlagOrConfig(constants.HOST, constants.HOST_FLAG)
	userConfig.port = utils.GetIntFromFlagOrConfig(constants.PORT, constants.PORT_FLAG)
	userConfig.user = utils.GetStringFromFlagOrConfig(constants.USER, constants.USERNAME_FLAG)
	userConfig.password = utils.GetStringFromFlagOrConfig(constants.PASSWORD, constants.PASSWORD_FLAG)
	userConfig.dbname = utils.GetStringFromFlagOrConfig(constants.DB_NAME, constants.DB_NAME_FLAG)
	userConfig.driver = utils.GetStringFromFlagOrConfig(constants.DRIVER, constants.DRIVER_FLAG)
	userConfig.uri = utils.GetStringFromFlagOrConfig(constants.URI, constants.URI_FLAG)
	userConfig.migrations_path = utils.GetStringFromFlagOrConfig(constants.MIGRATIONS_PATH, constants.MIGRATION_PATH_FLAG)

	logger.Log.Info("User Config Loaded Sucessfully.")
}
