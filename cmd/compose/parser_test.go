package compose

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"gopkg.in/yaml.v2"
)

const composeFile string = `version: "3.7"
services:
  service:
    image: service-image
    volumes:
    - service:/app:delegated
  service2:
    image: service-image2
    volumes:
    - service2:/app:delegated
volumes:
  service: null
  service2: null
networks:
  kool_local: null
  kool_global:
    external: true
    name: ${KOOL_GLOBAL_NETWORK:-kool_global}
`

func TestNewParser(t *testing.T) {
	p := NewParser()

	if _, assert := p.(*DefaultParser); !assert {
		t.Errorf("NewParser() did not return a *DefaultParser")
	}
}

func TestParseDefaultParser(t *testing.T) {
	p := NewParser()

	if err := p.Parse(composeFile); err != nil {
		t.Errorf("unexpected error loading docker compose file; error: %v", err)
	}

	yamlData := getYamlData(p.(*DefaultParser))

	parsed := new(Compose)
	_ = yaml.Unmarshal([]byte(composeFile), &parsed)

	if !reflect.DeepEqual(yamlData, parsed) {
		t.Error("failed loading docker compose file content")
	}
}

func TestStringDefaultParser(t *testing.T) {
	p := NewParser()

	_ = p.Parse(composeFile)

	content, err := p.String()

	if err != nil {
		t.Errorf("unexpected error getting docker compose file content; error: %v", err)
	}

	if strings.TrimSpace(content) != strings.TrimSpace(composeFile) {
		t.Errorf("expecting content '%s', got '%s'", strings.TrimSpace(composeFile), strings.TrimSpace(content))
	}
}

func TestServicesDefaultParser(t *testing.T) {
	p := NewParser()

	appService := yaml.MapSlice{
		yaml.MapItem{Key: "image", Value: "kooldev/image"},
	}

	services := yaml.MapSlice{
		yaml.MapItem{Key: "app", Value: appService},
	}

	p.SetService("app", appService)

	if getServices := p.GetServices(); !reflect.DeepEqual(services, getServices) {
		t.Error("failed handling services")
	}

	appService = yaml.MapSlice{
		yaml.MapItem{Key: "image", Value: "kooldev/image2"},
	}

	services = yaml.MapSlice{
		yaml.MapItem{Key: "app", Value: appService},
	}

	p.SetService("app", appService)

	if getServices := p.GetServices(); !reflect.DeepEqual(services, getServices) {
		t.Error("failed handling services")
	}
}

func TestVolumesDefaultParser(t *testing.T) {
	p := NewParser()

	volumes := yaml.MapSlice{
		yaml.MapItem{Key: "database"},
	}

	p.SetVolume("database")

	if getVolumes := p.GetVolumes(); !reflect.DeepEqual(volumes, getVolumes) {
		t.Error("failed handling volumes")
	}

	p.SetVolume("database")

	if getVolumes := p.GetVolumes(); !reflect.DeepEqual(volumes, getVolumes) {
		t.Error("failed handling volumes")
	}
}

func TestErrorStringDefaultParser(t *testing.T) {
	p := NewParser()
	_ = p.Parse(composeFile)

	originalYamlMarshalFn := yamlMarshalFn
	defer func() {
		yamlMarshalFn = originalYamlMarshalFn
	}()

	yamlMarshalFn = func(in interface{}) ([]byte, error) {
		return nil, errors.New("yaml marshal error")
	}

	_, err := p.String()

	if err == nil {
		t.Error("expecting error 'yaml marshal error', got none")
	} else if err.Error() != "yaml marshal error" {
		t.Errorf("expecting error 'yaml marshal error', got %v", err)
	}
}

func TestErrorParseDefaultParser(t *testing.T) {
	originalYamlUnmarshalFn := yamlUnmarshalFn
	defer func() {
		yamlUnmarshalFn = originalYamlUnmarshalFn
	}()

	yamlUnmarshalFn = func(in []byte, out interface{}) error {
		return errors.New("yaml unmarshal error")
	}

	p := NewParser()
	err := p.Parse(composeFile)

	if err == nil {
		t.Error("expecting error 'yaml unmarshal error', got none")
	} else if err.Error() != "yaml unmarshal error" {
		t.Errorf("expecting error 'yaml unmarshal error', got %v", err)
	}
}

func getYamlData(p *DefaultParser) *Compose {
	parserStruct := reflect.ValueOf(p).Elem()
	reflectYamlData := parserStruct.FieldByName("compose")
	return reflect.NewAt(reflectYamlData.Type(), unsafe.Pointer(reflectYamlData.UnsafeAddr())).Elem().Interface().(*Compose)
}
