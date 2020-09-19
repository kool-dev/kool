package cmd

import (
	"kool-dev/kool/cmd/shell"
	"os"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(selfUpdateCmd)
}

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update kool to latest version",
	Long:  "Checks for the latest release of Kool on Github Releases, downloads and replaces the local binary if a newer version is available.",
	Run:   runSelfUpdate,
	Args:  cobra.MaximumNArgs(0),
}

var selfUpdateCmdOutputWriter shell.OutputWriter = shell.NewOutputWriter()

func runSelfUpdate(cmf *cobra.Command, args []string) {
	var (
		currentVersion semver.Version
		updater        *selfupdate.Updater
		err            error
		latest         *selfupdate.Release
	)

	selfUpdateCmdOutputWriter.SetWriter(cmf.OutOrStdout())

	currentVersion = semver.MustParse(version)
	if updater, err = selfupdate.NewUpdater(selfupdate.Config{
		Validator: &selfupdate.SHA2Validator{},
	}); err != nil {
		selfUpdateCmdOutputWriter.Error("Failed to create updater controller", err)
		os.Exit(1)
	}

	if latest, err = updater.UpdateSelf(currentVersion, "kool-dev/kool"); err != nil {
		selfUpdateCmdOutputWriter.Error("Kool self-update failed:", err)
		os.Exit(1)
	}

	if latest.Version.Equals(currentVersion) {
		selfUpdateCmdOutputWriter.Warning("You already have the latest version.", version)
	} else {
		selfUpdateCmdOutputWriter.Success("Successfully updated to version", latest.Version)
	}
}
