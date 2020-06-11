package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// StartFlags holds the flags for the start command
type StartFlags struct {
	All      bool
	Services string
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start fwd environment containers",
	Run:   runStart,
}

var startFlags = &StartFlags{false, ""}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolVarP(&startFlags.All, "all", "a", false, "Start all services")
	startCmd.Flags().StringVarP(&startFlags.Services, "services", "", "", "Specific services to be started")
}

func runStart(cmd *cobra.Command, args []string) {
	handleGlobalNetwork()
	startContainers(startFlags.All, startFlags.Services)
}

func handleGlobalNetwork() {
	networkID, err := shellExec("docker", "network", "ls", "-q", "-f", fmt.Sprintf("NAME=%s", os.Getenv("FWD_NETWORK")))

	if err != nil {
		log.Fatal(err)
	}

	if networkID != "" {
		return
	}

	_, err = shellExec("docker", "network", "create", "--attachable", os.Getenv("FWD_NETWORK"))

	if err != nil {
		log.Fatal(err)
	}
}

func startContainers(all bool, services string) {
	var (
		args []string
		err  error
		out  string
	)

	args = []string{"up", "-d", "--force-recreate"}

	if !all {
		if services == "" {
			services = os.Getenv("FWD_START_DEFAULT_SERVICES")
		}

		args = append(args, strings.Split(services, " ")...)
	}

	out, err = shellExec("docker-compose", args...)

	if err != nil {
		execError(out, err)
		os.Exit(1)
	}
}
