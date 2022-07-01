package commands

import (
	"fmt"
	"kool-dev/kool/services/updater"
	"strings"

	"github.com/blang/semver"
	"github.com/spf13/cobra"
)

// KoolSelfUpdate holds handlers and functions to implement the self-update command logic
type KoolSelfUpdate struct {
	DefaultKoolService
	updater updater.Updater
}

func AddKoolSelfUpdate(root *cobra.Command) {
	var (
		selfUpdate    = NewKoolSelfUpdate()
		selfUpdateCmd = NewSelfUpdateCommand(selfUpdate)
	)

	root.AddCommand(selfUpdateCmd)
}

// NewKoolSelfUpdate creates a new handler for self-update logic with default dependencies
func NewKoolSelfUpdate() *KoolSelfUpdate {
	return &KoolSelfUpdate{
		*newDefaultKoolService(),
		&updater.DefaultUpdater{RootCommand: rootCmd},
	}
}

// Execute runs the self-update logic with incoming arguments.
func (s *KoolSelfUpdate) Execute(args []string) (err error) {
	if err = s.updater.CheckPermission(); err != nil {
		return
	}

	var currentVersion, latestVersion semver.Version

	currentVersion = s.updater.GetCurrentVersion()

	if latestVersion, err = s.updater.Update(currentVersion); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "permission denied") {
			s.Shell().Warning("Error: %s", err)

			return fmt.Errorf("kool self-update failed: (maybe try again with sudo)")
		}
		return fmt.Errorf("kool self-update failed: %v", err)
	}

	if latestVersion.Equals(currentVersion) {
		s.Shell().Warning("You already have the latest version ", currentVersion.String())
		return
	}

	s.Shell().Success("Successfully updated to version ", latestVersion.String())
	return
}

// NewSelfUpdateCommand initializes new kool self-update command
func NewSelfUpdateCommand(selfUpdate *KoolSelfUpdate) *cobra.Command {
	selfUpdateTask := NewKoolTask("Updating kool version", selfUpdate)
	selfUpdateTask.SetFrameOutput(false)

	return &cobra.Command{
		Use:   "self-update",
		Short: "Update kool to the latest version",
		Long:  "Checks the latest release of Kool in GitHub Releases, and downloads and replaces the local binary if a newer version is available.",
		Args:  cobra.NoArgs,
		RunE:  LongTaskCommandRunFunction(selfUpdateTask),

		DisableFlagsInUseLine: true,
	}
}
