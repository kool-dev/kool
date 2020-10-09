package presets

// FakeParser implements all fake behaviors for using parser in tests.
type FakeParser struct {
	CalledExists, CalledLookUpFiles, CalledWriteFiles, CalledGetPresets bool

	MockExists     bool
	MockFoundFiles []string
	MockFileError  string
	MockError      error
	MockPresets    []string
}

// Exists check if preset exists
func (f *FakeParser) Exists(preset string) (exists bool) {
	f.CalledExists = true
	exists = f.MockExists
	return
}

// GetPresets get all presets names
func (f *FakeParser) GetPresets() (presets []string) {
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
	f.CalledWriteFiles = true
	fileError = f.MockFileError
	err = f.MockError
	return
}
