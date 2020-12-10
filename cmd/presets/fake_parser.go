package presets

// FakeParser implements all fake behaviors for using parser in tests.
type FakeParser struct {
	CalledExists              bool
	CalledLookUpFiles         bool
	CalledWriteFiles          map[string]bool
	CalledGetPresets          bool
	CalledGetLanguages        bool
	CalledSetPresetKeyContent map[string]map[string]map[string]bool
	CalledGetTemplates        bool
	CalledLoadPresets         bool
	CalledLoadTemplates       bool
	CalledLoadConfigs         bool
	CalledGetConfig           map[string]bool

	MockExists         bool
	MockFoundFiles     []string
	MockFileError      string
	MockError          error
	MockLanguages      []string
	MockPresets        []string
	MockTemplates      map[string]map[string]string
	MockAllPresets     map[string]map[string]string
	MockAllTemplates   map[string]map[string]string
	MockAllConfigs     map[string]string
	MockConfig         map[string]*PresetConfig
	MockGetConfigError map[string]error
}

// Exists check if preset exists
func (f *FakeParser) Exists(preset string) (exists bool) {
	f.CalledExists = true
	exists = f.MockExists
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

// WriteFiles write preset files
func (f *FakeParser) WriteFiles(preset string) (fileError string, err error) {
	if f.CalledWriteFiles == nil {
		f.CalledWriteFiles = make(map[string]bool)
	}

	f.CalledWriteFiles[preset] = true
	fileError = f.MockFileError
	err = f.MockError
	return
}

// SetPresetKeyContent set preset key value
func (f *FakeParser) SetPresetKeyContent(preset string, key string, content string) {
	if f.CalledSetPresetKeyContent == nil {
		f.CalledSetPresetKeyContent = make(map[string]map[string]map[string]bool)
	}

	if _, ok := f.CalledSetPresetKeyContent[preset]; !ok {
		f.CalledSetPresetKeyContent[preset] = make(map[string]map[string]bool)
	}

	if _, ok := f.CalledSetPresetKeyContent[preset][key]; !ok {
		f.CalledSetPresetKeyContent[preset][key] = make(map[string]bool)
	}

	f.CalledSetPresetKeyContent[preset][key][content] = true
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

// LoadConfigs load the configs
func (f *FakeParser) LoadConfigs(configs map[string]string) {
	f.CalledLoadConfigs = true
	f.MockAllConfigs = configs
}

// GetConfig get preset config
func (f *FakeParser) GetConfig(preset string) (config *PresetConfig, err error) {
	if f.CalledGetConfig == nil {
		f.CalledGetConfig = make(map[string]bool)
	}

	if f.MockConfig == nil {
		f.MockConfig = make(map[string]*PresetConfig)
	}

	if f.MockGetConfigError == nil {
		f.MockGetConfigError = make(map[string]error)
	}

	f.CalledGetConfig[preset] = true
	if val, ok := f.MockConfig[preset]; ok {
		config = val
	}

	if val, ok := f.MockGetConfigError[preset]; ok {
		err = val
	}

	return
}
