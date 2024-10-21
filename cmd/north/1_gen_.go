package main

import (
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "code generation",
}

func init() {
	rootCmd.AddCommand(genCmd)
}
