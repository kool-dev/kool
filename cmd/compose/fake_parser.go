package compose

// FakeParser implements all fake behaviors for using parser in tests.
type FakeParser struct {
	CalledParse      map[string]bool
	CalledSetService map[string]bool
	CalledSetVolume  map[string]bool
	CalledString     bool
	MockParseError   error
	MockStringError  error
}

// Parse implements fake Parse behavior
func (f *FakeParser) Parse(content string) (err error) {
	if f.CalledParse == nil {
		f.CalledParse = make(map[string]bool)
	}

	f.CalledParse[content] = true
	err = f.MockParseError
	return
}

// SetService implements fake SetService behavior
func (f *FakeParser) SetService(service string, content interface{}) {
	if f.CalledSetService == nil {
		f.CalledSetService = make(map[string]bool)
	}

	f.CalledSetService[service] = true
}

// SetVolume implements fake SetVolume behavior
func (f *FakeParser) SetVolume(volume string) {
	if f.CalledSetVolume == nil {
		f.CalledSetVolume = make(map[string]bool)
	}

	f.CalledSetVolume[volume] = true
}

// String implements fake String behavior
func (f *FakeParser) String() (content string, err error) {
	f.CalledString = true
	err = f.MockStringError
	return
}
