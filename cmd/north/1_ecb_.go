package main

import (
	"github.com/spf13/cobra"
)

var ecbCmd = &cobra.Command{
	Use:   "ecb",
	Short: "Commands related to the ECB connector",
}

func init() {
	rootCmd.AddCommand(ecbCmd)
}
