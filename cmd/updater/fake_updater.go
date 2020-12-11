package updater

import (
	"github.com/blang/semver"
)

// FakeUpdater implements all fake behaviors for self-update
type FakeUpdater struct {
	CalledGetCurrentVersion, CalledUpdate, CalledCheckForUpdates bool
	MockCurrentVersion, MockLatestVersion                        string
	MockError                                                    error
}

// GetCurrentVersion get mocked current version
func (u *FakeUpdater) GetCurrentVersion() semver.Version {
	u.CalledGetCurrentVersion = true
	return semver.MustParse(u.MockCurrentVersion)
}

// Update implements fake update
func (u *FakeUpdater) Update(currentVersion semver.Version) (updatedVersion semver.Version, err error) {
	updatedVersion = semver.MustParse(u.MockLatestVersion)
	err = u.MockError
	u.CalledUpdate = true
	return
}

// CheckForUpdates implements fake check
func (u *FakeUpdater) CheckForUpdates(currentVersion semver.Version, ch chan bool) {
	u.CalledCheckForUpdates = true
	ch <- true
}
