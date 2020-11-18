package compose

import (
	"errors"
	"fmt"
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
`

const composeWithoutService = `version: "3.7"
services:
  service2:
    image: service-image2
    volumes:
    - service2:/app:delegated
volumes:
  service: null
  service2: null
`

const composeWithouServiceVolume string = `version: "3.7"
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
  service2: null
`

const newComposeService string = `image: new-service-image
volumes:
- service:/app:delegated
`

const composeWithNewService string = `version: "3.7"
services:
  service:
    image: new-service-image
    volumes:
    - service:/app:delegated
  service2:
    image: service-image2
    volumes:
    - service2:/app:delegated
volumes:
  service: null
  service2: null
`

func TestNewParser(t *testing.T) {
	p := NewParser()

	if _, assert := p.(*DefaultParser); !assert {
		t.Errorf("NewParser() did not return a *DefaultParser")
	}
}

func TestLoadDefaultParser(t *testing.T) {
	p := NewParser()

	if err := p.Load(composeFile); err != nil {
		t.Errorf("unexpected error loading docker compose file; error: %v", err)
	}

	yamlData := getYamlData(p.(*DefaultParser))

	parsed := yaml.MapSlice{}
	_ = yaml.Unmarshal([]byte(composeFile), &parsed)

	if !reflect.DeepEqual(yamlData, parsed) {
		t.Error("failed loading docker compose file content")
	}
}

func TestStringDefaultParser(t *testing.T) {
	p := NewParser()

	_ = p.Load(composeFile)

	content, err := p.String()

	if err != nil {
		t.Errorf("unexpected error getting docker compose file content; error: %v", err)
	}

	if strings.TrimSpace(content) != strings.TrimSpace(composeFile) {
		t.Errorf("expecting content '%s', got '%s'", strings.TrimSpace(composeFile), strings.TrimSpace(content))
	}
}

func TestRemoveServiceDefaultParser(t *testing.T) {
	p := NewParser()

	_ = p.Load(composeFile)

	p.RemoveService("service")

	yamlData := getYamlData(p.(*DefaultParser))

	parsed := yaml.MapSlice{}
	_ = yaml.Unmarshal([]byte(composeWithoutService), &parsed)

	if !reflect.DeepEqual(yamlData, parsed) {
		t.Error("failed removing docker compose file service")
	}
}

func TestRemoveVolumeDefaultParser(t *testing.T) {
	p := NewParser()

	_ = p.Load(composeFile)

	p.RemoveVolume("service")

	yamlData := getYamlData(p.(*DefaultParser))

	parsed := yaml.MapSlice{}
	_ = yaml.Unmarshal([]byte(composeWithouServiceVolume), &parsed)

	if !reflect.DeepEqual(yamlData, parsed) {
		t.Error("failed removing docker compose file volume")
	}
}

func TestSetServiceDefaultParser(t *testing.T) {
	p := NewParser()

	_ = p.Load(composeFile)

	if err := p.SetService("service", newComposeService); err != nil {
		t.Errorf("unexpected error setting docker compose service; error: %v", err)
	}

	yamlData := getYamlData(p.(*DefaultParser))

	parsed := yaml.MapSlice{}
	_ = yaml.Unmarshal([]byte(composeWithNewService), &parsed)

	if !reflect.DeepEqual(yamlData, parsed) {
		t.Error("failed setting docker compose file service")
	}
}

func TestErrorSetServiceDefaultParser(t *testing.T) {
	p := NewParser()
	_ = p.Load(composeFile)

	originalYamlUnmarshalFn := yamlUnmarshalFn
	defer func() {
		yamlUnmarshalFn = originalYamlUnmarshalFn
	}()

	yamlUnmarshalFn = func(in []byte, out interface{}) error {
		fmt.Println("unmarshal")
		return errors.New("yaml unmarshal error")
	}

	err := p.SetService("service", newComposeService)

	if err == nil {
		t.Error("expecting error 'yaml unmarshal error', got none")
	} else if err.Error() != "yaml unmarshal error" {
		t.Errorf("expecting error 'yaml unmarshal error', got %v", err)
	}
}

func TestErrorStringDefaultParser(t *testing.T) {
	p := NewParser()
	_ = p.Load(composeFile)

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

func getYamlData(p *DefaultParser) yaml.MapSlice {
	parserStruct := reflect.ValueOf(p).Elem()
	reflectYamlData := parserStruct.FieldByName("yamlData")
	return reflect.NewAt(reflectYamlData.Type(), unsafe.Pointer(reflectYamlData.UnsafeAddr())).Elem().Interface().(yaml.MapSlice)
}
