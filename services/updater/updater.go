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
	"time"

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
	client      *http.Client
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
	return u.update(context.Background(), currentVersion, selfupdate.Options{})
}

func (u *DefaultUpdater) update(ctx context.Context, currentVersion semver.Version, options selfupdate.Options) (updatedVersion semver.Version, err error) {
	client := u.httpClient()
	release, err := latestRelease(ctx, client)
	if err != nil {
		return
	}
	updatedVersion, err = parseVersion(release.GetTagName())
	if err != nil {
		return currentVersion, fmt.Errorf("invalid release version: %w", err)
	}
	if !updatedVersion.GT(currentVersion) {
		return currentVersion, nil
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

	sum, err := checksum(ctx, client, release, assetName)
	if err != nil {
		return
	}

	binary, err := download(ctx, client, asset.GetBrowserDownloadURL())
	if err != nil {
		return
	}

	options.Checksum = sum
	if err = selfupdate.Apply(bytes.NewReader(binary), options); err != nil {
		if rollbackErr := selfupdate.RollbackError(err); rollbackErr != nil {
			return currentVersion, fmt.Errorf("applying update: %w; rollback failed: %v; restore the binary manually", err, rollbackErr)
		}
		err = fmt.Errorf("applying update: %w", err)
		return
	}

	return
}

// CheckForUpdates checks if there is a new version
func (u *DefaultUpdater) CheckForUpdates(current semver.Version, chHasNewVersion chan bool) {
	defer close(chHasNewVersion)

	release, err := latestRelease(context.Background(), u.httpClient())
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

func (u *DefaultUpdater) httpClient() *http.Client {
	if u.client != nil {
		return u.client
	}
	return &http.Client{Timeout: 2 * time.Minute}
}

func latestRelease(ctx context.Context, httpClient *http.Client) (*github.RepositoryRelease, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	client, err := github.NewClient(github.WithHTTPClient(httpClient))
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

func checksum(ctx context.Context, client *http.Client, release *github.RepositoryRelease, assetName string) ([]byte, error) {
	name := assetName + ".sha256"

	asset := findAsset(release, name)
	if asset == nil {
		return nil, fmt.Errorf("missing checksum asset %s", name)
	}

	data, err := download(ctx, client, asset.GetBrowserDownloadURL())
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

func download(ctx context.Context, client *http.Client, url string) (data []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
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
		return []string{"arm6"}
	}
	if armVersion == 7 {
		return []string{"arm7", "arm6"}
	}
	return nil
}

func goarm() int {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "GOARM" {
				return parseGOARM(setting.Value)
			}
		}
	}

	return 0
}

func parseGOARM(value string) int {
	version, _, _ := strings.Cut(value, ",")
	v, _ := strconv.Atoi(version)
	return v
}
