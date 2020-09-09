package cmd

import (
	"fmt"
	"kool-dev/kool/cmd/shell"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// StartFlags holds the flags for the start command
type StartFlags struct {
	Services string
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Kool environment containers",
	Run:   runStart,
}

var startFlags = &StartFlags{""}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&startFlags.Services, "services", "", "", "Specific services to be started")
}

func runStart(cmd *cobra.Command, args []string) {
	checkKoolDependencies()
	handleGlobalNetwork()
	startContainers(startFlags.Services)
}

func handleGlobalNetwork() {
	networkID, err := shell.Exec("docker", "network", "ls", "-q", "-f", fmt.Sprintf("NAME=^%s$", os.Getenv("KOOL_GLOBAL_NETWORK")))

	if err != nil {
		log.Fatal(err)
	}

	if networkID != "" {
		return
	}

	err = shell.Interactive("docker", "network", "create", "--attachable", os.Getenv("KOOL_GLOBAL_NETWORK"))

	if err != nil {
		log.Fatal(err)
	}
}

func startContainers(services string) {
	var (
		args []string
		err  error
	)

	args = []string{"up", "-d", "--force-recreate"}

	if services != "" {
		args = append(args, strings.Split(services, " ")...)
	}

	err = shell.Interactive("docker-compose", args...)

	if err != nil {
		execError("", err)
		os.Exit(1)
	}
}
