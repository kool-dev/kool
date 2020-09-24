package parser

import (
	"kool-dev/kool/cmd/builder"
)

// FakeParser implements all fake behaviors for using parser in tests.
type FakeParser struct {
	CalledAddLookupPath bool
	TargetFiles         []string
	CalledParse         bool
	MockParsedCommands  []builder.Command
	MockParseError      error
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
