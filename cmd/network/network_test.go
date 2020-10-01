package network

import (
	"kool-dev/kool/cmd/builder"
	"testing"
)

func TestDefaultHandler(t *testing.T) {
	var c Handler = NewHandler()

	if _, assert := c.(*DefaultHandler); !assert {
		t.Errorf("NewHandler() did not return a *DefaultHandler")
	}
}

func TestGlobalNetworkExists(t *testing.T) {
	var h Handler

	checkNetCmd := &builder.FakeCommand{MockExecOut: "NetworkID"}
	createNetCmd := &builder.FakeCommand{}

	h = &DefaultHandler{checkNetCmd, createNetCmd}

	err := h.HandleGlobalNetwork("global_network")

	if !h.(*DefaultHandler).CheckNetworkCmd.(*builder.FakeCommand).CalledExec {
		t.Errorf("HandleGlobalNetwork() did not check if network exists.")
	}

	if h.(*DefaultHandler).CreateNetworkCmd.(*builder.FakeCommand).CalledInteractive {
		t.Errorf("HandleGlobalNetwork() should not try to create the global network if it already exists.")
	}

	if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	}
}

func TestGlobalNetworkNotExists(t *testing.T) {
	var h Handler

	checkNetCmd := &builder.FakeCommand{}
	createNetCmd := &builder.FakeCommand{}

	h = &DefaultHandler{checkNetCmd, createNetCmd}

	err := h.HandleGlobalNetwork("global_network")

	if !h.(*DefaultHandler).CreateNetworkCmd.(*builder.FakeCommand).CalledInteractive {
		t.Errorf("HandleGlobalNetwork() is not trying to create the global network when it not exists.")
	}

	if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	}
}
