package updater

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-github/v93/github"
)

func TestArchCandidatesFor(t *testing.T) {
	tests := []struct {
		goarch     string
		armVersion int
		want       []string
	}{
		{"amd64", 0, []string{"amd64"}},
		{"arm64", 0, []string{"arm64"}},
		{"arm", 6, []string{"arm6", "arm7"}},
		{"arm", 7, []string{"arm7", "arm6"}},
		{"arm", 0, []string{"arm7", "arm6"}},
	}

	for _, test := range tests {
		got := archCandidatesFor(test.goarch, test.armVersion)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("archCandidatesFor(%q, %d) = %v, want %v", test.goarch, test.armVersion, got, test.want)
		}
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
