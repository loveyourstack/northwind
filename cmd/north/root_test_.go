package main

import (
	"fmt"
	"os"

	"github.com/loveyourstack/northwind/pkg/ecbdata"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "General test cmd",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		defer cliApp.Db.Close()

		currs, err := ecbdata.GetCurrencies()
		if err != nil {
			cliApp.ErrorLog.Error("ecbdata.GetCurrencies failed: " + err.Error())
			os.Exit(1)
		}
		fmt.Println(currs)

		/*exRates, err := ecbdata.GetExchangeRates("EUR", ecbdata.Daily, time.Now().Add(-7*24*time.Hour), time.Now())
		if err != nil {
			cliApp.ErrorLog.Error("ecbdata.GetExchangeRates failed: " + err.Error())
			os.Exit(1)
		}
		fmt.Println(exRates)*/

		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
