package templates

import (
	"errors"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestFakeParser(t *testing.T) {
	f := &FakeParser{}

	f.MockParseError = errors.New("parse error")
	err := f.Parse("template")

	if err == nil {
		t.Error("expecting error on Parse, got none")
	} else if err.Error() != "parse error" {
		t.Errorf("expecting error 'parse error' on Parse, got %v", err)
	}

	if val, ok := f.CalledParse["template"]; !val || !ok {
		t.Error("failed calling Parse")
	}

	f.MockGetServices = yaml.MapSlice{
		yaml.MapItem{Key: "serviceKey", Value: "serviceContent"},
	}
	services := f.GetServices()

	if !f.CalledGetServices || !reflect.DeepEqual(services, f.MockGetServices) {
		t.Error("failed calling GetServices")
	}

	f.MockGetVolumes = yaml.MapSlice{
		yaml.MapItem{Key: "volumeKey", Value: "volumeContent"},
	}
	volumes := f.GetVolumes()

	if !f.CalledGetVolumes || !reflect.DeepEqual(volumes, f.MockGetVolumes) {
		t.Error("failed calling GetVolumes")
	}
}
