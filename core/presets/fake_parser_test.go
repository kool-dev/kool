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

	f.MockGetPresets = map[string]string{"preset": "preset"}
	presets := f.GetPresets("")

	if !f.CalledGetPresets || len(presets) != 1 || presets["preset"] != "preset" {
		t.Error("failed to use mocked GetPresets function on FakeParser")
	}

	f.MockGetTags = []string{"php"}
	tags := f.GetTags()

	if !f.CalledGetTags || len(tags) != 1 || tags[0] != "php" {
		t.Error("failed to use mocked GetTags function on FakeParser")
	}

	f.MockInstall = errors.New("Install")
	errInstall := f.Install("")

	if !f.CalledInstall || errInstall == nil || errInstall.Error() != "Install" {
		t.Error("failed to use mocked Install function on FakeParser")
	}

	f.MockCreate = errors.New("Create")
	errCreate := f.Create("")

	if !f.CalledCreate || errCreate == nil || errCreate.Error() != "Create" {
		t.Error("failed to use mocked Create function on FakeParser")
	}

	f.MockAdd = errors.New("Add")
	errAdd := f.Add("", nil)

	if !f.CalledAdd || errAdd == nil || errAdd.Error() != "Add" {
		t.Error("failed to use mocked Add function on FakeParser")
	}
}
