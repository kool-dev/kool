package updater

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"

	"kool-dev/kool/services/user"

	"github.com/blang/semver"
	"github.com/google/go-github/v93/github"
	"github.com/minio/selfupdate"
	"github.com/spf13/cobra"
)

const (
	updateOwner = "kool-dev"
	updateRepo  = "kool"
	updateCmd   = "kool"
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
	ctx := context.Background()

	release, err := latestRelease(ctx)
	if err != nil {
		return
	}

	var (
		assetName string
		asset     *github.ReleaseAsset
	)
	for _, arch := range archCandidates() {
		assetName = platformAssetName(runtime.GOOS, arch)
		if asset = findAsset(release, assetName); asset != nil {
			break
		}
	}
	if asset == nil {
		err = fmt.Errorf("no release asset found for %s/%s", runtime.GOOS, runtime.GOARCH)
		return
	}

	sum, err := checksum(ctx, release, assetName)
	if err != nil {
		return
	}

	binary, err := download(ctx, asset.GetBrowserDownloadURL())
	if err != nil {
		return
	}

	if err = selfupdate.Apply(bytes.NewReader(binary), selfupdate.Options{Checksum: sum}); err != nil {
		err = fmt.Errorf("applying update: %w", err)
		return
	}

	updatedVersion, err = parseVersion(release.GetTagName())
	return
}

// CheckForUpdates checks if there is a new version
func (u *DefaultUpdater) CheckForUpdates(current semver.Version, chHasNewVersion chan bool) {
	defer close(chHasNewVersion)

	release, err := latestRelease(context.Background())
	if err != nil {
		chHasNewVersion <- false
		return
	}

	latest, err := parseVersion(release.GetTagName())
	if err != nil {
		chHasNewVersion <- false
		return
	}

	chHasNewVersion <- latest.GT(current)
}

// CheckPermission will return an error if the running
// user has not enough privileges to perform this task,
// OS wise.
func (u *DefaultUpdater) CheckPermission() (err error) {
	var (
		binPath string
	)

	if binPath, err = os.Executable(); err != nil {
		return
	}

	if binPath, err = filepath.EvalSymlinks(binPath); err != nil {
		return
	}

	if isWriteable(binPath) && isWriteable(filepath.Dir(binPath)) {
		// the folder the binary file lives is IS writeable
		// for the current user
		return
	}

	// we need elevated privileges!
	var isAdmin = user.CurrentUserIsElevated()
	if !isAdmin {
		err = fmt.Errorf("you need to use 'sudo' to perform this task")
	}

	return
}

func latestRelease(ctx context.Context) (*github.RepositoryRelease, error) {
	client, err := github.NewClient()
	if err != nil {
		return nil, fmt.Errorf("creating GitHub client: %w", err)
	}

	release, _, err := client.Repositories.GetLatestRelease(ctx, updateOwner, updateRepo)
	if err != nil {
		return nil, fmt.Errorf("fetching latest release: %w", err)
	}
	return release, nil
}

func parseVersion(tag string) (semver.Version, error) {
	return semver.Parse(strings.TrimPrefix(tag, "v"))
}

func platformAssetName(goos, arch string) string {
	return fmt.Sprintf("%s-%s-%s", updateCmd, goos, arch)
}

func findAsset(release *github.RepositoryRelease, name string) *github.ReleaseAsset {
	for _, asset := range release.Assets {
		if asset.GetName() == name {
			return asset
		}
	}
	return nil
}

func checksum(ctx context.Context, release *github.RepositoryRelease, assetName string) ([]byte, error) {
	name := assetName + ".sha256"

	asset := findAsset(release, name)
	if asset == nil {
		return nil, fmt.Errorf("missing checksum asset %s", name)
	}

	data, err := download(ctx, asset.GetBrowserDownloadURL())
	if err != nil {
		return nil, err
	}

	sum, err := parseChecksum(data)
	if err != nil {
		return nil, fmt.Errorf("invalid checksum asset %s: %w", name, err)
	}
	return sum, nil
}

func parseChecksum(data []byte) ([]byte, error) {
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty checksum")
	}

	sum, err := hex.DecodeString(fields[0])
	if err != nil || len(sum) != sha256.Size {
		return nil, fmt.Errorf("invalid checksum")
	}
	return sum, nil
}

func download(ctx context.Context, url string) (data []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer func() {
		if closeErr := resp.Body.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("downloading %s: HTTP %d", url, resp.StatusCode)
	}

	data, err = io.ReadAll(resp.Body)
	return
}

func archCandidates() []string {
	return archCandidatesFor(runtime.GOARCH, goarm())
}

func archCandidatesFor(goarch string, armVersion int) []string {
	if goarch != "arm" {
		return []string{goarch}
	}

	if armVersion == 6 {
		return []string{"arm6", "arm7"}
	}
	return []string{"arm7", "arm6"}
}

func goarm() int {
	if v, err := strconv.Atoi(os.Getenv("GOARM")); err == nil && (v == 6 || v == 7) {
		return v
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "GOARM" {
				if v, err := strconv.Atoi(setting.Value); err == nil {
					return v
				}
			}
		}
	}

	return 7
}
