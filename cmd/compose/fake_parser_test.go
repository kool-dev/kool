package compose

import (
	"errors"
	"testing"
)

func TestFakeParser(t *testing.T) {
	var err error

	f := &FakeParser{}

	f.MockLoadError = errors.New("load error")
	err = f.Load("compose")

	if err == nil {
		t.Error("expecting error on Load, got none")
	} else if err.Error() != "load error" {
		t.Errorf("expecting error 'load error' on Load, got %v", err)
	}

	if val, ok := f.CalledLoad["compose"]; !ok || !val {
		t.Error("failed calling Load")
	}

	f.MockSetServiceError = errors.New("set service error")
	err = f.SetService("service", "content")

	if err == nil {
		t.Error("expecting error on SetService, got none")
	} else if err.Error() != "set service error" {
		t.Errorf("expecting error 'set service error' on SetService, got %v", err)
	}

	if val, ok := f.CalledSetService["service"]["content"]; !ok || !val {
		t.Error("failed calling SetService")
	}

	f.RemoveService("service")

	if val, ok := f.CalledRemoveService["service"]; !ok || !val {
		t.Error("failed calling RemoveService")
	}

	f.RemoveVolume("volume")

	if val, ok := f.CalledRemoveVolume["volume"]; !ok || !val {
		t.Error("failed calling RemoveVolume")
	}

	f.MockStringError = errors.New("string error")

	_, err = f.String()

	if err == nil {
		t.Error("expecting error on String, got none")
	} else if err.Error() != "string error" {
		t.Errorf("expecting error 'string error' on String, got %v", err)
	}

	if !f.CalledString {
		t.Error("failed calling String")
	}
}
