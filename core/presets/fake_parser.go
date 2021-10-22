package presets

// FakeParser implements all fake behaviors for using parser in tests.
type FakeParser struct {
	CalledExists     bool
	CalledGetTags    bool
	CalledGetPresets bool
	CalledInstall    bool
	CalledCreate     bool
	CalledAdd        bool

	MockExists     bool
	MockGetTags    []string
	MockGetPresets []string
	MockInstall    error
	MockCreate     error
	MockAdd        error
}

// Exists check if preset exists
func (f *FakeParser) Exists(preset string) (exists bool) {
	f.CalledExists = true
	exists = f.MockExists
	return
}

// GetTags get all presets tags
func (f *FakeParser) GetTags() (languages []string) {
	f.CalledGetTags = true
	languages = f.MockGetTags
	return
}

// GetPresets get all presets names
func (f *FakeParser) GetPresets(tag string) (presets []string) {
	f.CalledGetPresets = true
	presets = f.MockGetPresets
	return
}

// Install
func (f *FakeParser) Install(tag string) (err error) {
	f.CalledInstall = true
	err = f.MockInstall
	return
}

// Create
func (f *FakeParser) Create(tag string) (err error) {
	f.CalledCreate = true
	err = f.MockCreate
	return
}

// Add
func (f *FakeParser) Add(tag string) (err error) {
	f.CalledAdd = true
	err = f.MockAdd
	return
}
