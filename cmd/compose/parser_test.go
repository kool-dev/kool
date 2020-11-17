package compose

import (
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

func TestRemoveService(t *testing.T) {
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

func TestRemoveVolume(t *testing.T) {
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

func TestSetService(t *testing.T) {
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

func getYamlData(p *DefaultParser) yaml.MapSlice {
	parserStruct := reflect.ValueOf(p).Elem()
	reflectYamlData := parserStruct.FieldByName("yamlData")
	return reflect.NewAt(reflectYamlData.Type(), unsafe.Pointer(reflectYamlData.UnsafeAddr())).Elem().Interface().(yaml.MapSlice)
}
