package network

import (
	"errors"
	"testing"
)

func TestFakeHandler(t *testing.T) {
	f := &FakeHandler{}

	_ = f.HandleGlobalNetwork("testing_network")

	if !f.CalledHandleGlobalNetwork || f.NetworkNameArg != "testing_network" {
		t.Error("failed to use mocked HandleGlobalNetwork function on FakeHandler")
	}
}

func TestFailedFakeChecker(t *testing.T) {
	f := &FakeHandler{MockError: errors.New("fake error")}

	err := f.HandleGlobalNetwork("testing_network")

	if err == nil {
		t.Error("failed to use mocked failed HandleGlobalNetwork function on FakeHandler")
	} else if err.Error() != "fake error" {
		t.Error("failed to use mocked failed HandleGlobalNetwork function on FakeHandler")
	}
}
