package task

// FakeTask fake task runner for tests
type FakeTask struct {
	CalledRun bool
}

// Run fake task run
func (f *FakeTask) Run(message string, closure func() (interface{}, error)) (result interface{}, err error) {
	f.CalledRun = true
	result, err = closure()
	return
}
