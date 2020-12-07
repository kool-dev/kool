package parser

import (
	"kool-dev/kool/cmd/builder"
	"strings"
)

// FakeParser implements all fake behaviors for using parser in tests.
type FakeParser struct {
	CalledAddLookupPath            bool
	TargetFiles                    []string
	CalledParse                    bool
	CalledParseAvailableScripts    bool
	CalledLookUpVariables          bool
	MockParsedCommands             []builder.Command
	MockParseError                 error
	MockScripts                    []string
	MockParseAvailableScriptsError error
	MockVariables                  []string
}

// AddLookupPath implements fake AddLookupPath behavior
func (f *FakeParser) AddLookupPath(rootPath string) (err error) {
	f.CalledAddLookupPath = true
	f.TargetFiles = append(f.TargetFiles, "kool.yml")
	return
}

// Parse implements fake Parse behavior
func (f *FakeParser) Parse(script string) (commands []builder.Command, err error) {
	f.CalledParse = true
	commands = f.MockParsedCommands
	err = f.MockParseError
	return
}

// ParseAvailableScripts implements fake ParseAvailableScripts behavior
func (f *FakeParser) ParseAvailableScripts(filter string) (scripts []string, err error) {
	f.CalledParseAvailableScripts = true

	if filter == "" {
		scripts = f.MockScripts
	} else {
		for _, script := range f.MockScripts {
			if strings.HasPrefix(script, filter) {
				scripts = append(scripts, script)
			}
		}
	}

	err = f.MockParseAvailableScriptsError
	return
}

// LookUpVariables implements fake LookUpVariables behavior
func (f *FakeParser) LookUpVariables(script string) (variables []string) {
	f.CalledLookUpVariables = true
	variables = f.MockVariables
	return
}
