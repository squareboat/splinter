package cmd

import (
	"fmt"

	"github.com/the-e3n/migrator/parser"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var MigratorCommands = []*cobra.Command{
	{
		Use:     "migrate",
		Aliases: []string{"up", "m"},
		Short:   "Run all the migration.",
		Long:    `Run all the migration that are pending in the system to database.`,
		Run: func(cmd *cobra.Command, args []string) {
			querys := parser.ParseAllMigrations()
			for _, query := range querys {
				fmt.Println(query.Up)
			}

		},
	},
	{
		Use:   "rollback",
		Short: "Rollback all the migration.",
		Long:  `Rollback all the migration that are pending in the system to database.`,
		Run: func(cmd *cobra.Command, args []string) {
			queries, err := parser.ParseRollbackMigration()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(queries)
		},
	},
	{
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
	{
		Use:     "show",
		Short:   "Show specified config value.",
		Long:    `Show specified config value.`,
		Example: "splinter show <key1> <key2> <key3>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
			}
			for _, arg := range args {
				fmt.Printf("Value of %s :- %s", arg, viper.GetString(arg))
			}
		},
	},
}
