package main

import (
	"fmt"
	"os"

	"github.com/loveyourstack/lys/lysgen"
	"github.com/spf13/cobra"
)

var genImCmd = &cobra.Command{
	Use:   "im",
	Short: "generates Go Input and Model structs from the supplied schema + table",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		res, err := lysgen.InputModel(cmd.Context(), cliApp.Db, args[0], args[1])
		if err != nil {
			cliApp.ErrorLog.Error("lysgen.InputModel failed: " + err.Error())
			os.Exit(1)
		}

		fmt.Println(res)
	},
}

func init() {
	genCmd.AddCommand(genImCmd)
}
