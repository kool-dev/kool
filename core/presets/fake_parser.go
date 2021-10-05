package presets

// FakeParser implements all fake behaviors for using parser in tests.
type FakeParser struct {
	CalledExists       bool
	CalledLookUpFiles  bool
	CalledWriteFiles   map[string]bool
	CalledGetPresets   bool
	CalledGetLanguages bool
	CalledGetConfig    map[string]bool

	MockExists         bool
	MockFoundFiles     []string
	MockFileError      string
	MockError          error
	MockLanguages      []string
	MockPresets        []string
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
