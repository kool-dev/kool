package updater

import (
	"github.com/blang/semver"
)

// FakeUpdater implements all fake behaviors for self-update
type FakeUpdater struct {
	CalledGetCurrentVersion, CalledUpdate,
	CalledCheckForUpdates, CalledCheckPermission bool

	MockCurrentVersion, MockLatestVersion string
	MockErrorUpdate, MockErrorPermission  error
}

// GetCurrentVersion get mocked current version
func (u *FakeUpdater) GetCurrentVersion() semver.Version {
	u.CalledGetCurrentVersion = true
	return semver.MustParse(u.MockCurrentVersion)
}

// Update implements fake update
func (u *FakeUpdater) Update(currentVersion semver.Version) (updatedVersion semver.Version, err error) {
	updatedVersion = semver.MustParse(u.MockLatestVersion)
	err = u.MockErrorUpdate
	u.CalledUpdate = true
	return
}

// CheckForUpdates implements fake available update check
func (u *FakeUpdater) CheckForUpdates(currentVersion semver.Version, ch chan bool) {
	u.CalledCheckForUpdates = true
	ch <- true
}

// CheckPermission implements fake permission check
func (u *FakeUpdater) CheckPermission() (err error) {
	u.CalledCheckPermission = true
	err = u.MockErrorPermission
	return
}
