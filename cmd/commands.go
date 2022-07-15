package cmd

import (
	"fmt"

	"github.com/the-e3n/splinter/constants"
	"github.com/the-e3n/splinter/logger"
	"github.com/the-e3n/splinter/parser"
	"github.com/the-e3n/splinter/runner"

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

			connURL, err := cmd.Flags().GetString("conn")
			if err != nil {
				logger.Log.WithError(err)
				return
			}
			runner.Postgres(connURL, constants.MIGRATION_UP)
			logger.Log.Info("Migrate Command")

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
				fmt.Println(viper.AllSettings())
			}
			for _, arg := range args {
				fmt.Printf("Value of %s :- %s", arg, viper.GetString(arg))
			}
		},
	},
}

func SetFlags(rootCmd *cobra.Command) {
	// Migrate Command Flags
	MigratorCommands["migrate"].PersistentFlags().String("conn", "", "connection URI DB")

	// Global Flags
	rootCmd.PersistentFlags().StringVar(&configFile, "config", ".env", "Path of config file [env|YAML|JSON|TOML|INI]")
}
