package northgen

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/loveyourstack/lys/lysgen"
	"github.com/spf13/cobra"
)

var genInputModelCmd = &cobra.Command{
	Use:   "im",
	Short: "im - generate Go Input and Model structs from the supplied schema.table",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		defer cliApp.Db.Close()

		argA := strings.Split(args[0], ".")

		res, stmt, err := lysgen.InputModel(context.Background(), cliApp.Db, argA[0], argA[1])
		if err != nil {
			cliApp.ErrorLog.Error("lysgen.InputModel failed: "+err.Error(), slog.String("stmt", stmt))
			os.Exit(1)
		}

		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(genInputModelCmd)
}
