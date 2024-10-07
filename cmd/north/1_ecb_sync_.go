package main

import (
	"github.com/spf13/cobra"
)

var ecbSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync API data into database",
}

func init() {
	ecbCmd.AddCommand(ecbSyncCmd)
}
