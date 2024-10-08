package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "General test cmd",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		defer cliApp.Db.Close()

		/*currs, err := ecbClient.GetApiCurrencies()
		if err != nil {
			cliApp.ErrorLog.Error("ecbClient.GetApiCurrencies failed: " + err.Error())
			os.Exit(1)
		}
		fmt.Println(currs)*/

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
	rootCmd.AddCommand(testCmd)
}
