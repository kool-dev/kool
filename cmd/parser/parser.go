package parser

import (
	"errors"
	"kool-dev/kool/cmd/builder"
	"os"
	"path"
)

// Parser defines the functions required for handling kool.yml files.
type Parser interface {
	AddLookupPath(string) error
	Parse(string) ([]builder.Command, error)
}

// DefaultParser implements all default behavior for using kool.yml files.
type DefaultParser struct {
	targetFiles []string
}

// NewParser initialises a Parser to be used for handling kool.yml scripts.
func NewParser() Parser {
	return &DefaultParser{}
}

// AddLookupPath adds a folder to look for kool.yml scripts file.
func (p *DefaultParser) AddLookupPath(rootPath string) (err error) {
	var koolFile string

	if _, err = os.Stat(path.Join(rootPath, "kool.yml")); err == nil {
		koolFile = path.Join(rootPath, "kool.yml")
	} else if _, err = os.Stat(path.Join(rootPath, "kool.yaml")); err == nil {
		koolFile = path.Join(rootPath, "kool.yaml")
	}

	if koolFile == "" {
		err = ErrKoolYmlNotFound
	} else {
		p.targetFiles = append(p.targetFiles, koolFile)
	}

	return
}

// Parse looks up for the given script name on all of the kool.yml files available
// on the configured lookup paths. If the script exists in more than one file
// this function will return the first occurrence and an ErrMultipleDefinedScript
// error just to let the user know and avoid confusing.
func (p *DefaultParser) Parse(script string) (commands []builder.Command, err error) {
	var (
		koolFile        string
		parsedFile      *KoolYaml
		found           bool
		previouslyFound bool
	)

	if len(p.targetFiles) == 0 {
		err = errors.New("kool.yml not found")
		return
	}

	for _, koolFile = range p.targetFiles {
		if parsedFile, err = ParseKoolYaml(koolFile); err != nil {
			return
		}

		if found = parsedFile.HasScript(script); found {
			if !previouslyFound {
				// this is the first time we find the script we want!
				previouslyFound = true

				if commands, err = parsedFile.ParseCommands(script); err != nil {
					return
				}
			} else {
				// so we already found once, and now found again the same script
				// in another file! let's warn about that
				err = ErrMultipleDefinedScript
			}
		}
	}
	return
}
