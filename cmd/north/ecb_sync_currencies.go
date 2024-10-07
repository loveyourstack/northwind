package main

import (
	"context"
	"fmt"
	"os"

	"github.com/loveyourstack/connectors/csyncdb"
	"github.com/spf13/cobra"
)

var ecbSyncCurrCmd = &cobra.Command{
	Use:   "currencies",
	Short: "Sync currencies from ECB API into database",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		defer cliApp.Db.Close()

		_, err := csyncdb.EcbCurrencies(context.Background(), cliApp.Db, cliApp.EcbClient)
		if err != nil {
			cliApp.ErrorLog.Error("csyncdb.EcbCurrencies failed: " + err.Error())
			os.Exit(1)
		}

		fmt.Println("done")
	},
}

func init() {
	ecbSyncCmd.AddCommand(ecbSyncCurrCmd)
}
