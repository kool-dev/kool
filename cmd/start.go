package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start fwd environment containers",
	Run:   runStart,
}

// StartFlags holds the flags for the start command
type StartFlags struct {
	All      bool
	Services string
}

var startFlags = &StartFlags{false, ""}

func init() {
	startCmd.Flags().BoolVarP(&startFlags.All, "all", "a", false, "Start all services")
	startCmd.Flags().StringVarP(&startFlags.Services, "services", "", "", "Specific services to be started")
}

func runStart(cmf *cobra.Command, args []string) {
	if startFlags.Services == "" {
		// we should default to the environment variable
		startFlags.Services = os.Getenv("FWD_START_DEFAULT_SERVICES")
	}

	startGlobalNetwork()
	startContainers()
}

func startGlobalNetwork() {
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

func startContainers() {
	var (
		args []string
		err  error
		out  string
	)

	args = []string{"up", "-d", "--force-recreate"}

	if !startFlags.All {
		if startFlags.Services == "" {
			startFlags.Services = os.Getenv("FWD_START_DEFAULT_SERVICES")
		}

		args = append(args, strings.Split(startFlags.Services, " ")...)
	}

	out, err = shellExec("docker-compose", args...)

	if err != nil {
		log.Println("ERROR: ", err)
		log.Println("Output:")
		fmt.Println(out)
	}
}
