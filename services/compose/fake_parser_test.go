package compose

import (
	"errors"
	"testing"
)

func TestFakeParser(t *testing.T) {
	var err error

	f := &FakeParser{}

	f.MockParseError = errors.New("parse error")
	err = f.Parse("compose")

	if err == nil {
		t.Error("expecting error on Parse, got none")
	} else if err.Error() != "parse error" {
		t.Errorf("expecting error 'parse error' on Parse, got %v", err)
	}

	if val, ok := f.CalledParse["compose"]; !ok || !val {
		t.Error("failed calling Parse")
	}

	f.SetService("service", "content")

	if val, ok := f.CalledSetService["service"]; !ok || !val {
		t.Error("failed calling SetService")
	}

	f.SetVolume("volume")

	if val, ok := f.CalledSetVolume["volume"]; !ok || !val {
		t.Error("failed calling SetVolume")
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
