package utils

import (
	"os"

	"github.com/spf13/viper"
	"github.com/squareboat/splinter/constants"
)

func CheckDirExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func GetStringFromFlagOrConfig(key string, flag string) string {
	if viper.IsSet(flag) {
		return viper.GetString(flag)
	}
	return viper.GetString(key)
}
func GetIntFromFlagOrConfig(key string, flag string) int {
	if viper.IsSet(flag) {
		return viper.GetInt(flag)
	}
	return viper.GetInt(key)
}

func GetConfigFile() string {
	defaultConfig := constants.DEFAULT_USER_CONFIG_FILE
	if viper.IsSet("default_config") {
		defaultConfig = viper.GetString("default_config")
	}
	return defaultConfig
}
