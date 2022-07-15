package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/the-e3n/splinter/constants"
	"github.com/the-e3n/splinter/logger"
)

var rootCmd = &cobra.Command{
	Use:   "splinter",
	Short: "Splinter is a tool for migration",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
var configFile string

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(onInit)

	// Add SubCommands to Root Command
	for _, cmd := range MigratorCommands {
		rootCmd.AddCommand(cmd)
	}
	SetFlags(rootCmd)
}

func onInit() {

	// Check if user provided config file exists
	wd, _ := os.Getwd()
	logger.Log.Info("Current working directory: ", wd)
	exists, _ := os.Stat(configFile)
	logger.Log.Info("Config file: ", configFile)
	logger.Log.Info("Config file Exists")

	if exists != nil {
		// viper.SetEnvPrefix("")
		viper.AddConfigPath(wd)
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		if err != nil {
			logger.Log.Fatal("Error reading user config file, "+configFile, "\n", err)
		}
	} else {
		logger.Log.Fatal("User Provided Config file not found.")
		os.Exit(1)
	}
	LoadSplinterConfig()

}

func LoadSplinterConfig() {
	homeDir, osErr := os.UserHomeDir()
	if osErr != nil {
		logger.Log.Fatal("Error getting user home directory.")
		os.Exit(1)
	}
	viper.SetConfigName(constants.CONFIG_FILE_NAME)
	viper.AddConfigPath(homeDir)
	err := viper.MergeInConfig()
	if err != nil {
		logger.Log.Fatal("Error reading splinter config file. \n", err)
	}
	viper.AutomaticEnv()
}
