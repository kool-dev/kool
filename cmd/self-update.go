package cmd

import (
	"fmt"
	"kool-dev/kool/cmd/task"
	"kool-dev/kool/cmd/updater"

	"github.com/blang/semver"
	"github.com/spf13/cobra"
)

// KoolSelfUpdate holds handlers and functions to implement the self-update command logic
type KoolSelfUpdate struct {
	DefaultKoolService
	taskRunner task.Runner
	updater    updater.Updater
}

func init() {
	var (
		selfUpdate    = NewKoolSelfUpdate()
		selfUpdateCmd = NewSelfUpdateCommand(selfUpdate)
	)

	rootCmd.AddCommand(selfUpdateCmd)
}

// NewKoolSelfUpdate creates a new handler for self-update logic with default dependencies
func NewKoolSelfUpdate() *KoolSelfUpdate {
	return &KoolSelfUpdate{
		*newDefaultKoolService(),
		task.NewRunner(),
		&updater.DefaultUpdater{RootCommand: rootCmd},
	}
}

// Execute runs the self-update logic with incoming arguments.
func (s *KoolSelfUpdate) Execute(args []string) (err error) {
	var currentVersion, latestVersion semver.Version

	err = s.taskRunner.Run("Updating kool version", func() (taskError error) {
		currentVersion = s.updater.GetCurrentVersion()
		latestVersion, taskError = s.updater.Update(currentVersion)
		return
	})

	if err != nil {
		return fmt.Errorf("kool self-update failed: %v", err)
	}

	if latestVersion.Equals(currentVersion) {
		s.Warning("You already have the latest version ", currentVersion.String())
	} else {
		s.Success("Successfully updated to version ", latestVersion.String())
	}

	return
}

// NewSelfUpdateCommand initializes new kool self-update command
func NewSelfUpdateCommand(selfUpdate *KoolSelfUpdate) *cobra.Command {
	return &cobra.Command{
		Use:   "self-update",
		Short: "Update kool to latest version",
		Long:  "Checks for the latest release of Kool on Github Releases, downloads and replaces the local binary if a newer version is available.",
		Args:  cobra.NoArgs,
		Run:   DefaultCommandRunFunction(selfUpdate),
	}
}
