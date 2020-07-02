package cmd

import (
	"github.com/spf13/cobra"
)

// DbFlags holds the flags for the db command
type DbFlags struct {
	ServiceName string
}

var dbFlags = &DbFlags{""}

var dbCmd = &cobra.Command{
	Use:              "db",
	Short:            "Useful database service related actions",
	TraverseChildren: true,
}

func init() {
	rootCmd.AddCommand(dbCmd)

	dbCmd.PersistentFlags().StringVarP(&dbFlags.ServiceName, "service", "s", "database", "The service name for the database container.")
}
