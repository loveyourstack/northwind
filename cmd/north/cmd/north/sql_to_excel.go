package north

import (
	"os"

	"github.com/loveyourstack/lys/lysexcel"
	"github.com/loveyourstack/lys/lyspg"
	"github.com/loveyourstack/northwind/internal/stores/core/corecategory"
	"github.com/spf13/cobra"
)

var sqlToExcelCmd = &cobra.Command{
	Use:   "sqlToExcel",
	Short: "cmd for testing Excel file output",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		defer cliApp.Db.Close()

		// select from db
		catStore := corecategory.Store{Db: cliApp.Db}
		items, _, _, err := catStore.Select(cmd.Context(), lyspg.SelectParams{
			Fields: catStore.GetJsonFields(),
		})
		if err != nil {
			cliApp.ErrorLog.Error("catStore.Select failed: " + err.Error())
			os.Exit(1)
		}

		// create temp file
		f, err := os.CreateTemp("", "categories.*.xlsx")
		if err != nil {
			cliApp.ErrorLog.Error("os.CreateTemp failed: " + err.Error())
			os.Exit(1)
		}
		f.Close()
		defer os.Remove(f.Name())

		// write items to Excel
		err = lysexcel.WriteItemsToFile(items, catStore.GetJsonTagTypeMap(), f.Name(), "")
		if err != nil {
			cliApp.ErrorLog.Error("lysexcel.WriteItemsToFile failed: " + err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(sqlToExcelCmd)
}
