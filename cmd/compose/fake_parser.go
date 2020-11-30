package compose

// FakeParser implements all fake behaviors for using parser in tests.
type FakeParser struct {
	CalledLoad                              map[string]bool
	CalledSetService                        map[string]map[string]bool
	CalledRemoveService, CalledRemoveVolume map[string]bool
	CalledString                            bool
	MockLoadError                           error
	MockSetServiceError                     error
	MockStringError                         error
}

// Load implements fake Load behavior
func (f *FakeParser) Load(compose string) (err error) {
	if f.CalledLoad == nil {
		f.CalledLoad = make(map[string]bool)
	}

	f.CalledLoad[compose] = true
	err = f.MockLoadError
	return
}

// SetService implements fake SetService behavior
func (f *FakeParser) SetService(service string, content string) (err error) {
	if f.CalledSetService == nil {
		f.CalledSetService = make(map[string]map[string]bool)
	}

	if f.CalledSetService[service] == nil {
		f.CalledSetService[service] = make(map[string]bool)
	}

	f.CalledSetService[service][content] = true
	err = f.MockSetServiceError
	return
}

// RemoveService implements fake RemoveService behavior
func (f *FakeParser) RemoveService(service string) {
	if f.CalledRemoveService == nil {
		f.CalledRemoveService = make(map[string]bool)
	}

	f.CalledRemoveService[service] = true
}

// RemoveVolume implements fake RemoveVolume behavior
func (f *FakeParser) RemoveVolume(volume string) {
	if f.CalledRemoveVolume == nil {
		f.CalledRemoveVolume = make(map[string]bool)
	}

	f.CalledRemoveVolume[volume] = true
}

// String implements fake String behavior
func (f *FakeParser) String() (content string, err error) {
	f.CalledString = true
	err = f.MockStringError
	return
}
