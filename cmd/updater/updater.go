package updater

import (
	"fmt"
	"kool-dev/kool/cmd/user"
	"runtime"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

// DefaultUpdater holds data for updating kool
type DefaultUpdater struct {
	RootCommand *cobra.Command
}

// Updater holds logic for updating kool
type Updater interface {
	GetCurrentVersion() semver.Version
	Update(semver.Version) (semver.Version, error)
	CheckForUpdates(semver.Version, chan bool)
	CheckPermission() error
}

// GetCurrentVersion get current version
func (u *DefaultUpdater) GetCurrentVersion() semver.Version {
	return semver.MustParse(u.RootCommand.Version)
}

// Update updates kool version
func (u *DefaultUpdater) Update(currentVersion semver.Version) (updatedVersion semver.Version, err error) {
	var (
		updater *selfupdate.Updater
		latest  *selfupdate.Release
	)

	if updater, err = selfupdate.NewUpdater(selfupdate.Config{
		Validator: &selfupdate.SHA2Validator{},
	}); err != nil {
		return
	}

	if latest, err = updater.UpdateSelf(currentVersion, "kool-dev/kool"); err != nil {
		return
	}

	updatedVersion = latest.Version
	return
}

// CheckForUpdates checks if there is a new version
func (u *DefaultUpdater) CheckForUpdates(current semver.Version, chHasNewVersion chan bool) {
	var (
		latest *selfupdate.Release
		err    error
	)

	if latest, _, err = selfupdate.DetectLatest("kool-dev/kool"); err != nil {
		chHasNewVersion <- false
		return
	}

	if !latest.Version.Equals(current) {
		chHasNewVersion <- true
	}

	close(chHasNewVersion)
}

// CheckPermission will return an error if the running
// user has not enough privileges to perform this task,
// OS wise.
func (u *DefaultUpdater) CheckPermission() (err error) {
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
