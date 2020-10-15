package cmd

import (
	"errors"
	"fmt"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/parser"
	"kool-dev/kool/cmd/shell"
	"kool-dev/kool/environment"
	"strings"
	"testing"
)

func newFakeKoolRun(mockParsedCommands []builder.Command, mockParseError error) *KoolRun {
	return &KoolRun{
		*newFakeKoolService(),
		&parser.FakeParser{MockParsedCommands: mockParsedCommands, MockParseError: mockParseError},
		environment.NewFakeEnvStorage(),
		[]builder.Command{},
	}
}

func TestNewKoolRun(t *testing.T) {
	k := NewKoolRun()

	if _, ok := k.DefaultKoolService.out.(*shell.DefaultOutputWriter); !ok {
		t.Errorf("unexpected shell.OutputWriter on default KoolRun instance")
	}

	if _, ok := k.DefaultKoolService.exiter.(*shell.DefaultExiter); !ok {
		t.Errorf("unexpected shell.Exiter on default KoolRun instance")
	}

	if _, ok := k.DefaultKoolService.in.(*shell.DefaultInputReader); !ok {
		t.Errorf("unexpected shell.InputReader on default KoolRun instance")
	}

	if _, ok := k.parser.(*parser.DefaultParser); !ok {
		t.Errorf("unexpected parser.Parser on default KoolRun instance")
	}
}

func TestNewRunCommand(t *testing.T) {
	fakeParsedCommands := []builder.Command{&builder.FakeCommand{}}

	f := newFakeKoolRun(fakeParsedCommands, nil)
	cmd := NewRunCommand(f)

	cmd.SetArgs([]string{"script"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.out.(*shell.FakeOutputWriter).CalledSetWriter {
		t.Errorf("did not call SetWriter")
	}

	if !f.parser.(*parser.FakeParser).CalledAddLookupPath {
		t.Errorf("did not call AddLookupPath")
	}

	targetFiles := f.parser.(*parser.FakeParser).TargetFiles

	if len(targetFiles) != 2 {
		t.Errorf("did not call AddLookupPath twice (global and local)")
	}

	if !f.parser.(*parser.FakeParser).CalledParse {
		t.Errorf("did not call Parse")
	}

	if len(f.commands) != 1 {
		t.Errorf("did not parse the commands")
	}

	for _, command := range f.commands {
		if command.(*builder.FakeCommand).CalledAppendArgs {
			t.Errorf("unexpected AppendArgs call by parsed command")
		}

		if !command.(*builder.FakeCommand).CalledInteractive {
			t.Errorf("parsed command did not call Interactive")
		}
	}
}

func TestNewRunCommandMultipleScriptsWarning(t *testing.T) {
	f := newFakeKoolRun([]builder.Command{}, parser.ErrMultipleDefinedScript)
	cmd := NewRunCommand(f)

	cmd.SetArgs([]string{"script"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.out.(*shell.FakeOutputWriter).CalledWarning {
		t.Errorf("did not call Warning for multiple scripts")
	}

	expectedWarning := "Attention: the script was found in more than one kool.yml file"

	if gotWarning := fmt.Sprint(f.out.(*shell.FakeOutputWriter).WarningOutput...); gotWarning != expectedWarning {
		t.Errorf("expecting warning '%s', got '%s'", expectedWarning, gotWarning)
	}
}

func TestNewRunCommandParseError(t *testing.T) {
	f := newFakeKoolRun([]builder.Command{}, errors.New("parse error"))
	cmd := NewRunCommand(f)

	cmd.SetArgs([]string{"script"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.out.(*shell.FakeOutputWriter).CalledError {
		t.Error("did not call Error for parse error")
	}

	expectedError := "parse error"

	if gotError := f.out.(*shell.FakeOutputWriter).Err.Error(); gotError != expectedError {
		t.Errorf("expecting error '%s', got '%s'", expectedError, gotError)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("got an parse error, but command did not exit")
	}
}

func TestNewRunCommandExtraArgsError(t *testing.T) {
	fakeParsedCommands := []builder.Command{&builder.FakeCommand{}, &builder.FakeCommand{}}
	f := newFakeKoolRun(fakeParsedCommands, nil)
	cmd := NewRunCommand(f)

	cmd.SetArgs([]string{"script", "extraArg"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.out.(*shell.FakeOutputWriter).CalledError {
		t.Error("did not call Error for extra arguments")
	}

	expectedError := ErrExtraArguments.Error()

	if gotError := f.out.(*shell.FakeOutputWriter).Err.Error(); gotError != expectedError {
		t.Errorf("expecting error '%s', got '%s'", expectedError, gotError)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("got an extra arguments error, but command did not exit")
	}
}

func TestNewRunCommandErrorInteractive(t *testing.T) {
	f := newFakeKoolRun([]builder.Command{&builder.FakeCommand{MockError: errors.New("interactive error")}}, nil)
	cmd := NewRunCommand(f)

	cmd.SetArgs([]string{"script"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.out.(*shell.FakeOutputWriter).CalledError {
		t.Error("did not call Error for parsed command failure")
	}

	expectedError := "interactive error"

	if gotError := f.out.(*shell.FakeOutputWriter).Err.Error(); gotError != expectedError {
		t.Errorf("expecting error '%s', got '%s'", expectedError, gotError)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("got an error executing parsed command, but command did not exit")
	}
}

func TestNewRunCommandScriptNotFound(t *testing.T) {
	f := newFakeKoolRun([]builder.Command{}, nil)
	cmd := NewRunCommand(f)

	cmd.SetArgs([]string{"script"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.out.(*shell.FakeOutputWriter).CalledError {
		t.Error("did not call Error for not found script error")
	}

	expectedError := ErrKoolScriptNotFound.Error()

	if gotError := f.out.(*shell.FakeOutputWriter).Err.Error(); gotError != expectedError {
		t.Errorf("expecting error '%s', got '%s'", expectedError, gotError)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("got an not found script error, but command did not exit")
	}
}

func TestNewRunCommandWithArguments(t *testing.T) {
	fakeParsedCommands := []builder.Command{&builder.FakeCommand{}}
	f := newFakeKoolRun(fakeParsedCommands, nil)
	cmd := NewRunCommand(f)

	cmd.SetArgs([]string{"script", "arg1", "arg2"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.commands[0].(*builder.FakeCommand).CalledAppendArgs {
		t.Error("did not call AppendArgs for parsed command")
	}

	fakeCommandArgs := f.commands[0].(*builder.FakeCommand).ArgsAppend

	if len(fakeCommandArgs) != 2 || fakeCommandArgs[0] != "arg1" || fakeCommandArgs[1] != "arg2" {
		t.Error("did not call AppendArgs properly for parsed command")
	}
}

func TestNewRunCommandUsageTemplate(t *testing.T) {
	f := newFakeKoolRun([]builder.Command{}, nil)
	f.parser.(*parser.FakeParser).MockScripts = []string{"testing_script"}
	cmd := NewRunCommand(f)
	SetRunUsageFunc(f, cmd)

	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.out.(*shell.FakeOutputWriter).CalledPrintln {
		t.Error("did not call Println for command usage")
	}

	usage := fmt.Sprintln(f.out.(*shell.FakeOutputWriter).Out...)

	if !strings.Contains(usage, "testing_script") {
		t.Error("did not find testing_script as available script on usage text")
	}
}

func TestNewRunCommandFailingUsageTemplate(t *testing.T) {
	f := newFakeKoolRun([]builder.Command{}, nil)
	f.parser.(*parser.FakeParser).MockScripts = []string{"testing_script"}
	f.parser.(*parser.FakeParser).MockParseAvailableScriptsError = errors.New("error parse avaliable scripts")
	f.envStorage.(*environment.FakeEnvStorage).Envs["KOOL_VERBOSE"] = "1"

	cmd := NewRunCommand(f)
	SetRunUsageFunc(f, cmd)

	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	output := strings.TrimSpace(fmt.Sprintln(f.out.(*shell.FakeOutputWriter).Out...))

	if strings.Contains(output, "testing_script") {
		t.Error("should not find testing_script as available script on usage text due to error on parsing scripts")
	}

	if !f.out.(*shell.FakeOutputWriter).CalledPrintln {
		t.Error("did not call Println to output error on getting available scripts when KOOL_VERBOSE is true")
	}

	expected := "$ got an error trying to add available scripts to command usage template; error: error parse avaliable scripts"

	if expected != output {
		t.Errorf("expecting message '%s', got '%s'", expected, output)
	}
}
