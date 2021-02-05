// +build !windows

package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/parser"
	"kool-dev/kool/cmd/shell"
	"kool-dev/kool/environment"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newFakeKoolRun(mockParsedCommands map[string][]builder.Command, mockParseError map[string]error) *KoolRun {
	return &KoolRun{
		*newFakeKoolService(),
		&parser.FakeParser{MockParsedCommands: mockParsedCommands, MockParseError: mockParseError},
		environment.NewFakeEnvStorage(),
		&shell.FakePromptSelect{},
		[]builder.Command{},
	}
}

func TestNewKoolRun(t *testing.T) {
	k := NewKoolRun()

	if _, ok := k.DefaultKoolService.shell.(*shell.DefaultShell); !ok {
		t.Errorf("unexpected shell.Shell on default KoolRun instance")
	}

	if _, ok := k.DefaultKoolService.exiter.(*shell.DefaultExiter); !ok {
		t.Errorf("unexpected shell.Exiter on default KoolRun instance")
	}

	if _, ok := k.DefaultKoolService.term.(*shell.DefaultTerminalChecker); !ok {
		t.Errorf("unexpected shell.TerminalChecker on default KoolRun instance")
	}

	if _, ok := k.parser.(*parser.DefaultParser); !ok {
		t.Errorf("unexpected parser.Parser on default KoolRun instance")
	}
}

func TestNewRunCommand(t *testing.T) {
	fakeParsedCommands := map[string][]builder.Command{
		"script": {
			&builder.FakeCommand{MockCmd: "cmd1"},
		},
	}

	f := newFakeKoolRun(fakeParsedCommands, nil)
	cmd := NewRunCommand(f)

	cmd.SetArgs([]string{"script"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
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

		if val, ok := f.shell.(*shell.FakeShell).CalledInteractive[command.Cmd()]; !ok || !val {
			t.Errorf("parsed command did not call Interactive")
		}
	}
}

func TestNewRunCommandMultipleScriptsWarning(t *testing.T) {
	f := newFakeKoolRun(nil, map[string]error{"script": parser.ErrMultipleDefinedScript})
	cmd := NewRunCommand(f)

	cmd.SetArgs([]string{"script"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledWarning {
		t.Errorf("did not call Warning for multiple scripts")
	}

	expectedWarning := "Attention: the script was found in more than one kool.yml file"

	if gotWarning := fmt.Sprint(f.shell.(*shell.FakeShell).WarningOutput...); gotWarning != expectedWarning {
		t.Errorf("expecting warning '%s', got '%s'", expectedWarning, gotWarning)
	}
}

func TestNewRunCommandParseError(t *testing.T) {
	f := newFakeKoolRun(nil, map[string]error{"script": errors.New("parse error")})
	cmd := NewRunCommand(f)

	cmd.SetArgs([]string{"script"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error for parse error")
	}

	expectedError := "parse error"

	if gotError := f.shell.(*shell.FakeShell).Err.Error(); gotError != expectedError {
		t.Errorf("expecting error '%s', got '%s'", expectedError, gotError)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("got an parse error, but command did not exit")
	}
}

func TestNewRunCommandExtraArgsError(t *testing.T) {
	fakeParsedCommands := map[string][]builder.Command{
		"script": {
			&builder.FakeCommand{},
			&builder.FakeCommand{},
		},
	}
	f := newFakeKoolRun(fakeParsedCommands, nil)
	cmd := NewRunCommand(f)

	cmd.SetArgs([]string{"script", "extraArg"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error for extra arguments")
	}

	expectedError := ErrExtraArguments.Error()

	if gotError := f.shell.(*shell.FakeShell).Err.Error(); gotError != expectedError {
		t.Errorf("expecting error '%s', got '%s'", expectedError, gotError)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("got an extra arguments error, but command did not exit")
	}
}

func TestNewRunCommandErrorInteractive(t *testing.T) {
	fakeParsedCommands := map[string][]builder.Command{
		"script": {
			&builder.FakeCommand{MockInteractiveError: errors.New("interactive error")},
		},
	}
	f := newFakeKoolRun(fakeParsedCommands, nil)
	cmd := NewRunCommand(f)

	cmd.SetArgs([]string{"script"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error for parsed command failure")
	}

	expectedError := "interactive error"

	if gotError := f.shell.(*shell.FakeShell).Err.Error(); gotError != expectedError {
		t.Errorf("expecting error '%s', got '%s'", expectedError, gotError)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("got an error executing parsed command, but command did not exit")
	}
}

func TestNewRunCommandScriptNotFound(t *testing.T) {
	f := newFakeKoolRun(nil, nil)
	cmd := NewRunCommand(f)

	cmd.SetArgs([]string{"script"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error for not found script error")
	}

	expectedError := ErrKoolScriptNotFound.Error()

	if gotError := f.shell.(*shell.FakeShell).Err.Error(); gotError != expectedError {
		t.Errorf("expecting error '%s', got '%s'", expectedError, gotError)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("got an not found script error, but command did not exit")
	}
}

func TestNewRunCommandWithArguments(t *testing.T) {
	fakeParsedCommands := map[string][]builder.Command{
		"script": {
			&builder.FakeCommand{},
		},
	}
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
	f := newFakeKoolRun(nil, nil)
	f.parser.(*parser.FakeParser).MockScripts = []string{"testing_script"}
	cmd := NewRunCommand(f)
	SetRunUsageFunc(f, cmd)

	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledPrintln {
		t.Error("did not call Println for command usage")
	}

	usage := f.shell.(*shell.FakeShell).OutLines[0]

	if !strings.Contains(usage, "testing_script") {
		t.Error("did not find testing_script as available script on usage text")
	}
}

func TestNewRunCommandFailingUsageTemplate(t *testing.T) {
	f := newFakeKoolRun(nil, nil)
	f.parser.(*parser.FakeParser).MockScripts = []string{"testing_script"}
	f.parser.(*parser.FakeParser).MockParseAvailableScriptsError = errors.New("error parse avaliable scripts")
	f.env.(*environment.FakeEnvStorage).Envs["KOOL_VERBOSE"] = "1"

	cmd := NewRunCommand(f)
	SetRunUsageFunc(f, cmd)

	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	output := f.shell.(*shell.FakeShell).OutLines[0]

	if strings.Contains(output, "testing_script") {
		t.Error("should not find testing_script as available script on usage text due to error on parsing scripts")
	}

	if !f.shell.(*shell.FakeShell).CalledPrintln {
		t.Error("did not call Println to output error on getting available scripts when KOOL_VERBOSE is true")
	}

	expected := "$ got an error trying to add available scripts to command usage template; error: error parse avaliable scripts"

	if expected != output {
		t.Errorf("expecting message '%s', got '%s'", expected, output)
	}
}

func TestNewRunCommandCompletion(t *testing.T) {
	var scripts []string
	f := newFakeKoolRun(nil, nil)
	f.parser.(*parser.FakeParser).MockScripts = []string{"testing_script"}
	cmd := NewRunCommand(f)

	scripts, _ = cmd.ValidArgsFunction(cmd, []string{}, "")

	if len(scripts) != 1 || scripts[0] != "testing_script" {
		t.Errorf("expecting suggestions [testing_script], got %v", scripts)
	}

	scripts, _ = cmd.ValidArgsFunction(cmd, []string{}, "tes")

	if len(scripts) != 1 || scripts[0] != "testing_script" {
		t.Errorf("expecting suggestions [testing_script], got %v", scripts)
	}

	scripts, _ = cmd.ValidArgsFunction(cmd, []string{}, "invalid")

	if len(scripts) != 0 {
		t.Errorf("expecting no suggestion, got %v", scripts)
	}

	scripts, _ = cmd.ValidArgsFunction(cmd, []string{"testing_script"}, "")

	if scripts != nil {
		t.Errorf("expecting no suggestion, got %v", scripts)
	}
}

func TestNewRunCommandFailingCompletion(t *testing.T) {
	var scripts []string
	f := newFakeKoolRun(nil, nil)
	f.parser.(*parser.FakeParser).MockScripts = []string{"testing_script"}
	f.parser.(*parser.FakeParser).MockParseAvailableScriptsError = errors.New("parsing error")
	cmd := NewRunCommand(f)

	scripts, _ = cmd.ValidArgsFunction(cmd, []string{}, "")

	if scripts != nil {
		t.Errorf("expecting no suggestion, got %v", scripts)
	}
}

func TestRunRecursiveCalls(t *testing.T) {
	makeKoolRoot := func() *cobra.Command {
		k := NewKoolRun()
		k.env = environment.NewFakeEnvStorage()
		k.env.Set("HOME", "")
		tmp := t.TempDir()
		k.env.Set("PWD", tmp)

		kooYml := []byte(`scripts:
  show-version: kool -v
  recursive: kool run show-version
  recursive:multi:
    - kool run recursive
    - kool run recursive
`)

		if err := ioutil.WriteFile(fmt.Sprintf("%s/kool.yml", tmp), kooYml, os.ModePerm); err != nil {
			t.Fatalf("failed creating temp kool.yml for testing: %v", err)
		}

		root := NewRootCmd(k.env)
		root.AddCommand(NewRunCommand(k))

		shell.RecursiveCall = func(args []string, in io.Reader, out, err io.Writer) error {
			fmt.Printf("called RecursiveCall args: %v\n", args)
			root.SetArgs(args)
			return root.Execute()
		}

		return root
	}

	defer func() {
		// clear up shell.RecursiveCall
		shell.RecursiveCall = nil
	}()

	root := makeKoolRoot()

	root.SetArgs([]string{"run", "show-version"})

	if err := root.Execute(); err != nil {
		t.Errorf("unexpected error executing run show-version; error: %v", err)
	}

	root = makeKoolRoot()
	root.SetArgs([]string{"run", "recursive"})

	if err := root.Execute(); err != nil {
		t.Errorf("unexpected error executing run recursive; error: %v", err)
	}

	root = makeKoolRoot()
	root.SetArgs([]string{"run", "recursive:multi"})

	if err := root.Execute(); err != nil {
		t.Errorf("unexpected error executing run recursive; error: %v", err)
	}
}

const inputContent string = "input file"

func TestRunRecursiveCallsWithInputRedirecting(t *testing.T) {
	makeKoolRoot := func() *cobra.Command {
		k := NewKoolRun()
		k.env = environment.NewFakeEnvStorage()
		k.env.Set("HOME", "")
		tmp := t.TempDir()
		inputFilePath := fmt.Sprintf("%s/input_file", tmp)
		k.env.Set("PWD", tmp)

		kooYml := []byte(fmt.Sprintf(`scripts:
  input: kool receive-file < %s
`, inputFilePath))

		if err := ioutil.WriteFile(fmt.Sprintf("%s/kool.yml", tmp), kooYml, os.ModePerm); err != nil {
			t.Fatalf("failed creating temp kool.yml for testing: %v", err)
		}
		if err := ioutil.WriteFile(inputFilePath, []byte(inputContent), os.ModePerm); err != nil {
			t.Fatalf("failed creating temp %v for testing: %v", inputFilePath, err)
		}

		root := NewRootCmd(k.env)
		root.AddCommand(NewRunCommand(k))
		root.AddCommand(&cobra.Command{
			Use: "receive-file",
			Run: func(cmd *cobra.Command, args []string) {
				if shell.NewTerminalChecker().IsTerminal(cmd.InOrStdin()) {
					t.Errorf("unexpected input - TTY - %T", cmd.InOrStdin())
				}
				if file, isFile := cmd.InOrStdin().(*os.File); !isFile {
					t.Errorf("unexpected input - should be a file; but is %T", cmd.InOrStdin())
				} else if input, err := ioutil.ReadAll(file); err != nil {
					t.Errorf("failed reading input file: %v", err)
				} else if string(input) != inputContent {
					t.Errorf("unexpected content on input file: %v", input)
				}
			},
		})

		setRecursiveCall(root)

		return root
	}

	defer func() {
		// clear up shell.RecursiveCall
		shell.RecursiveCall = nil
	}()

	root := makeKoolRoot()

	root.SetArgs([]string{"run", "input"})

	if err := root.Execute(); err != nil {
		t.Errorf("unexpected error executing run show-version; error: %v", err)
	}
}

func TestNewRunCommandWithTypoError(t *testing.T) {
	fakeParsedCommands := map[string][]builder.Command{
		"script1": {
			&builder.FakeCommand{MockCmd: "cmd"},
		},
	}

	possibleTypoError := &parser.ErrPossibleTypo{}
	possibleTypoError.SetSimilars([]string{"script1"})

	fakeParsedError := map[string]error{
		"script": possibleTypoError,
	}
	f := newFakeKoolRun(fakeParsedCommands, fakeParsedError)

	f.promptSelect.(*shell.FakePromptSelect).MockAnswer = map[string]string{
		"did you mean 'script1'?": "Yes",
	}

	cmd := NewRunCommand(f)

	cmd.SetArgs([]string{"script"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.promptSelect.(*shell.FakePromptSelect).CalledAsk {
		t.Error("did not call Ask on PromptSelect")
	}

	if val, ok := f.shell.(*shell.FakeShell).CalledInteractive["cmd"]; !ok || !val {
		t.Errorf("did not call Interactive for command 'cmd'")
	}
}

func TestNewRunCommandWithTypoErrorMultipleSimilar(t *testing.T) {
	fakeParsedCommands := map[string][]builder.Command{
		"script2": {
			&builder.FakeCommand{MockCmd: "cmd"},
		},
	}

	possibleTypoError := &parser.ErrPossibleTypo{}
	possibleTypoError.SetSimilars([]string{"script1", "script2"})

	fakeParsedError := map[string]error{
		"script": possibleTypoError,
	}

	f := newFakeKoolRun(fakeParsedCommands, fakeParsedError)

	f.promptSelect.(*shell.FakePromptSelect).MockAnswer = map[string]string{
		"did you mean one of ['script1', 'script2']?": "Yes",
		"which one did you mean?":                     "script2",
	}

	cmd := NewRunCommand(f)

	cmd.SetArgs([]string{"script"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.promptSelect.(*shell.FakePromptSelect).CalledAsk {
		t.Error("did not call Ask on PromptSelect")
	}

	if val, ok := f.shell.(*shell.FakeShell).CalledInteractive["cmd"]; !ok || !val {
		t.Errorf("did not call Interactive for command 'cmd'")
	}
}

func TestNewRunCommandWithTypoErrorNoAnswer(t *testing.T) {
	possibleTypoError := &parser.ErrPossibleTypo{}
	possibleTypoError.SetSimilars([]string{"script1"})

	fakeParsedError := map[string]error{
		"script": possibleTypoError,
	}

	f := newFakeKoolRun(nil, fakeParsedError)

	f.promptSelect.(*shell.FakePromptSelect).MockAnswer = map[string]string{
		"did you mean 'script1'?": "No",
	}

	cmd := NewRunCommand(f)

	cmd.SetArgs([]string{"script"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing run command; error: %v", err)
	}

	if !f.promptSelect.(*shell.FakePromptSelect).CalledAsk {
		t.Error("did not call Ask on PromptSelect")
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error for not found script error")
	}

	expectedError := ErrKoolScriptNotFound.Error()

	if gotError := f.shell.(*shell.FakeShell).Err.Error(); gotError != expectedError {
		t.Errorf("expecting error '%s', got '%s'", expectedError, gotError)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("got an not found script error, but command did not exit")
	}
}
