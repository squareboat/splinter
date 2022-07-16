package cmd

import (
	"fmt"
	"log"

	"github.com/the-e3n/splinter/constants"
	"github.com/the-e3n/splinter/logger"
	"github.com/the-e3n/splinter/parser"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var MigratorCommands = map[string]*cobra.Command{
	"migrate": {
		Use:     "migrate",
		Aliases: []string{"up"},
		Short:   "Run all the migration.",
		Long:    `Run all the migration that are pending in the system to database.`,
		Run: func(cmd *cobra.Command, args []string) {
			totalQueries := []string{}
			filenames, _ := parser.GetMigrationFileNames()
			for _, filename := range filenames {
				queries, err := parser.ParseFile(filename, constants.MIGRATION_UP)
				if err != nil {
					log.Fatal(err)
				}
				totalQueries = append(totalQueries, queries...)
			}
			logger.Log.Info(fmt.Sprintf("Total Queries: %d", len(totalQueries)))
			for _, query := range totalQueries {
				logger.Log.Info(query)
			}

		},
	},
	"rollback": {
		Use:     "rollback",
		Short:   "Rollback all the migration.",
		Aliases: []string{"down"},
		Long:    `Rollback all the migration that are pending in the system to database.`,
		Run: func(cmd *cobra.Command, args []string) {
			logger.Log.Info("Running rollback")
		},
	},
	"create": {
		Use:     "create",
		Short:   "Create a new migration file.",
		Long:    `Create a new migration file.`,
		Example: "splinter create <filename1> <filename2>\nsplinter create create_user_table",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Please provide a migration name.")
				return
			}
			fmt.Println("Creating a new migration file.")
			parser.CreateMigrationFile(args)
		},
	},
	"config": {
		Use:     "config",
		Short:   "Show specified config value.",
		Long:    `Show specified config value.`,
		Example: "splinter show <key1> <key2> <key3>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				for key, value := range viper.AllSettings() {
					fmt.Printf("Value of %s = %#v\n", key, value)
				}
			}
			for _, arg := range args {
				fmt.Printf("Value of %s = %#v\n", arg, viper.GetString(arg))
			}
		},
	},
}

func SetFlags(rootCmd *cobra.Command) {
	// Sub Commands Flags Go Here

	// Global Flags
	rootCmd.PersistentFlags().String(constants.URI_FLAG, "", "connection URI DB")
	rootCmd.PersistentFlags().String(constants.USERNAME_FLAG, "", "DB Connection Username")
	rootCmd.PersistentFlags().String(constants.PASSWORD_FLAG, "", "DB Connection Password")
	rootCmd.PersistentFlags().String(constants.HOST_FLAG, "", "DB Connection Host")
	rootCmd.PersistentFlags().Int(constants.PORT_FLAG, 0, "DB Connection Port")
	rootCmd.PersistentFlags().String(constants.DB_NAME_FLAG, "", "DB Connection Database Name")
	rootCmd.PersistentFlags().String(constants.MIGRATION_PATH_FLAG, constants.DEFAULT_MIGRATION_PATH, "Path Where Migrations are stored")
	rootCmd.PersistentFlags().String(constants.USER_CONFIG_FLAG, constants.DEFAULT_USER_CONFIG_FILE, "Path of config file [env|YAML|JSON|TOML|INI]")

	// Bind Flags to viper in order to get the value of flags from command line
	viper.BindPFlags(rootCmd.PersistentFlags())
	for _, cmd := range MigratorCommands {
		viper.BindPFlags(cmd.PersistentFlags())
	}
}
