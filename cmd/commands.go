package cmd

import (
	"fmt"
	"log"

	"github.com/squareboat/splinter/config"
	"github.com/squareboat/splinter/constants"
	"github.com/squareboat/splinter/logger"
	"github.com/squareboat/splinter/migrate"
	"github.com/squareboat/splinter/parser"
	"github.com/squareboat/splinter/runner"
	"github.com/squareboat/splinter/utils"

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
			if len(args) > 0 {
				logger.Log.Infof("Migrating files %v", args)
				// Run the migrations from the files passed in the command line
			}

			migrator, err := migrate.NewMigrate(constants.MIGRATION_UP)
			if err != nil {
				log.Fatal(err)
			}
			err = migrator.Up()
			if err != nil {
				log.Fatal(err)
			}
		},
	},
	"rollback": {
		Use:     "rollback",
		Short:   "Rollback the last migration.",
		Aliases: []string{"down"},
		Long:    `Rollback the last migration that was applied to the database.`,
		Run: func(cmd *cobra.Command, args []string) {
			logger.Log.Info("Running rollback")
			migrator, err := migrate.NewMigrate(constants.MIGRATION_DOWN)
			if err != nil {
				log.Fatal(err)
			}
			err = migrator.Down()
			if err != nil {
				log.Fatal(err)
			}
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
		Example: "splinter config <key1> <key2> <key3>\nsplinter config ",
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
	"unlock": {
		Use:   "unlock",
		Short: "Unlock the database.",
		Long:  `Unlock the database if in case previous locks were not removed due to a crash.`,
		Run: func(cmd *cobra.Command, args []string) {
			runner.UnlockDB(config.GetDbUri())
		},
	},
}

func SetFlags(rootCmd *cobra.Command) {
	// Sub Commands Flags Go Here
	MigratorCommands["rollback"].PersistentFlags().Int("n", 1, "Limit the number of rollback migrations.")

	// Global Flags
	rootCmd.PersistentFlags().String(constants.URI_FLAG, "", "DB Connection URI")
	rootCmd.PersistentFlags().String(constants.USERNAME_FLAG, "", "DB Connection Username")
	rootCmd.PersistentFlags().String(constants.PASSWORD_FLAG, "", "DB Connection Password")
	rootCmd.PersistentFlags().String(constants.HOST_FLAG, "", "DB Connection Host")
	rootCmd.PersistentFlags().Int(constants.PORT_FLAG, 0, "DB Connection Port")
	rootCmd.PersistentFlags().String(constants.DB_NAME_FLAG, "", "DB Connection Database Name")
	rootCmd.PersistentFlags().String(constants.MIGRATION_PATH_FLAG, constants.DEFAULT_MIGRATION_PATH, "Path Where Migrations are stored")
	rootCmd.PersistentFlags().String(constants.USER_CONFIG_FLAG, utils.GetConfigFile(), "Path of config file [env|YAML|JSON|TOML|INI]")

	// Bind Flags to viper in order to get the value of flags from command line
	viper.BindPFlags(rootCmd.PersistentFlags())
	for _, cmd := range MigratorCommands {
		viper.BindPFlags(cmd.PersistentFlags())
	}
}
