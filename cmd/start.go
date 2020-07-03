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
	handleGlobalNetwork()
	startContainers(startFlags.Services)
}

func handleGlobalNetwork() {
	networkID, err := shellExec("docker", "network", "ls", "-q", "-f", fmt.Sprintf("NAME=%s", os.Getenv("KOOL_GLOBAL_NETWORK")))

	if err != nil {
		log.Fatal(err)
	}

	if networkID != "" {
		return
	}

	_, err = shellExec("docker", "network", "create", "--attachable", os.Getenv("KOOL_GLOBAL_NETWORK"))

	if err != nil {
		log.Fatal(err)
	}
}

func startContainers(services string) {
	var (
		args []string
		err  error
		out  string
	)

	args = []string{"up", "-d", "--force-recreate"}

	if services != "" {
		args = append(args, strings.Split(services, " ")...)
	}

	out, err = shellExec("docker-compose", args...)

	if err != nil {
		execError(out, err)
		os.Exit(1)
	}
}
