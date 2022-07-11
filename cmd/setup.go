package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	wd, _ := os.Getwd()
	logger.Log.Info("Current working directory: ", wd)
	exists, _ := os.Stat(configFile)
	logger.Log.Info("Config file: ", configFile)
	logger.Log.Info("Config file Exists: ", exists)
	if exists != nil {
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		if err != nil {
			logger.Log.Fatal("Error reading config file, " + configFile)
		}
	} else {
		logger.Log.Warn("Config file not found.")
		os.Exit(1)
		// fmt.Println("Created a default config file.")
		// ioutil.WriteFile(".env", []byte("SPLINTER_PATH=./migrations"), 0644)
		// viper.SetConfigFile(".env")
	}
	viper.AutomaticEnv()

}
