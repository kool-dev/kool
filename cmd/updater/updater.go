package updater

import (
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
