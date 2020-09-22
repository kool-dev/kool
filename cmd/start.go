package cmd

import (
	"kool-dev/kool/cmd/checker"
	"kool-dev/kool/cmd/network"
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
	var dependenciesChecker = checker.NewChecker()

	if err := dependenciesChecker.VerifyDependencies(); err != nil {
		shell.ExecError("", err)
		os.Exit(1)
	}

	var globalNetworkHandler = network.NewHandler()

	if err := globalNetworkHandler.HandleGlobalNetwork(os.Getenv("KOOL_GLOBAL_NETWORK")); err != nil {
		log.Fatal(err)
	}

	startContainers(startFlags.Services)
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
		shell.ExecError("", err)
		os.Exit(1)
	}
}
