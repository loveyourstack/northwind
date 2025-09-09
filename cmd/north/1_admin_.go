package main

import (
	"github.com/spf13/cobra"
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "TODO",
}

func init() {
	rootCmd.AddCommand(adminCmd)
}
