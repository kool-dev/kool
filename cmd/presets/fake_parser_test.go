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
	fileError, err := f.WriteFile("filename", "filecontent")

	if !f.CalledWriteFile || fileError != f.MockFileError || f.MockError.Error() != err.Error() {
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

	f.MockPresetKeys = []string{"key"}
	keys := f.GetPresetKeys("preset")

	if !f.CalledGetPresetKeys || len(keys) != 1 || keys[0] != "key" {
		t.Error("failed to use mocked GetPresetKeys function on FakeParser")
	}

	f.MockPresetKeyContent = "content"
	content := f.GetPresetKeyContent("preset", "key")

	if !f.CalledGetPresetKeyContent || content != "content" {
		t.Error("failed to use mocked GetPresetKeyContent function on FakeParser")
	}
}
