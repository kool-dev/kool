package checker

// FakeChecker implements all fake behaviors for using checker in tests.
type FakeChecker struct {
	CalledCheck bool
	MockError   error
}

// Check implements fake Check behavior
func (f *FakeChecker) Check() (err error) {
	f.CalledCheck = true
	err = f.MockError
	return
}
