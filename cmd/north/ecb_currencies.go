package main

import (
	"context"
	"fmt"
	"os"

	"github.com/loveyourstack/connectors/apiclients/ecbapi"
	"github.com/loveyourstack/connectors/csyncdb"
	"github.com/spf13/cobra"
)

var ecbCurrCmd = &cobra.Command{
	Use:   "currencies",
	Short: "S",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		defer cliApp.Db.Close()

		ecbClient := ecbapi.NewClient(cliApp.InfoLog, cliApp.ErrorLog)

		/*currs, err := ecbClient.GetApiCurrencies()
		if err != nil {
			cliApp.ErrorLog.Error("ecbClient.GetApiCurrencies failed: " + err.Error())
			os.Exit(1)
		}
		fmt.Println(currs)*/

		_, err := csyncdb.EcbCurrencies(context.Background(), cliApp.Db, ecbClient)
		if err != nil {
			cliApp.ErrorLog.Error("csyncdb.Currencies failed: " + err.Error())
			os.Exit(1)
		}

		/*exRates, err := ecbClient.GetAPIExchangeRates("EUR", ecbdata.Daily, time.Now().Add(-7*24*time.Hour), time.Now())
		if err != nil {
			cliApp.ErrorLog.Error("ecbClient.GetExchangeRates failed: " + err.Error())
			os.Exit(1)
		}
		fmt.Println(exRates)*/

		fmt.Println("done")
	},
}

func init() {
	ecbCmd.AddCommand(ecbCurrCmd)
}
