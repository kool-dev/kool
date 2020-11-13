package presets

// FakeParser implements all fake behaviors for using parser in tests.
type FakeParser struct {
	CalledExists, CalledLookUpFiles, CalledWriteFile, CalledGetPresets, CalledGetLanguages, CalledGetPresetKeys, CalledGetPresetKeyContent bool

	MockExists           bool
	MockFoundFiles       []string
	MockFileError        string
	MockError            error
	MockLanguages        []string
	MockPresets          []string
	MockPresetKeys       []string
	MockPresetKeyContent string
	MockTemplates        map[string]map[string]string
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

// WriteFile write preset files
func (f *FakeParser) WriteFile(fileName string, fileContent string) (fileError string, err error) {
	f.CalledWriteFile = true
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
	f.CalledGetPresetKeyContent = true
	value = f.MockPresetKeyContent
	return
}

// GetTemplates get all templates
func (f *FakeParser) GetTemplates() (templates map[string]map[string]string) {
	if f.MockTemplates == nil {
		f.MockTemplates = make(map[string]map[string]string)
	}
	templates = f.MockTemplates
	return
}
