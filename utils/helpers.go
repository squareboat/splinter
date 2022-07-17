package utils

import (
	"os"

	"github.com/spf13/viper"
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
	if viper.GetInt(flag) != 0 {
		return viper.GetInt(flag)
	}
	return viper.GetInt(key)
}
