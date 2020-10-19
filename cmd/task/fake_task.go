package task

// FakeTask fake task runner for tests
type FakeTask struct {
	CalledRun bool
}

// Run fake task run
func (f *FakeTask) Run(message string, closure func() error) (err error) {
	f.CalledRun = true
	err = closure()
	return
}
