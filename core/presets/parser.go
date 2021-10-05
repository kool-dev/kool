package presets

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"kool-dev/kool/core/automate"
	"kool-dev/kool/core/shell"
	"os"
	"sort"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

const presetConfigFile = "presets/%s/config.yml"

var source embed.FS

// SetSource informs the package about the
// source of template files and configs
func SetSource(src embed.FS) {
	source = src
}

// DefaultParser holds presets parsing data
type DefaultParser struct {
	local afero.Fs
}

// Parser holds presets parsing logic
type Parser interface {
	Exists(string) bool
	GetTags() []string
	GetPresets(string) []string
	LookUpFiles(string) []string
	Install(string, shell.OutputWritter) error
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

// GetPresets get all presets names
func (p *DefaultParser) GetPresets(tag string) (presets []string) {
	var (
		entries []fs.DirEntry
		folder  fs.DirEntry
		data    []byte
		config  = new(PresetConfig)
	)

	entries, _ = source.ReadDir("presets")

	for _, folder = range entries {
		data, _ = source.ReadFile(
			fmt.Sprintf(presetConfigFile, folder.Name()),
		)

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

// LookUpFiles check if preset files exist
func (p *DefaultParser) LookUpFiles(preset string) (foundFiles []string) {
	for _, fileName := range p.presetFiles(preset) {
		if _, err := p.local.Stat(fileName); !os.IsNotExist(err) {
			foundFiles = append(foundFiles, fileName)
		}
	}
	return
}

// Install write preset files
func (p *DefaultParser) Install(preset string, output shell.OutputWritter) (err error) {
	var (
		config *PresetConfig
		step   *automate.ActionStep
		action *automate.Action
	)

	if config, err = p.getConfig(preset); err != nil {
		return
	}

	for _, step = range config.Preset {
		output.Println("-", step.Name)

		for _, action = range step.Actions {
			switch action.Type() {
			case automate.TypeAdd:
				// the 'add' operation will run a new recipe
				// that is composed by a new array of ActionStep
				// action.Recipe
				output.Println("add:", action.Recipe)
				break
			case automate.TypeCopy:
				// io.Copy(dst, src)
				// action.Src
				// action.Dst
				output.Println("copy:", action.Src, action.Dst)
				break
			case automate.TypeScripts:
				// action.Scripts
				output.Println("scripts:", len(action.Scripts))
				break
			case automate.TypePrompt:
				// action.Prompt
				output.Println("prompt:", action.Prompt)
				break
			default:
				err = fmt.Errorf("ops, something is wrong with this preset config (%d)", action.Type())
				return
			}
		}
	}

	// for _, fileName := range p.presetFiles(preset) {
	// 	fileContent, _ = source.ReadFile(
	// 		fmt.Sprintf("presets/%s/%s", preset, fileName),
	// 	)

	// 	var (
	// 		file afero.File
	// 		size int
	// 	)

	// 	if _, statErr := p.local.Stat(fileName); !os.IsNotExist(statErr) {
	// 		if err = p.local.Rename(fileName, fmt.Sprintf("%s.bak.%s", fileName, time.Now().Format("20060102"))); err != nil {
	// 			fileError = fileName
	// 			return
	// 		}
	// 	}

	// 	file, err = p.local.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)

	// 	if err != nil {
	// 		fileError = fileName
	// 		return
	// 	}

	// 	if size, err = file.Write(fileContent); err != nil {
	// 		fileError = fileName
	// 		return
	// 	}

	// 	if len(fileContent) != size {
	// 		fileError = fileName
	// 		err = ErrPresetWriteAllBytes
	// 		return
	// 	}

	// 	if err = file.Sync(); err != nil {
	// 		fileError = fileName
	// 		return
	// 	}

	// 	file.Close()
	// }

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
