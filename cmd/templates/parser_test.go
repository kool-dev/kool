package templates

import (
	"gopkg.in/yaml.v2"
	"reflect"
	"testing"
	"unsafe"
)

const templateFile string = `version: "3.7"
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
scripts:
  script: command
  script2:
    - command
    - command2
`

func TestNewParser(t *testing.T) {
	p := NewParser()

	if _, assert := p.(*DefaultParser); !assert {
		t.Errorf("NewParser() did not return a *DefaultParser")
	}
}

func TestDefaultParser(t *testing.T) {
	p := NewParser()

	if err := p.Parse(templateFile); err != nil {
		t.Errorf("unexpected error parsing template file; error: %v", err)
	}

	yamlData := getYamlData(p.(*DefaultParser))

	parsed := new(TemplateFile)
	_ = yaml.Unmarshal([]byte(templateFile), &parsed)

	if !reflect.DeepEqual(yamlData, parsed) {
		t.Error("failed parsing template file content")
	}

	if services := p.GetServices(); !reflect.DeepEqual(yamlData.Services, services) {
		t.Error("failed getting services from template file content")
	}

	if volumes := p.GetVolumes(); !reflect.DeepEqual(yamlData.Volumes, volumes) {
		t.Error("failed getting volumes from template file content")
	}

	expectedScripts := make(map[string][]string)

	expectedScripts["script"] = []string{"command"}
	expectedScripts["script2"] = []string{"command", "command2"}

	if scripts := p.GetScripts(); !reflect.DeepEqual(expectedScripts, scripts) {
		t.Error("failed getting scripts from template file content")
	}
}

func getYamlData(p *DefaultParser) *TemplateFile {
	parserStruct := reflect.ValueOf(p).Elem()
	reflectYamlData := parserStruct.FieldByName("template")
	return reflect.NewAt(reflectYamlData.Type(), unsafe.Pointer(reflectYamlData.UnsafeAddr())).Elem().Interface().(*TemplateFile)
}
