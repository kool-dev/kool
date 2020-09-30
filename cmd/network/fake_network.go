package network

// FakeHandler implements all fake behaviors for using network handler in tests.
type FakeHandler struct {
	CalledHandleGlobalNetwork bool
	NetworkNameArg            string
	MockError                 error
}

// HandleGlobalNetwork implements fake HandleGlobalNetwork behavior
func (f *FakeHandler) HandleGlobalNetwork(networkName string) (err error) {
	f.CalledHandleGlobalNetwork = true
	f.NetworkNameArg = networkName
	err = f.MockError
	return
}
