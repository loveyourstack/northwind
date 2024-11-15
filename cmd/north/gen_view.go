package main

import (
	"context"
	"fmt"
	"os"

	"github.com/loveyourstack/lys/lysgen"
	"github.com/loveyourstack/lys/lyspgdb"
	"github.com/spf13/cobra"
)

var genViewCmd = &cobra.Command{
	Use:   "view",
	Short: "generates PostgreSQL view from the supplied schema + table",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		defer cliApp.Db.Close()

		// connect as owner (for permission to view constraints)
		dbOwnerUser := lyspgdb.User{
			Name:     cliApp.Config.DbOwnerUser.Name,
			Password: cliApp.Config.DbOwnerUser.Password,
		}

		ctx := context.Background()
		ownerDb, err := lyspgdb.GetPool(ctx, cliApp.Config.Db, dbOwnerUser)
		if err != nil {
			cliApp.ErrorLog.Error("lyspgdb.GetPool (db owner) failed: " + err.Error())
			os.Exit(1)
		}
		defer ownerDb.Close()

		res, err := lysgen.View(cmd.Context(), ownerDb, args[0], args[1])
		if err != nil {
			cliApp.ErrorLog.Error("lysgen.View failed: " + err.Error())
			os.Exit(1)
		}

		fmt.Println(res)
	},
}

func init() {
	genCmd.AddCommand(genViewCmd)
}
