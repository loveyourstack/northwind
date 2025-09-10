package main

import (
	"context"
	"os"

	"github.com/loveyourstack/lys/lyspgdb"
	"github.com/loveyourstack/lys/lyspgmon"
	"github.com/spf13/cobra"
)

var adminCheckDdlCmd = &cobra.Command{
	Use:   "checkDDL",
	Short: "Checks database DDL",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		defer cliApp.Db.Close()

		// connect as owner (for permission to perform admin tasks)
		dbOwnerUser := lyspgdb.User{
			Name:     cliApp.Config.DbOwnerUser.Name,
			Password: cliApp.Config.DbOwnerUser.Password,
		}

		ctx := context.Background()
		ownerDb, err := lyspgdb.GetPool(ctx, cliApp.Config.Db, dbOwnerUser, cliApp.Config.General.AppName+" Cli")
		if err != nil {
			cliApp.ErrorLog.Error("lyspgdb.GetPool (db owner) failed: " + err.Error())
			os.Exit(1)
		}
		defer ownerDb.Close()

		err = lyspgmon.CheckDDL(ctx, ownerDb, cliApp.InfoLog, cliApp.ErrorLog)
		if err != nil {
			cliApp.ErrorLog.Error("lyspgmon.CheckDDL failed: " + err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	adminCmd.AddCommand(adminCheckDdlCmd)
}
