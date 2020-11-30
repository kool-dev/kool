package presets

// FakeParser implements all fake behaviors for using parser in tests.
type FakeParser struct {
	CalledExists              bool
	CalledLookUpFiles         bool
	CalledWriteFile           map[string]map[string]bool
	CalledGetPresets          bool
	CalledGetLanguages        bool
	CalledGetPresetKeys       bool
	CalledGetPresetKeyContent map[string]map[string]bool
	CalledGetTemplates        bool
	CalledGetCreateCommand    bool
	CalledLoadPresets         bool
	CalledLoadTemplates       bool

	MockExists           bool
	MockFoundFiles       []string
	MockFileError        string
	MockError            error
	MockCreateCommand    string
	MockLanguages        []string
	MockPresets          []string
	MockPresetKeys       []string
	MockPresetKeyContent map[string]map[string]string
	MockTemplates        map[string]map[string]string
	MockAllPresets       map[string]map[string]string
	MockAllTemplates     map[string]map[string]string
}

// Exists check if preset exists
func (f *FakeParser) Exists(preset string) (exists bool) {
	f.CalledExists = true
	exists = f.MockExists
	return
}

// GetCreateCommand gets the command to create a new project
func (f *FakeParser) GetCreateCommand(preset string) (cmd string, err error) {
	f.CalledGetCreateCommand = true
	cmd = f.MockCreateCommand
	return
}

// GetLanguages get all presets languages
func (f *FakeParser) GetLanguages() (languages []string) {
	f.CalledGetLanguages = true
	languages = f.MockLanguages
	return
}

// GetPresets get all presets names
func (f *FakeParser) GetPresets(language string) (presets []string) {
	f.CalledGetPresets = true
	presets = f.MockPresets
	return
}

// LookUpFiles check if preset files exist
func (f *FakeParser) LookUpFiles(preset string) (foundFiles []string) {
	f.CalledLookUpFiles = true
	foundFiles = f.MockFoundFiles
	return
}

// WriteFile write preset files
func (f *FakeParser) WriteFile(fileName string, fileContent string) (fileError string, err error) {
	if f.CalledWriteFile == nil {
		f.CalledWriteFile = make(map[string]map[string]bool)
	}

	if _, ok := f.CalledWriteFile[fileName]; !ok {
		f.CalledWriteFile[fileName] = make(map[string]bool)
	}

	f.CalledWriteFile[fileName][fileContent] = true
	fileError = f.MockFileError
	err = f.MockError
	return
}

// GetPresetKeys get preset file contents
func (f *FakeParser) GetPresetKeys(preset string) (keys []string) {
	f.CalledGetPresetKeys = true
	keys = f.MockPresetKeys
	return
}

// GetPresetKeyContent get preset key value
func (f *FakeParser) GetPresetKeyContent(preset string, key string) (value string) {
	if f.CalledGetPresetKeyContent == nil {
		f.CalledGetPresetKeyContent = make(map[string]map[string]bool)
	}

	if _, ok := f.CalledGetPresetKeyContent[preset]; !ok {
		f.CalledGetPresetKeyContent[preset] = make(map[string]bool)
	}

	f.CalledGetPresetKeyContent[preset][key] = true
	value = f.MockPresetKeyContent[preset][key]
	return
}

// GetTemplates get all templates
func (f *FakeParser) GetTemplates() (templates map[string]map[string]string) {
	f.CalledGetTemplates = true
	if f.MockTemplates == nil {
		f.MockTemplates = make(map[string]map[string]string)
	}
	templates = f.MockTemplates
	return
}

//LoadPresets loads all presets
func (f *FakeParser) LoadPresets(presets map[string]map[string]string) {
	f.CalledLoadPresets = true
	f.MockAllPresets = presets
}

//LoadTemplates loads all templates
func (f *FakeParser) LoadTemplates(templates map[string]map[string]string) {
	f.CalledLoadTemplates = true
	f.MockAllTemplates = templates
}
