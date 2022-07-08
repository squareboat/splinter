package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		if cmd.Name() == "migrate" {
			cmd.PersistentFlags().String("conn", "", "connection URI DB")
		}
	}

	// Add Flags to Root Command
	rootCmd.PersistentFlags().StringVar(&configFile, "config", ".env", "Path of config file [env|YAML|JSON|TOML|INI]")
}

func onInit() {
	exists, _ := os.Stat(configFile)
	if exists != nil {
		viper.SetConfigFile(configFile)
	} else {
		fmt.Println("Config file not found.")
		viper.SetConfigFile(".env")
		ioutil.WriteFile(".env", []byte("SPLINTER_PATH=./migrations"), 0644)
		fmt.Println("Created a default config file.")
	}
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s\n", configFile)
	}
}
