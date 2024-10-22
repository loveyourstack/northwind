package main

import (
	"fmt"
	"os"

	"github.com/loveyourstack/lys/lysgen"
	"github.com/spf13/cobra"
)

var genViewCmd = &cobra.Command{
	Use:   "view",
	Short: "generates PostgreSQL view from the supplied schema + table",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		res, err := lysgen.View(cmd.Context(), cliApp.Db, args[0], args[1])
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
