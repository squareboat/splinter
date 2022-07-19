package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/squareboat/splinter/config"
	"github.com/squareboat/splinter/constants"
	"github.com/squareboat/splinter/logger"
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

	// Add SubCommands to Root Command
	for _, cmd := range MigratorCommands {
		rootCmd.AddCommand(cmd)
	}
	SetFlags(rootCmd)
	cobra.OnInitialize(onInit)
}

func onInit() {
	configFile := viper.GetString("config")
	// Check if user provided config file exists
	wd, _ := os.Getwd()

	logger.Log.Info("Current working directory: ", wd)
	exists, _ := os.Stat(configFile)

	logger.Log.Info("Config file: ", configFile)
	logger.Log.Info("Config file Exists")

	if exists != nil {
		viper.SetEnvPrefix(constants.SPLINTER_KEY_PREFIX)
		viper.AddConfigPath(wd)
		viper.SetConfigFile(configFile)
		err := viper.MergeInConfig()
		if err != nil {
			logger.Log.Fatal("Error reading user config file, "+configFile, "\n", err)
		}
	} else {
		logger.Log.Info("Config file does not exist. ", configFile)
		logger.Log.Fatal("User Provided Config file not found.")
		os.Exit(1)
	}
	logger.Log.Info()
	for k, v := range viper.AllSettings() {
		logger.Log.Infof("Config - %s : %#v", k, v)
	}
	config.Load()
	MergeSplinterConfig()
	CreateMigrationsPathIfNotExists()

}

func MergeSplinterConfig() {
	homeDir, osErr := os.UserHomeDir()
	if osErr != nil {
		logger.Log.Fatal("Error getting user home directory.")
		os.Exit(1)
	}
	viper.SetConfigName(constants.SETTINGS_CONFIG_FILE_NAME)
	viper.AddConfigPath(homeDir)
	err := viper.MergeInConfig()
	if err != nil {
		os.WriteFile(homeDir+"/"+constants.SETTINGS_CONFIG_FILE_NAME+".json", []byte("{}"), 0644)
		logger.Log.Warn("Error reading splinter config file. Reseting..  \n", err)
	}
	viper.AutomaticEnv()
}

func CreateMigrationsPathIfNotExists() {
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
}
