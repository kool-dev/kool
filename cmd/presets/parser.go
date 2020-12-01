package presets

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// PresetConfigQuestion preset config question
type PresetConfigQuestion struct {
	Message string   `yaml:"message"`
	Options []string `yaml:"options"`
}

// PresetConfig preset config
type PresetConfig struct {
	Language  string              `yaml:"language"`
	Commands  map[string][]string `yaml:"commands"`
	Questions map[string]PresetConfigQuestion
}

// DefaultParser holds presets parsing data
type DefaultParser struct {
	Presets   map[string]map[string]string
	Templates map[string]map[string]string
	Configs   map[string]string
	fs        afero.Fs
}

// Parser holds presets parsing logic
type Parser interface {
	Exists(string) bool
	GetLanguages() []string
	GetPresets(string) []string
	LookUpFiles(string) []string
	LoadPresets(map[string]map[string]string)
	LoadTemplates(map[string]map[string]string)
	LoadConfigs(map[string]string)
	WriteFile(string, string) (string, error)
	GetPresetKeys(string) []string
	GetPresetKeyContent(string, string) string
	GetTemplates() map[string]map[string]string
	GetConfig(string) (*PresetConfig, error)
}

// NewParser creates a new preset default parser
func NewParser() Parser {
	return &DefaultParser{
		fs: afero.NewOsFs(),
	}
}

// NewParserFS creates a new preset default parser with file system
func NewParserFS(fs afero.Fs) Parser {
	return &DefaultParser{
		fs: fs,
	}
}

// Exists check if preset exists
func (p *DefaultParser) Exists(preset string) (exists bool) {
	_, exists = p.Presets[preset]
	return
}

// GetLanguages get all presets languages
func (p *DefaultParser) GetLanguages() (languages []string) {
	if len(p.Presets) == 0 {
		return
	}

	var lookedLangs map[string]bool = make(map[string]bool)

	for preset := range p.Presets {
		config, err := p.GetConfig(preset)

		if err == nil && config != nil && config.Language != "" && !lookedLangs[config.Language] {
			languages = append(languages, config.Language)
			lookedLangs[config.Language] = true
		}
	}
	sort.Strings(languages)
	return
}

// GetPresets get all presets names
func (p *DefaultParser) GetPresets(language string) (presets []string) {
	if len(p.Presets) == 0 {
		return
	}

	for preset := range p.Presets {
		config, err := p.GetConfig(preset)

		if err != nil || config == nil {
			presets = append(presets, preset)
		} else if config.Language == language {
			presets = append(presets, preset)
		}
	}
	sort.Strings(presets)
	return
}

// ErrPresetWriteAllBytes error throwed when did not write all preset file bytes
var ErrPresetWriteAllBytes = errors.New("failed to write all bytes")

// LookUpFiles check if preset files exist
func (p *DefaultParser) LookUpFiles(preset string) (foundFiles []string) {
	presetFiles := p.Presets[preset]

	for fileName := range presetFiles {
		if _, err := p.fs.Stat(fileName); !os.IsNotExist(err) {
			foundFiles = append(foundFiles, fileName)
		}
	}
	return
}

// WriteFile write preset file
func (p *DefaultParser) WriteFile(fileName string, fileContent string) (fileError string, err error) {
	var (
		file  afero.File
		lines int
	)

	file, err = p.fs.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)

	if err != nil {
		fileError = fileName
		return
	}

	defer file.Close()

	if lines, err = file.Write([]byte(fileContent)); err != nil {
		fileError = fileName
		return
	}

	if len([]byte(fileContent)) != lines {
		fileError = fileName
		err = ErrPresetWriteAllBytes
		return
	}

	if err = file.Sync(); err != nil {
		fileError = fileName
		return
	}

	return
}

// GetPresetKeys get preset file contents
func (p *DefaultParser) GetPresetKeys(preset string) (keys []string) {
	presetData := p.Presets[preset]

	for dataKey := range presetData {
		keys = append(keys, dataKey)
	}

	sort.Strings(keys)

	return
}

// GetPresetKeyContent get preset key value
func (p *DefaultParser) GetPresetKeyContent(preset string, key string) (value string) {
	presetData := p.Presets[preset]

	for dataKey, dataContent := range presetData {
		if dataKey == key {
			value = dataContent
			return
		}
	}

	return
}

// GetTemplates get all templates
func (p *DefaultParser) GetTemplates() map[string]map[string]string {
	return p.Templates
}

// LoadPresets loads the presets
func (p *DefaultParser) LoadPresets(allPresets map[string]map[string]string) {
	p.Presets = allPresets
}

// LoadTemplates loads the templates
func (p *DefaultParser) LoadTemplates(allTemplates map[string]map[string]string) {
	p.Templates = allTemplates
}

// LoadConfigs load the configs
func (p *DefaultParser) LoadConfigs(allConfigs map[string]string) {
	p.Configs = allConfigs
}

// GetConfig get preset config
func (p *DefaultParser) GetConfig(preset string) (config *PresetConfig, err error) {
	var (
		configValue string
		hasConfig   bool
	)

	if configValue, hasConfig = p.Configs[preset]; !hasConfig {
		err = fmt.Errorf("configuration for preset %s not found", preset)
		return
	}

	config = new(PresetConfig)
	err = yaml.Unmarshal([]byte(configValue), config)
	return
}
