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

	f.MockPresets = []string{"preset"}
	presets := f.GetPresets("")

	if !f.CalledGetPresets || len(presets) != 1 || presets[0] != "preset" {
		t.Error("failed to use mocked GetPresets function on FakeParser")
	}

	f.MockLanguages = []string{"php"}
	languages := f.GetLanguages()

	if !f.CalledGetLanguages || len(languages) != 1 || languages[0] != "php" {
		t.Error("failed to use mocked GetPresets function on FakeParser")
	}

	f.MockCreateCommand = "create"
	createCommand, _ := f.GetCreateCommand("")

	if !f.CalledGetCreateCommand || createCommand == "" || createCommand != "create" {
		t.Error("failed to use mocked GetCreateCommand function on FakeParser")
	}
}
