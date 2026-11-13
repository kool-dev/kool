package updater

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/blang/semver"
	"github.com/google/go-github/v93/github"
	"github.com/minio/selfupdate"
)

func TestArchCandidatesFor(t *testing.T) {
	tests := []struct {
		goarch     string
		armVersion int
		want       []string
	}{
		{"amd64", 0, []string{"amd64"}},
		{"arm64", 0, []string{"arm64"}},
		{"arm", 6, []string{"arm6"}},
		{"arm", 7, []string{"arm7", "arm6"}},
		{"arm", 0, nil},
		{"arm", 5, nil},
	}

	for _, test := range tests {
		got := archCandidatesFor(test.goarch, test.armVersion)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("archCandidatesFor(%q, %d) = %v, want %v", test.goarch, test.armVersion, got, test.want)
		}
	}
}

func TestParseGOARM(t *testing.T) {
	for value, want := range map[string]int{"6": 6, "6,hardfloat": 6, "7,hardfloat": 7, "7,softfloat": 7, "5": 5, "": 0} {
		if got := parseGOARM(value); got != want {
			t.Errorf("parseGOARM(%q) = %d, want %d", value, got, want)
		}
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func githubTestClient(server *httptest.Server) *http.Client {
	base, _ := url.Parse(server.URL)
	return &http.Client{
		Timeout: time.Second,
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Host == "api.github.com" {
				req = req.Clone(req.Context())
				address := *req.URL
				address.Scheme, address.Host = base.Scheme, base.Host
				req.URL = &address
			}
			return http.DefaultTransport.RoundTrip(req)
		}),
	}
}

func TestUpdateValidatesReleaseBeforeReplacingBinary(t *testing.T) {
	newBinary := []byte("#!/bin/sh\necho updated\n")
	sum := sha256.Sum256(newBinary)
	validChecksum := fmt.Sprintf("%x  dist/kool-%s-%s\n", sum, runtime.GOOS, runtime.GOARCH)
	tests := []struct {
		name, tag, current, checksum, wantError string
		omitBinary, omitChecksum                bool
		binaryStatus, checksumStatus            int
		updated                                 bool
	}{
		{name: "upgrade", tag: "v3.8.1", current: "3.8.0", checksum: validChecksum, updated: true},
		{name: "already current", tag: "3.8.1", current: "3.8.1", omitBinary: true, omitChecksum: true},
		{name: "no downgrade", tag: "3.8.1", current: "3.9.0", omitBinary: true, omitChecksum: true},
		{name: "invalid version", tag: "invalid", current: "3.8.0", checksum: validChecksum, wantError: "invalid release version"},
		{name: "missing binary", tag: "3.8.1", current: "3.8.0", omitBinary: true, wantError: "no release asset"},
		{name: "missing checksum", tag: "3.8.1", current: "3.8.0", omitChecksum: true, wantError: "missing checksum"},
		{name: "invalid checksum", tag: "3.8.1", current: "3.8.0", checksum: "not-hex", wantError: "invalid checksum"},
		{name: "checksum mismatch", tag: "3.8.1", current: "3.8.0", checksum: strings.Repeat("0", 64), wantError: "wrong checksum"},
		{name: "checksum unavailable", tag: "3.8.1", current: "3.8.0", checksumStatus: 404, wantError: "HTTP 404"},
		{name: "binary unavailable", tag: "3.8.1", current: "3.8.0", checksum: validChecksum, binaryStatus: 500, wantError: "HTTP 500"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name := platformAssetName(runtime.GOOS, archCandidates()[0])
			var downloads atomic.Int32
			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				switch req.URL.Path {
				case "/repos/kool-dev/kool/releases/latest":
					release := &github.RepositoryRelease{TagName: test.tag}
					if !test.omitBinary {
						binaryURL := server.URL + "/binary"
						release.Assets = append(release.Assets, &github.ReleaseAsset{Name: &name, BrowserDownloadURL: &binaryURL})
					}
					if !test.omitChecksum {
						checksumName, checksumURL := name+".sha256", server.URL+"/checksum"
						release.Assets = append(release.Assets, &github.ReleaseAsset{Name: &checksumName, BrowserDownloadURL: &checksumURL})
					}
					_ = json.NewEncoder(w).Encode(release)
				case "/binary":
					downloads.Add(1)
					if test.binaryStatus != 0 {
						w.WriteHeader(test.binaryStatus)
					}
					_, _ = w.Write(newBinary)
				case "/checksum":
					downloads.Add(1)
					if test.checksumStatus != 0 {
						w.WriteHeader(test.checksumStatus)
					}
					_, _ = fmt.Fprint(w, test.checksum)
				default:
					http.NotFound(w, req)
				}
			}))
			defer server.Close()
			target := filepath.Join(t.TempDir(), "kool")
			original := []byte("original executable")
			if err := os.WriteFile(target, original, 0755); err != nil {
				t.Fatal(err)
			}
			current := semver.MustParse(test.current)
			u := &DefaultUpdater{client: githubTestClient(server)}
			version, err := u.update(context.Background(), current, selfupdate.Options{TargetPath: target})
			if test.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantError) {
					t.Fatalf("expected %q, got %v", test.wantError, err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			content, err := os.ReadFile(target)
			if err != nil {
				t.Fatal(err)
			}
			if test.updated {
				if string(content) != string(newBinary) || version.String() != "3.8.1" {
					t.Fatalf("update not installed: %q, %s", content, version)
				}
			} else if string(content) != string(original) {
				t.Fatalf("original binary changed: %q", content)
			}
			if test.wantError == "" && !test.updated && !version.Equals(current) {
				t.Fatalf("no-op changed version: %s", version)
			}
			if (test.name == "already current" || test.name == "no downgrade" || test.name == "invalid version") && downloads.Load() != 0 {
				t.Fatalf("downloaded %d assets for a non-upgrade", downloads.Load())
			}
			files, err := os.ReadDir(filepath.Dir(target))
			if err != nil || len(files) != 1 {
				t.Fatalf("left update staging files: %v (%v)", files, err)
			}
		})
	}
}

func TestUpdateTimeoutAndCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		<-req.Context().Done()
	}))
	defer server.Close()
	u := &DefaultUpdater{client: githubTestClient(server)}
	u.client.Timeout = 20 * time.Millisecond
	if _, err := u.Update(semver.MustParse("3.8.0")); err == nil {
		t.Fatal("expected stalled release request to time out")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := u.update(ctx, semver.MustParse("3.8.0"), selfupdate.Options{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
	if timeout := (&DefaultUpdater{}).httpClient().Timeout; timeout <= 0 {
		t.Fatal("default update client has no timeout")
	}
}

func TestCheckForUpdatesComparesVersionsAndClosesChannel(t *testing.T) {
	for _, tag := range []string{"3.8.1", "3.8.0", "3.7.0", "invalid"} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_ = json.NewEncoder(w).Encode(&github.RepositoryRelease{TagName: tag})
		}))
		u := &DefaultUpdater{client: githubTestClient(server)}
		ch := make(chan bool, 1)
		u.CheckForUpdates(semver.MustParse("3.8.0"), ch)
		if got := <-ch; got != (tag == "3.8.1") {
			t.Errorf("tag %s: unexpected update notice %v", tag, got)
		}
		if _, open := <-ch; open {
			t.Error("update-check channel not closed")
		}
		server.Close()
	}
}

func TestPlatformAssetName(t *testing.T) {
	if got := platformAssetName("linux", "amd64"); got != "kool-linux-amd64" {
		t.Errorf("unexpected asset name: %s", got)
	}
	if got := platformAssetName("linux", "arm6"); got != "kool-linux-arm6" {
		t.Errorf("unexpected asset name: %s", got)
	}
}

func TestParseVersion(t *testing.T) {
	for _, tag := range []string{"3.8.1", "v3.8.1"} {
		got, err := parseVersion(tag)
		if err != nil {
			t.Fatalf("parseVersion(%q): %v", tag, err)
		}
		if got.String() != "3.8.1" {
			t.Errorf("parseVersion(%q) = %s, want 3.8.1", tag, got)
		}
	}

	if _, err := parseVersion("not-a-version"); err == nil {
		t.Error("expected an error parsing a non-semver tag")
	}
}

func TestFindAsset(t *testing.T) {
	name := "kool-linux-amd64"
	release := &github.RepositoryRelease{
		Assets: []*github.ReleaseAsset{
			{Name: &name},
		},
	}

	if asset := findAsset(release, name); asset == nil {
		t.Fatal("expected to find the release asset")
	}
	if asset := findAsset(release, "kool-linux-arm64"); asset != nil {
		t.Errorf("expected no asset, got %q", asset.GetName())
	}
}

func TestParseChecksum(t *testing.T) {
	sum := sha256.Sum256([]byte("kool"))
	valid := hex.EncodeToString(sum[:]) + "  kool-linux-amd64\n"

	got, err := parseChecksum([]byte(valid))
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(got) != hex.EncodeToString(sum[:]) {
		t.Errorf("unexpected checksum: %x", got)
	}

	for _, input := range []string{"", "   \n", "not-hex  kool-linux-amd64", "abcd  kool-linux-amd64"} {
		if _, err := parseChecksum([]byte(input)); err == nil || !strings.Contains(err.Error(), "checksum") {
			t.Errorf("parseChecksum(%q) expected checksum error, got %v", input, err)
		}
	}
}
