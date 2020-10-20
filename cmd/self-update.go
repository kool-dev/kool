package cmd

import (
	"fmt"
	"kool-dev/kool/cmd/task"
	"kool-dev/kool/cmd/updater"
	"strings"

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
	newVersion, err := s.taskRunner.Run("Updating kool version", func() (taskResult interface{}, taskError error) {
		var currentVersion, latestVersion semver.Version
		currentVersion = s.updater.GetCurrentVersion()

		if latestVersion, taskError = s.updater.Update(currentVersion); taskError != nil {
			return
		}

		if latestVersion.Equals(currentVersion) {
			taskError = fmt.Errorf("You already have the latest version %s", currentVersion.String())
		}

		taskResult = latestVersion

		return
	})

	if err != nil {
		if strings.Contains(err.Error(), "You already have the latest version") {
			s.Warning(err.Error())
			return nil
		}

		return fmt.Errorf("kool self-update failed: %v", err)
	}

	s.Success("Successfully updated to version ", newVersion.(semver.Version).String())
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
