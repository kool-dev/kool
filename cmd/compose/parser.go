package compose

import "gopkg.in/yaml.v2"

// Parser holds logic for handling docker-compose
type Parser interface {
	Load(string) error
	SetService(string, string) error
	RemoveService(string)
	RemoveVolume(string)
	String() (string, error)
}

// DefaultParser holds data for docker-compose
type DefaultParser struct {
	yamlData yaml.MapSlice
}

// NewParser creates new docker-compose parser
func NewParser() Parser {
	return &DefaultParser{}
}

// Load loads docker-compose into Parser
func (p *DefaultParser) Load(compose string) (err error) {
	p.yamlData, err = parseYaml(compose)
	return
}

// SetService set docker-compose service
func (p *DefaultParser) SetService(serviceName string, serviceContent string) (err error) {
	for sectionKey, section := range p.yamlData {
		if section.Key == "services" {
			for serviceKey, service := range section.Value.(yaml.MapSlice) {
				if service.Key == serviceName {
					var template yaml.MapSlice

					if template, err = parseYaml(serviceContent); err != nil {
						return
					}

					p.yamlData[sectionKey].Value.(yaml.MapSlice)[serviceKey].Value = template
					return
				}
			}
		}
	}
	return
}

// RemoveService remove a docker-compose service
func (p *DefaultParser) RemoveService(service string) {
	p.yamlData = removeSubItem(p.yamlData, "services", service)
}

// RemoveVolume remove a docker-compose volume
func (p *DefaultParser) RemoveVolume(volume string) {
	p.yamlData = removeSubItem(p.yamlData, "volumes", volume)
}

// String returns docker-compose as string
func (p *DefaultParser) String() (content string, err error) {
	var parsedBytes []byte

	if parsedBytes, err = yaml.Marshal(p.yamlData); err != nil {
		return
	}

	content = string(parsedBytes)
	return
}

func parseYaml(content string) (yaml.MapSlice, error) {
	parsed := yaml.MapSlice{}

	if err := yaml.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, err
	}

	return parsed, nil
}

func removeSubItem(originalCompose yaml.MapSlice, item string, subItem string) (compose yaml.MapSlice) {
	for _, section := range originalCompose {
		if section.Key != item {
			compose = append(compose, section)
			continue
		}

		var services yaml.MapSlice
		for _, service := range section.Value.(yaml.MapSlice) {
			if service.Key != subItem {
				services = append(services, service)
			}
		}

		compose = append(compose, yaml.MapItem{Key: "services", Value: services})
	}

	return
}
