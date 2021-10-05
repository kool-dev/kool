package presets

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"time"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

var source embed.FS

func SetSource(src embed.FS) {
	source = src
}

// PresetConfigQuestion preset config question
type PresetConfigQuestion struct {
	Key           string                       `yaml:"key"`
	DefaultAnswer string                       `yaml:"default_answer"`
	Message       string                       `yaml:"message"`
	Options       []PresetConfigQuestionOption `yaml:"options"`
}

// PresetConfigQuestionOption preset config question option
type PresetConfigQuestionOption struct {
	Name     string `yaml:"name"`
	Template string `yaml:"template"`
}

// PresetConfigTemplate default templates for preset
type PresetConfigTemplate struct {
	Key      string `yaml:"key"`
	Template string `yaml:"template"`
}

// ConfigLanguage holds the field to parse
// language out of preset configuration
type ConfigLanguage struct {
	Language string `yaml:"language"`
}

// PresetConfig preset config
type PresetConfig struct {
	Language  string                            `yaml:"language"`
	Commands  map[string][]string               `yaml:"commands"`
	Questions map[string][]PresetConfigQuestion `yaml:"questions"`
	Templates []PresetConfigTemplate            `yaml:"templates"`
}

// DefaultParser holds presets parsing data
type DefaultParser struct {
	local afero.Fs
}

// Parser holds presets parsing logic
type Parser interface {
	Exists(string) bool
	GetLanguages() []string
	GetPresets(string) []string
	LookUpFiles(string) []string
	WriteFiles(string) (string, error)
	GetConfig(string) (*PresetConfig, error)
}

// NewParser creates a new preset default parser
func NewParser() Parser {
	return &DefaultParser{
		local: afero.NewOsFs(),
	}
}

// NewParserFS creates a new preset default parser with file system
func NewParserFS(fs afero.Fs) Parser {
	return &DefaultParser{
		local: fs,
	}
}

// Exists check if preset exists
func (p *DefaultParser) Exists(preset string) bool {
	var (
		err error
	)

	if _, err = source.ReadDir(fmt.Sprintf("presets/%s", preset)); err != nil {
		return false
	}

	return true
}

// GetLanguages get all presets languages
func (p *DefaultParser) GetLanguages() (languages []string) {
	var (
		entries     []fs.DirEntry
		folder      fs.DirEntry
		data        []byte
		lang        = new(ConfigLanguage)
		lookedLangs = make(map[string]bool)
	)

	entries, _ = source.ReadDir("presets")

	for _, folder = range entries {
		data, _ = source.ReadFile(
			fmt.Sprintf("presets/%s/preset-config.yml", folder.Name()),
		)

		_ = yaml.Unmarshal(data, lang)

		if !lookedLangs[lang.Language] {
			languages = append(languages, lang.Language)
			lookedLangs[lang.Language] = true
		}
	}
	sort.Strings(languages)
	return
}

// GetPresets get all presets names
func (p *DefaultParser) GetPresets(language string) (presets []string) {
	var (
		entries []fs.DirEntry
		folder  fs.DirEntry
		data    []byte
		lang    = new(ConfigLanguage)
	)

	entries, _ = source.ReadDir("presets")

	for _, folder = range entries {
		data, _ = source.ReadFile(
			fmt.Sprintf("presets/%s/preset-config.yml", folder.Name()),
		)

		_ = yaml.Unmarshal(data, lang)

		if lang.Language == language {
			presets = append(presets, folder.Name())
		}
	}
	sort.Strings(presets)
	return
}

// ErrPresetWriteAllBytes error throwed when did not write all preset file bytes
var ErrPresetWriteAllBytes = errors.New("failed to write all bytes")

// LookUpFiles check if preset files exist
func (p *DefaultParser) LookUpFiles(preset string) (foundFiles []string) {
	for _, fileName := range p.presetFiles(preset) {
		if _, err := p.local.Stat(fileName); !os.IsNotExist(err) {
			foundFiles = append(foundFiles, fileName)
		}
	}
	return
}

// WriteFiles write preset files
func (p *DefaultParser) WriteFiles(preset string) (fileError string, err error) {
	var (
		fileContent []byte
	)

	for _, fileName := range p.presetFiles(preset) {
		fileContent, _ = source.ReadFile(
			fmt.Sprintf("presets/%s/%s", preset, fileName),
		)

		var (
			file afero.File
			size int
		)

		if _, statErr := p.local.Stat(fileName); !os.IsNotExist(statErr) {
			if err = p.local.Rename(fileName, fmt.Sprintf("%s.bak.%s", fileName, time.Now().Format("20060102"))); err != nil {
				fileError = fileName
				return
			}
		}

		file, err = p.local.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)

		if err != nil {
			fileError = fileName
			return
		}

		if size, err = file.Write(fileContent); err != nil {
			fileError = fileName
			return
		}

		if len(fileContent) != size {
			fileError = fileName
			err = ErrPresetWriteAllBytes
			return
		}

		if err = file.Sync(); err != nil {
			fileError = fileName
			return
		}

		file.Close()
	}

	return
}

// GetConfig get preset config
func (p *DefaultParser) GetConfig(preset string) (config *PresetConfig, err error) {
	var data []byte

	data, err = source.ReadFile(
		fmt.Sprintf("presets/%s/preset-config.yml", preset),
	)

	if err != nil {
		err = fmt.Errorf("configuration for preset %s not found (%v)", preset, err)
		return
	}

	config = new(PresetConfig)
	err = yaml.Unmarshal(data, config)
	return
}

func (p *DefaultParser) presetFiles(preset string) (presetFiles []string) {
	var (
		entries []fs.DirEntry
		file    fs.DirEntry
		err     error
	)

	if entries, err = source.ReadDir(fmt.Sprintf("presets/%s", preset)); err != nil {
		return
	}

	for _, file = range entries {
		presetFiles = append(presetFiles, file.Name())
	}
	return
}
