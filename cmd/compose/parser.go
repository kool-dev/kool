package compose

import "gopkg.in/yaml.v2"

type yamlUnmarshalFnType func([]byte, interface{}) error
type yamlMarshalFnType func(interface{}) ([]byte, error)

// Compose represents a docker-compose file
type Compose struct {
	Version  string        `yaml:"version"`
	Services yaml.MapSlice `yaml:"services"`
	Volumes  yaml.MapSlice `yaml:"volumes,omitempty"`
	Networks yaml.MapSlice `yaml:"networks,omitempty"`
}

// Parser holds logic for handling docker-compose
type Parser interface {
	Parse(string) error
	GetServices() yaml.MapSlice
	SetService(string, interface{})
	GetVolumes() yaml.MapSlice
	SetVolume(string)
	String() (string, error)
}

// DefaultParser holds data for docker-compose
type DefaultParser struct {
	compose *Compose
}

var (
	yamlUnmarshalFn yamlUnmarshalFnType = yaml.Unmarshal
	yamlMarshalFn   yamlMarshalFnType   = yaml.Marshal
)

// NewParser creates new docker-compose parser
func NewParser() Parser {
	compose := &Compose{
		Version: "3.7",
		Networks: yaml.MapSlice{
			yaml.MapItem{Key: "kool_local"},
			yaml.MapItem{
				Key: "kool_global",
				Value: yaml.MapSlice{
					yaml.MapItem{Key: "external", Value: true},
					yaml.MapItem{Key: "name", Value: "${KOOL_GLOBAL_NETWORK:-kool_global}"},
				},
			},
		},
	}
	return &DefaultParser{compose}
}

// Parse parse content to yaml
func (p *DefaultParser) Parse(content string) (err error) {
	compose := new(Compose)
	if err = yamlUnmarshalFn([]byte(content), &compose); err != nil {
		return
	}

	p.compose = compose
	return
}

// GetServices get compose services
func (p *DefaultParser) GetServices() yaml.MapSlice {
	return p.compose.Services
}

// SetService set docker-compose service
func (p *DefaultParser) SetService(serviceName string, serviceContent interface{}) {
	for index, service := range p.compose.Services {
		if service.Key == serviceName {
			p.compose.Services[index].Value = serviceContent
			return
		}
	}

	p.compose.Services = append(p.compose.Services, yaml.MapItem{
		Key:   serviceName,
		Value: serviceContent,
	})
}

// GetVolumes get compose volumes
func (p *DefaultParser) GetVolumes() yaml.MapSlice {
	return p.compose.Volumes
}

// SetVolume remove a docker-compose volume
func (p *DefaultParser) SetVolume(volume string) {
	for _, v := range p.compose.Volumes {
		if v.Key == volume {
			return
		}
	}

	p.compose.Volumes = append(p.compose.Volumes, yaml.MapItem{
		Key: volume,
	})
}

// String returns docker-compose as string
func (p *DefaultParser) String() (content string, err error) {
	var parsedBytes []byte

	if parsedBytes, err = yamlMarshalFn(p.compose); err != nil {
		return
	}

	content = string(parsedBytes)
	return
}
