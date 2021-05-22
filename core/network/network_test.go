package network

import (
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/shell"
	"testing"
)

func TestDefaultHandler(t *testing.T) {
	var c Handler = NewHandler(&shell.FakeShell{})

	if _, assert := c.(*DefaultHandler); !assert {
		t.Errorf("NewHandler() did not return a *DefaultHandler")
	}
}

func TestGlobalNetworkExists(t *testing.T) {
	var h Handler

	checkNetCmd := &builder.FakeCommand{MockExecOut: "NetworkID"}
	checkNetCmd.MockCmd = "check"

	createNetCmd := &builder.FakeCommand{}
	createNetCmd.MockCmd = "create"

	s := &shell.FakeShell{}
	h = &DefaultHandler{checkNetCmd, createNetCmd, s}

	err := h.HandleGlobalNetwork("global_network")

	if val, ok := h.(*DefaultHandler).shell.(*shell.FakeShell).CalledExec["check"]; !val || !ok {
		t.Errorf("HandleGlobalNetwork() did not check if network exists.")
	}

	if val, ok := h.(*DefaultHandler).shell.(*shell.FakeShell).CalledInteractive["create"]; val && ok {
		t.Errorf("HandleGlobalNetwork() should not try to create the global network if it already exists.")
	}

	if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	}
}

func TestGlobalNetworkNotExists(t *testing.T) {
	var h Handler

	checkNetCmd := &builder.FakeCommand{}
	checkNetCmd.MockCmd = "check"

	createNetCmd := &builder.FakeCommand{}
	createNetCmd.MockCmd = "create"

	s := &shell.FakeShell{}
	h = &DefaultHandler{checkNetCmd, createNetCmd, s}

	err := h.HandleGlobalNetwork("global_network")

	if err != nil {
		t.Errorf("Expected no errors, got %v", err)
	}

	if val, ok := h.(*DefaultHandler).shell.(*shell.FakeShell).CalledInteractive["create"]; !val || !ok {
		t.Errorf("HandleGlobalNetwork() is not trying to create the global network when it not exists.")
	}
}
