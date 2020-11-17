package compose

// FakeParser implements all fake behaviors for using parser in tests.
type FakeParser struct {
	CalledLoad                                            map[string]bool
	CalledSetService                                      map[string]map[string]bool
	CalledRemoveService, CalledRemoveVolume, CalledString bool
}

// Load implements fake Load behavior
func (f *FakeParser) Load(compose string) (err error) {
	if f.CalledLoad == nil {
		f.CalledLoad = make(map[string]bool)
	}

	f.CalledLoad[compose] = true
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
	return
}
