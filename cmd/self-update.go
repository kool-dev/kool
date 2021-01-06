package cmd

import (
	"fmt"
	"kool-dev/kool/cmd/updater"
	"kool-dev/kool/cmd/user"
	"runtime"

	"github.com/blang/semver"
	"github.com/spf13/cobra"
)

// KoolSelfUpdate holds handlers and functions to implement the self-update command logic
type KoolSelfUpdate struct {
	DefaultKoolService
	updater updater.Updater
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
		&updater.DefaultUpdater{RootCommand: rootCmd},
	}
}

// Execute runs the self-update logic with incoming arguments.
func (s *KoolSelfUpdate) Execute(args []string) (err error) {
	if err = s.checkPermission(); err != nil {
		return
	}

	var currentVersion, latestVersion semver.Version

	currentVersion = s.updater.GetCurrentVersion()

	if latestVersion, err = s.updater.Update(currentVersion); err != nil {
		return fmt.Errorf("kool self-update failed: %v", err)
	}

	if latestVersion.Equals(currentVersion) {
		s.Warning("You already have the latest version ", currentVersion.String())
		return
	}

	s.Success("Successfully updated to version ", latestVersion.String())
	return
}

func (s *KoolSelfUpdate) checkPermission() (err error) {
	if runtime.GOOS != "windows" && runtime.GOOS != "linux" {
		// we should be fine in other plataforms, permission-wise
		return
	}

	// we need elevated privileges!
	var isAdmin = user.CurrentUserIsElevated()
	if !isAdmin {
		if runtime.GOOS == "linux" {
			err = fmt.Errorf("you need to use 'sudo' to perform this task")
		} else if runtime.GOOS == "windows" {
			err = fmt.Errorf("you need to Run as Administrator to perform this task")
		}
	}

	return
}

// NewSelfUpdateCommand initializes new kool self-update command
func NewSelfUpdateCommand(selfUpdate *KoolSelfUpdate) *cobra.Command {
	selfUpdateTask := NewKoolTask("Updating kool version", selfUpdate)

	return &cobra.Command{
		Use:   "self-update",
		Short: "Update kool to latest version",
		Long:  "Checks for the latest release of Kool on Github Releases, downloads and replaces the local binary if a newer version is available.",
		Args:  cobra.NoArgs,
		Run:   LongTaskCommandRunFunction(selfUpdateTask),
	}
}
