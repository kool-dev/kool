package presets

import (
	"os"
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

func TestExistsParser(t *testing.T) {
	presets := make(map[string]map[string]string)
	laravelPreset := make(map[string]string)

	laravelPreset["kool.yml"] = ""
	presets["laravel"] = laravelPreset

	p := NewParser()
	p.LoadPresets(presets)

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
	var allPresets []string
	p := NewParser()

	allPresets = p.GetPresets("")

	if len(allPresets) != 0 {
		t.Error("presets should be empty")
	}

	presets := make(map[string]map[string]string)
	laravelPreset := make(map[string]string)
	symfonyPreset := make(map[string]string)

	laravelPreset["kool.yml"] = ""
	symfonyPreset["kool.yml"] = ""

	presets["laravel"] = laravelPreset
	presets["symfony"] = symfonyPreset

	p.LoadPresets(presets)

	allPresets = p.GetPresets("")

	if len(allPresets) != 2 || allPresets[0] != "laravel" || allPresets[1] != "symfony" {
		t.Error("failed to get all presets")
	}
}

func TestGetLanguagesParser(t *testing.T) {
	var allLanguages []string
	p := NewParser()

	allLanguages = p.GetLanguages()

	if len(allLanguages) != 0 {
		t.Error("languages should be empty")
	}

	presets := make(map[string]map[string]string)

	laravelPreset := make(map[string]string)
	symfonyPreset := make(map[string]string)

	laravelPreset["preset_language"] = "php"
	laravelPreset["kool.yml"] = ""
	symfonyPreset["preset_language"] = "php"
	symfonyPreset["kool.yml"] = ""

	presets["laravel"] = laravelPreset
	presets["symfony"] = symfonyPreset

	p.LoadPresets(presets)

	allLanguages = p.GetLanguages()

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

	p := NewParser()
	p.LoadPresets(presets)

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

	p := NewParser()
	p.LoadPresets(presets)

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

	p := NewParser()
	p.LoadPresets(presets)

	_, err := p.GetCreateCommand("laravel")

	if err != ErrCreateCommandtNotFoundOrEmpty {
		t.Errorf("failed, expected to get error %v got %v", ErrCreateCommandtNotFoundOrEmpty, err.Error())
	}
}

func TestIgnorePresetMetaKeysParser(t *testing.T) {
	presets := make(map[string]map[string]string)

	testingPreset := make(map[string]string)

	testingPreset["preset_language"] = "php"
	testingPreset["kool.yml"] = ""

	presets["preset"] = testingPreset

	fs := afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "kool.yml", []byte("scripts"), os.ModePerm)

	p := NewParserFS(fs)
	p.LoadPresets(presets)

	foundFiles := p.LookUpFiles("preset")

	if len(foundFiles) != 1 || foundFiles[0] != "kool.yml" {
		t.Errorf("expecting to find only 'kool.yml', found %v", foundFiles)
	}
}

func TestGetPresetKeysAndContentsParser(t *testing.T) {
	presets := make(map[string]map[string]string)

	preset := make(map[string]string)

	preset["key1"] = "value1"
	preset["key2"] = "value2"
	preset["key3"] = "value3"

	presets["preset"] = preset

	p := NewParser()
	p.LoadPresets(presets)

	keys := p.GetPresetKeys("preset")

	if len(keys) != 3 || keys[0] != "key1" || keys[1] != "key2" || keys[2] != "key3" {
		t.Errorf("expecting to find keys '[key1 key2 key3]', found %v", keys)
	}

	content := p.GetPresetKeyContent("preset", "key2")

	if content != "value2" {
		t.Errorf("expecting to find value 'value2', found %s", content)
	}
}

func TestGetTemplatesParser(t *testing.T) {
	var allTemplates map[string]map[string]string
	p := NewParser()

	allTemplates = p.GetTemplates()

	if allTemplates != nil {
		t.Error("templates should be empty")
	}

	templates := make(map[string]map[string]string)

	template := make(map[string]string)

	template["key1"] = "value1"
	template["key2"] = "value2"
	template["key3"] = "value3"

	templates["template"] = template

	p.LoadTemplates(templates)

	allTemplates = p.GetTemplates()

	if !reflect.DeepEqual(allTemplates, templates) {
		t.Error("failed to get all presets")
	}
}

func TestWriteFileParser(t *testing.T) {
	fs := afero.NewMemMapFs()

	p := NewParserFS(fs)

	if _, err := p.WriteFile("kool.yml", "scripts"); err != nil {
		t.Errorf("unexpected error writing file, err: %v", err)
	}

	if _, err := fs.Stat("kool.yml"); os.IsNotExist(err) {
		t.Error("could not write the file 'kool.yml'")
	}
}

func TestLoadPresetsParser(t *testing.T) {
	presets := map[string]map[string]string{
		"laravel": {"preset_create": "command"},
	}

	p := &DefaultParser{}
	p.LoadPresets(presets)

	if ok := reflect.DeepEqual(p.Presets, presets); !ok {
		t.Error("did not load the presets correctly")
	}
}

func TestLoadTemplatesParser(t *testing.T) {
	templates := map[string]map[string]string{
		"service": {"serviceKey": "content"},
	}

	p := &DefaultParser{}
	p.LoadTemplates(templates)

	if ok := reflect.DeepEqual(p.Templates, templates); !ok {
		t.Error("did not load the templates correctly")
	}
}
