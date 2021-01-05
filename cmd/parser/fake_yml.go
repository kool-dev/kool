package parser

import "kool-dev/kool/cmd/builder"

// FakeKoolYaml implements all fake behaviors for using yml parser in tests.
type FakeKoolYaml struct {
	CalledParse, CalledHasScript, CalledParseCommands, CalledSetScript map[string]bool
	CalledString                                                       bool
	ScriptCommands                                                     map[string][]string
	MockParseError                                                     map[string]error
	MockHasScript                                                      map[string]bool
	MockCommands                                                       map[string][]builder.Command
	MockParseCommandsError                                             map[string]error
	MockString                                                         string
	MockStringError                                                    error
}

// Parse decodes the target kool.yml
func (f *FakeKoolYaml) Parse(filePath string) (err error) {
	if f.CalledParse == nil {
		f.CalledParse = make(map[string]bool)
	}
	f.CalledParse[filePath] = true
	err = f.MockParseError[filePath]
	return
}

// HasScript tells if the given script exists on this parsed YAML.
func (f *FakeKoolYaml) HasScript(script string) (has bool) {
	if f.CalledHasScript == nil {
		f.CalledHasScript = make(map[string]bool)
	}
	f.CalledHasScript[script] = true
	has = f.MockHasScript[script]
	return
}

// ParseCommands parsed the given script from kool.yml file onto a list
// of commands parsed.
func (f *FakeKoolYaml) ParseCommands(script string) (commands []builder.Command, err error) {
	if f.CalledParseCommands == nil {
		f.CalledParseCommands = make(map[string]bool)
	}

	f.CalledParseCommands[script] = true

	commands = f.MockCommands[script]
	err = f.MockParseCommandsError[script]
	return
}

// SetScript set script into kool yaml
func (f *FakeKoolYaml) SetScript(key string, commands []string) {
	if f.CalledSetScript == nil {
		f.CalledSetScript = make(map[string]bool)
	}

	if f.ScriptCommands == nil {
		f.ScriptCommands = make(map[string][]string)
	}

	f.CalledSetScript[key] = true
	f.ScriptCommands[key] = commands
}

// String returns docker-compose as string
func (f *FakeKoolYaml) String() (content string, err error) {
	f.CalledString = true
	content = f.MockString
	err = f.MockStringError
	return
}
