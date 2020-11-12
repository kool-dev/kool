package compose

// FakeParser implements all fake behaviors for using parser in tests.
type FakeParser struct {
	CalledLoad, CalledSetService, CalledRemoveService, CalledRemoveVolume, CalledString bool
	MockLoadError, MockSetServiceError, MockStringError error
	MockString string
}

// Load implements fake Load behavior
func (f *FakeParser) Load(compose string) (err error) {
	f.CalledLoad = true
	err = f.MockLoadError
	return
}

// SetService implements fake SetService behavior
func (f *FakeParser) SetService(service string, content string) (err error) {
	f.CalledSetService = true
	err = f.MockSetServiceError
	return
}

// RemoveService implements fake RemoveService behavior
func (f *FakeParser) RemoveService(service string) {
	f.CalledRemoveService = true
}

// RemoveVolume implements fake RemoveVolume behavior
func (f *FakeParser) RemoveVolume(volume string) {
	f.CalledRemoveVolume = true
}

// String implements fake String behavior
func (f *FakeParser) String() (content string, err error) {
	f.CalledString = true
	err = f.MockStringError
	content = f.MockString
	return
}
