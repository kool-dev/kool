package presets

import (
	"errors"
	"testing"
)

func TestFakeParser(t *testing.T) {
	f := &FakeParser{}

	f.MockExists = true
	exists := f.Exists("preset")

	if !f.CalledExists || exists != f.MockExists {
		t.Error("failed to use mocked Exists function on FakeParser")
	}

	f.MockFoundFiles = []string{"kool.yml"}
	foundFiles := f.LookUpFiles("preset")

	if !f.CalledLookUpFiles || len(foundFiles) != 1 || foundFiles[0] != "kool.yml" {
		t.Error("failed to use mocked LookUpFiles function on FakeParser")
	}

	f.MockFileError = "kool.yml"
	f.MockError = errors.New("error")
	fileError, err := f.WriteFiles("preset")

	if !f.CalledWriteFiles || fileError != f.MockFileError || f.MockError.Error() != err.Error() {
		t.Error("failed to use mocked WriteFiles function on FakeParser")
	}
}
