package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/squareboat/splinter/config"
	"github.com/squareboat/splinter/constants"
	"github.com/squareboat/splinter/logger"
	"github.com/squareboat/splinter/utils"
)

var rootCmd = &cobra.Command{
	Use:   "splinter",
	Short: "Splinter is a tool for migration",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	LoadSplinterConfig()
	// Add SubCommands to Root Command
	for _, cmd := range MigratorCommands {
		rootCmd.AddCommand(cmd)
	}
	cobra.OnInitialize(onInit)
	SetFlags(rootCmd)
}

func onInit() {
	configFile := utils.GetConfigFile()
	// Check if user provided config file exists
	wd, _ := os.Getwd()
	logger.Log.Debugf("Current working directory: %s", wd)
	exists, _ := os.Stat(configFile)

	if exists != nil {
		viper.SetEnvPrefix(constants.SPLINTER_KEY_PREFIX)
		viper.AddConfigPath(wd)
		viper.SetConfigFile(configFile)
		err := viper.MergeInConfig()
		if err != nil {
			logger.Log.Fatal("Error reading user config file, "+configFile, "\n", err)
		}
	} else {
		logger.Log.Fatal("Config file you provided does not exist. ", configFile)
		os.Exit(1)
	}
	for k, v := range viper.AllSettings() {
		logger.Log.Debugf("Config - %s : %#v", k, v)
	}
	config.Load()
	CreateMigrationsPathIfNotExists()

}

func LoadSplinterConfig() {
	logger.Log.Debug("Loading Splinter Config")
	homeDir, homeDirErr := os.UserHomeDir()
	if homeDirErr != nil {
		logger.Log.Fatal("Error getting user home directory. \n", homeDirErr)
	}
	configPath := fmt.Sprintf("%s/%s", homeDir, constants.SETTINGS_CONFIG_FILE_NAME)
	logger.Log.Debug("Config Path: ", configPath)
	exists, _ := os.Stat(configPath)
	if exists == nil {
		logger.Log.Errorf("Splinter Config file not found. %v", configPath)
	}
	viper.SetConfigFile(configPath)
	err := viper.MergeInConfig()
	if err != nil {
		logger.Log.Fatalf("Error reading config file, %v \n %v", configPath, err)
	}
	viper.AutomaticEnv()
}

func CreateMigrationsPathIfNotExists() {
	logger.Log.Debug("Creating Migrations Path if not exists")
	migrationsPath := config.GetMigrationsPath()
	if migrationsPath == "" {
		log.Fatal("Migrations path is not set.")
	}
	_, err := os.Stat(migrationsPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(migrationsPath, 0755)
		if err != nil {
			logger.Log.Fatal("Error creating migrations path. \n", err)
		}
	}
	logger.Log.Debug("Migrations Path: ", migrationsPath)
}
