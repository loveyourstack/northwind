package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/loveyourstack/connectors/apiclients/ecbapi"
	"github.com/loveyourstack/connectors/csyncdb"
	"github.com/spf13/cobra"
)

// example: north ecb sync xr 2024-10-10 2024-10-21
var ecbSyncExRatesCmd = &cobra.Command{
	Use:   "xr",
	Short: "Sync exchange rates with base currency EUR from ECB API into database. Arguments are from and to date, in format YYYY-MM-DD.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		defer cliApp.Db.Close()

		startDate, err := time.Parse("2006-01-02", args[0])
		if err != nil {
			cliApp.ErrorLog.Error("time.Parse (start) failed: " + err.Error())
			os.Exit(1)
		}
		endDate, err := time.Parse("2006-01-02", args[1])
		if err != nil {
			cliApp.ErrorLog.Error("time.Parse (end) failed: " + err.Error())
			os.Exit(1)
		}

		// daily
		err = csyncdb.EcbExchangeRates(context.Background(), cliApp.Db, cliApp.EcbClient, "EUR", ecbapi.Daily, startDate, endDate)
		if err != nil {
			cliApp.ErrorLog.Error("csyncdb.EcbExchangeRates (Daily) failed: " + err.Error())
			os.Exit(1)
		}

		// monthly
		/*err = csyncdb.EcbExchangeRates(context.Background(), cliApp.Db, cliApp.EcbClient, "EUR", ecbapi.Monthly, startDate, endDate)
		if err != nil {
			cliApp.ErrorLog.Error("csyncdb.EcbExchangeRates (Monthly) failed: " + err.Error())
			os.Exit(1)
		}*/

		fmt.Println("done")
	},
}

func init() {
	ecbSyncCmd.AddCommand(ecbSyncExRatesCmd)
}
