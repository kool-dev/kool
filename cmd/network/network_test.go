package network

import "testing"

var createNetworkCalled bool

type FakeCommand struct{}

func (c *FakeCommand) LookPath() (err error) {
	return
}

func (c *FakeCommand) Interactive() (err error) {
	createNetworkCalled = true
	return
}

func (c *FakeCommand) Exec() (outStr string, err error) {
	return
}

type NetworkExistsCheckCmd struct {
	FakeCommand
}

func (c *NetworkExistsCheckCmd) Exec() (outStr string, err error) {
	outStr = "NetworkID"
	return
}

func TestDefaultHandler(t *testing.T) {
	var c Handler = NewHandler()

	if _, assert := c.(*DefaultHandler); !assert {
		t.Errorf("NewHandler() did not return a *DefaultHandler")
	}
}

func TestGlobalNetworkExists(t *testing.T) {
	var h Handler

	createNetworkCalled = false

	checkNetCmd := &NetworkExistsCheckCmd{}
	createNetCmd := &FakeCommand{}

	h = &DefaultHandler{checkNetCmd, createNetCmd}

	err := h.HandleGlobalNetwork()

	if createNetworkCalled {
		t.Errorf("HandleGlobalNetwork() should not try to create the global network if it already exists.")
	}

	if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	}
}

func TestGlobalNetworkNotExists(t *testing.T) {
	var h Handler

	createNetworkCalled = false

	checkNetCmd := &FakeCommand{}
	createNetCmd := &FakeCommand{}

	h = &DefaultHandler{checkNetCmd, createNetCmd}

	err := h.HandleGlobalNetwork()

	if !createNetworkCalled {
		t.Errorf("HandleGlobalNetwork() is not trying to create the global network when it not exists.")
	}

	if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	}
}
