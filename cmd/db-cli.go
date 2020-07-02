package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// DbCliFlags holds the flags for the db cli command
type DbCliFlags struct {
}

var dbCliCmd = &cobra.Command{
	Use:              "cli",
	Short:            "Opens a CLI default client session within the main database service in use, if one exists",
	Run:              runDbCli,
	TraverseChildren: true,
}

var dbCliFlags = &DbCliFlags{}

func init() {
	dbCmd.AddCommand(dbCliCmd)

	// dbCliCmd.Flags().BoolVarP(&dbCliFlags.Purge, "purge", "", false, "Remove all persistent data from containers")
}

func runDbCli(cmd *cobra.Command, args []string) {
	dbCliOpen(args...)
}

func dbCliOpen(extraArgs ...string) {
	var (
		args []string
		err  error
	)

	fmt.Println("dbFlags.ServiceName", dbFlags.ServiceName)
	args = []string{"exec", "-e", "MYSQL_PWD=" + os.Getenv("DB_PASSWORD"), dbFlags.ServiceName, "mysql", "-u", os.Getenv("DB_USERNAME")}
	args = append(args, os.Getenv("DB_DATABASE"))
	args = append(args, extraArgs...)

	err = shellInteractive("docker-compose", args...)

	if err != nil {
		execError("", err)
		os.Exit(1)
	}
}
