package presets

import (
	"os"
	"testing"
)

func TestExistsParser(t *testing.T) {
	presets := make(map[string]map[string]string)
	laravelPreset := make(map[string]string)

	laravelPreset["kool.yml"] = ""
	presets["laravel"] = laravelPreset

	p := &DefaultParser{
		Presets: presets,
	}

	exists := p.Exists("laravel")

	if !exists {
		t.Error("did not find existing preset laravel")
	}

	exists = p.Exists("javascript")

	if exists {
		t.Error("found not existing preset javascript")
	}
}

func TestGetAllParser(t *testing.T) {
	presets := make(map[string]map[string]string)
	laravelPreset := make(map[string]string)
	symfonyPreset := make(map[string]string)

	laravelPreset["kool.yml"] = ""
	symfonyPreset["kool.yml"] = ""

	presets["laravel"] = laravelPreset
	presets["symfony"] = symfonyPreset

	p := &DefaultParser{
		Presets: presets,
	}

	allPresets := p.GetPresets("")

	if len(allPresets) != 2 || allPresets[0] != "laravel" || allPresets[1] != "symfony" {
		t.Error("failed to get all presets")
	}
}

func TestGetLanguagesParser(t *testing.T) {
	presets := make(map[string]map[string]string)

	laravelPreset := make(map[string]string)
	symfonyPreset := make(map[string]string)

	laravelPreset["preset_language"] = "php"
	laravelPreset["kool.yml"] = ""
	symfonyPreset["preset_language"] = "php"
	symfonyPreset["kool.yml"] = ""

	presets["laravel"] = laravelPreset
	presets["symfony"] = symfonyPreset

	p := &DefaultParser{
		Presets: presets,
	}

	allLanguages := p.GetLanguages()

	if len(allLanguages) != 1 || allLanguages[0] != "php" {
		t.Error("failed to get languages")
	}
}

func TestGetPresetByLanguageParser(t *testing.T) {
	presets := make(map[string]map[string]string)

	phpPreset := make(map[string]string)
	jsPreset := make(map[string]string)

	phpPreset["preset_language"] = "php"
	phpPreset["kool.yml"] = ""
	jsPreset["preset_language"] = "javascript"
	jsPreset["kool.yml"] = ""

	presets["php_language"] = phpPreset
	presets["javascript_language"] = jsPreset

	p := &DefaultParser{
		Presets: presets,
	}

	phpPresets := p.GetPresets("php")

	if len(phpPresets) != 1 || phpPresets[0] != "php_language" {
		t.Error("failed to get preset by language")
	}
}

func TestGetCreateCommandParser(t *testing.T) {
	presets := make(map[string]map[string]string)

	laravelPreset := make(map[string]string)

	laravelPreset["preset_create"] = "command"
	laravelPreset["kool.yml"] = ""

	presets["laravel"] = laravelPreset

	p := &DefaultParser{
		Presets: presets,
	}

	laravelCmd, _ := p.GetCreateCommand("laravel")

	if laravelCmd != laravelPreset["preset_create"] {
		t.Error("failed to get command")
	}
}

func TestFailGetCreateCommandParser(t *testing.T) {
	presets := make(map[string]map[string]string)

	laravelPreset := make(map[string]string)

	laravelPreset["preset_create"] = ""
	laravelPreset["kool.yml"] = ""

	presets["laravel"] = laravelPreset

	p := &DefaultParser{
		Presets: presets,
	}

	_, err := p.GetCreateCommand("laravel")

	if err != ErrCreateCommandtNotFoundOrEmpty {
		t.Errorf("failed, expected to get error %v got %v", ErrCreateCommandtNotFoundOrEmpty, err.Error())
	}
}

func TestIgnorePresetMetaKeysParser(t *testing.T) {
	originalStat := osStat

	defer func() { osStat = originalStat }()

	osStat = func(name string) (os.FileInfo, error) {
		return nil, nil
	}

	presets := make(map[string]map[string]string)

	testingPreset := make(map[string]string)

	testingPreset["preset_language"] = "php"
	testingPreset["kool.yml"] = ""

	presets["preset"] = testingPreset

	p := &DefaultParser{
		Presets: presets,
	}

	foundFiles := p.LookUpFiles("preset")

	if len(foundFiles) != 1 || foundFiles[0] != "kool.yml" {
		t.Errorf("expecting to find only 'kool.yml', found %v", foundFiles)
	}
}
