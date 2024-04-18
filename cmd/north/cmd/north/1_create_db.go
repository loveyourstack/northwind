package north

import (
	"context"
	"os"

	"github.com/loveyourstack/lys/lyspgdb"
	"github.com/loveyourstack/northwind/internal/const/appenv"
	"github.com/loveyourstack/northwind/sql/ddl"
	"github.com/spf13/cobra"
)

// expects db users to have been created first (create_users.sql)

var createDbCmd = &cobra.Command{
	Use:   "createDb",
	Short: "Creates database. Drops existing if present.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		defer cliApp.Db.Close()

		// only possible in dev env
		if cliApp.Config.General.Env != appenv.Dev {
			cliApp.ErrorLog.Error("this command may only be used in dev environment")
			os.Exit(1)
		}

		dbSuperUser := lyspgdb.User{
			Name:     cliApp.Config.DbSuperUser.Name,
			Password: cliApp.Config.DbSuperUser.Password,
		}

		// (re-)create test db
		ctx := context.Background()
		if err := lyspgdb.CreateLocalDb(ctx, ddl.SQLAssets, cliApp.Config.Db, dbSuperUser, cliApp.Config.DbOwnerUser, true, false,
			nil, cliApp.InfoLog); err != nil {
			cliApp.ErrorLog.Error("lyspgdb.CreateLocalDb failed: " + err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(createDbCmd)
}
