package presets

import (
	"errors"
	"os"
	"sort"
	"strings"
)

var osStat func(string) (os.FileInfo, error) = os.Stat

// DefaultParser holds presets parsing data
type DefaultParser struct {
	Presets map[string]map[string]string
}

// Parser holds presets parsing logic
type Parser interface {
	Exists(string) bool
	GetLanguages() []string
	GetPresets(string) []string
	LookUpFiles(string) []string
	WriteFile(string, string) (string, error)
	GetPresetKeys(string) []string
	GetPresetKeyContent(string, string) string
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
	for _, content := range p.Presets {
		if presetLang, ok := content["preset_language"]; ok && !lookedLangs[presetLang] {
			languages = append(languages, presetLang)
			lookedLangs[presetLang] = true
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

	for key, content := range p.Presets {
		if language == "" {
			presets = append(presets, key)
		} else if presetLang, ok := content["preset_language"]; ok && presetLang == language {
			presets = append(presets, key)
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
		if strings.HasPrefix(fileName, "preset_") {
			continue
		}

		if _, err := osStat(fileName); !os.IsNotExist(err) {
			foundFiles = append(foundFiles, fileName)
		}
	}
	return
}

// WriteFile write preset file
func (p *DefaultParser) WriteFile(fileName string, fileContent string) (fileError string, err error) {
	var (
		file  *os.File
		lines int
	)

	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)

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
