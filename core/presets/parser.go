package presets

import (
	"errors"
	"fmt"
	"io/fs"
	"kool-dev/kool/core/automate"
	"kool-dev/kool/core/shell"
	"sort"

	"gopkg.in/yaml.v2"
)

const presetConfigFile = "presets/%s/config.yml"

// SourceFS componds all required interfaces for managing
// the sourcing of presets and templates on a filesystem
type SourceFS interface {
	fs.FS
	fs.ReadDirFS
	fs.ReadFileFS
}

var source SourceFS

// SetSource informs the package about the
// source of template files and configs
func SetSource(src SourceFS) {
	source = src
}

// DefaultParser holds presets parsing data
type DefaultParser struct {
	presetID string
}

// Parser holds presets parsing logic
type Parser interface {
	Exists(string) bool
	GetTags() []string
	GetPresets(string) []string
	Install(string, shell.Shell) error
	Create(string, shell.Shell) error
	Add(string, shell.Shell) error
}

// NewParser creates a new preset default parser
func NewParser() Parser {
	return &DefaultParser{}
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

// GetTags get all presets tags
func (p *DefaultParser) GetTags() (tags []string) {
	var (
		entries []fs.DirEntry
		folder  fs.DirEntry
		data    []byte
		config  = new(PresetConfig)
		exists  = make(map[string]bool)
	)

	entries, _ = source.ReadDir("presets")

	for _, folder = range entries {
		data, _ = source.ReadFile(
			fmt.Sprintf(presetConfigFile, folder.Name()),
		)

		_ = yaml.Unmarshal(data, config)

		for _, tag := range config.Tags {
			if !exists[tag] {
				tags = append(tags, tag)
				exists[tag] = true
			}
		}
	}
	sort.Strings(tags)
	return
}

// GetPresets look up all presets IDs with the given tag
func (p *DefaultParser) GetPresets(tag string) (presets []string) {
	var (
		entries []fs.DirEntry
		folder  fs.DirEntry
		data    []byte
		config  *PresetConfig
	)

	entries, _ = source.ReadDir("presets")

	for _, folder = range entries {
		data, _ = source.ReadFile(
			fmt.Sprintf(presetConfigFile, folder.Name()),
		)

		config = new(PresetConfig)
		_ = yaml.Unmarshal(data, config)

		if config.HasTag(tag) {
			presets = append(presets, folder.Name())
		}
	}

	sort.Strings(presets)
	return
}

// ErrPresetWriteAllBytes error throwed when did not write all preset file bytes
var ErrPresetWriteAllBytes = errors.New("failed to write all bytes")

// Install executes the preset installation actions
func (p *DefaultParser) Install(preset string, sh shell.Shell) (err error) {
	var (
		config *PresetConfig
	)

	if config, err = p.getConfig(preset); err != nil {
		return
	}

	p.presetID = preset
	if err = automate.NewExecutor(sh, p.getSourceFile).Do(config.Preset); err != nil {
		return
	}

	return
}

// Create executes the preset installation actions
func (p *DefaultParser) Create(preset string, sh shell.Shell) (err error) {
	var (
		config *PresetConfig
	)

	if config, err = p.getConfig(preset); err != nil {
		return
	}

	p.presetID = preset
	if err = automate.NewExecutor(sh, p.getSourceFile).Do(config.Create); err != nil {
		return
	}

	return
}

func (p *DefaultParser) Add(recipe string, sh shell.Shell) (err error) {
	var steps = []*automate.ActionSet{
		{
			Name: fmt.Sprintf("Running recipe %s", recipe),
			Actions: []*automate.Action{
				{
					Recipe: recipe,
				},
			},
		},
	}

	if err = automate.NewExecutor(sh, p.getSourceFile).Do(steps); err != nil {
		return
	}

	return
}

func (p *DefaultParser) getSourceFile(path string) (data []byte, err error) {
	if p.presetID != "" {
		// look up in the preset folder
		if data, err = source.ReadFile(fmt.Sprintf("presets/%s/%s", p.presetID, path)); err == nil {
			return
		}
	}

	// fallback looking at the global templates
	if data, err = source.ReadFile(fmt.Sprintf("templates/%s", path)); err != nil {
		err = fmt.Errorf("could not find %s on within preset or global templates (err: %v)", path, err)
	}

	return
}

// getConfig parses the preset config data for usage
func (p *DefaultParser) getConfig(preset string) (config *PresetConfig, err error) {
	var data []byte

	data, err = source.ReadFile(
		fmt.Sprintf(presetConfigFile, preset),
	)

	if err != nil {
		err = fmt.Errorf("configuration for preset %s not found (%v)", preset, err)
		return
	}

	config = new(PresetConfig)
	err = yaml.Unmarshal(data, config)
	config.presetID = preset
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
