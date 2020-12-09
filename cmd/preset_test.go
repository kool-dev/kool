package cmd

import (
	"errors"
	"fmt"
	"kool-dev/kool/cmd/compose"
	"kool-dev/kool/cmd/presets"
	"kool-dev/kool/cmd/shell"
	"testing"

	"gopkg.in/yaml.v2"
)

const defaultCompose string = `version: "3.7"
services:
  app:
    image: kooldev/php:7.4-nginx
    ports:
     - "${KOOL_APP_PORT:-80}:80"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:delegated
    #  - $HOME/.ssh:/home/kool/.ssh:delegated
    networks:
     - kool_local
     - kool_global
  database:
    image: mysql:8.0 # can change to: mysql:5.7
    command: --default-authentication-plugin=mysql_native_password # remove this line if you change to: mysql:5.7
    ports:
     - "${KOOL_DATABASE_PORT:-3306}:3306"
    environment:
      MYSQL_ROOT_PASSWORD: "${DB_PASSWORD:-rootpass}"
      MYSQL_DATABASE: "${DB_DATABASE:-database}"
      MYSQL_USER: "${DB_USERNAME:-user}"
      MYSQL_PASSWORD: "${DB_PASSWORD:-pass}"
      MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
    volumes:
     - database:/var/lib/mysql:delegated
    networks:
     - kool_local
  cache:
    image: redis:6-alpine
    volumes:
     - cache:/data:delegated
    networks:
     - kool_local

volumes:
  database:
  cache:

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"
`

const mysqlTemplate string = `image: mysql:8.0
command: --default-authentication-plugin=mysql_native_password
ports:
  - "${KOOL_DATABASE_PORT:-3306}:3306"
environment:
  MYSQL_ROOT_PASSWORD: "${DB_PASSWORD:-rootpass}"
  MYSQL_DATABASE: "${DB_DATABASE:-database}"
  MYSQL_USER: "${DB_USERNAME:-user}"
  MYSQL_PASSWORD: "${DB_PASSWORD:-pass}"
  MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
volumes:
 - database:/var/lib/mysql:delegated
networks:
 - kool_local
 `

func newFakeKoolPreset() *KoolPreset {
	return &KoolPreset{
		*newFakeKoolService(),
		&KoolPresetFlags{false},
		&presets.FakeParser{},
		&compose.FakeParser{},
		&shell.FakePromptSelect{},
	}
}

func TestNewKoolPreset(t *testing.T) {
	k := NewKoolPreset()

	if _, ok := k.DefaultKoolService.shell.(*shell.DefaultShell); !ok {
		t.Errorf("unexpected shell.Shell on default KoolPreset instance")
	}

	if _, ok := k.DefaultKoolService.exiter.(*shell.DefaultExiter); !ok {
		t.Errorf("unexpected shell.Exiter on default KoolPreset instance")
	}

	if _, ok := k.DefaultKoolService.term.(*shell.DefaultTerminalChecker); !ok {
		t.Errorf("unexpected shell.TerminalChecker on default KoolPreset instance")
	}

	if k.Flags == nil {
		t.Errorf("Flags not initialized on default KoolPreset instance")
	} else if k.Flags.Override {
		t.Errorf("bad default value for Override flag on default KoolPreset instance")
	}

	if _, ok := k.presetsParser.(*presets.DefaultParser); !ok {
		t.Errorf("unexpected presets.Parser on default KoolPreset instance")
	}

	if _, ok := k.composeParser.(*compose.DefaultParser); !ok {
		t.Errorf("unexpected compose.Parser on default KoolPreset instance")
	}

	if _, ok := k.promptSelect.(*shell.DefaultPromptSelect); !ok {
		t.Errorf("unexpected shell.PromptSelect on default KoolPreset instance")
	}

	if _, ok := k.DefaultKoolService.term.(*shell.DefaultTerminalChecker); !ok {
		t.Errorf("unexpected shell.TerminalChecker on default KoolPreset instance")
	}
}

func TestPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.presetsParser.(*presets.FakeParser).MockExists = true
	f.presetsParser.(*presets.FakeParser).MockPresetKeyContent = map[string]map[string]string{
		"laravel": map[string]string{
			"kool.yml": "kool.yml content",
		},
	}
	f.presetsParser.(*presets.FakeParser).MockConfig = map[string]*presets.PresetConfig{
		"laravel": &presets.PresetConfig{},
	}

	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.presetsParser.(*presets.FakeParser).CalledExists {
		t.Error("did not call parser.Exists")
	}

	if !f.presetsParser.(*presets.FakeParser).CalledLoadTemplates {
		t.Error("did not call parser.LoadTemplates")
	}

	if !f.presetsParser.(*presets.FakeParser).CalledLoadPresets {
		t.Error("did not call parser.LoadPresets")
	}

	if !f.presetsParser.(*presets.FakeParser).CalledLoadConfigs {
		t.Error("did not call parser.LoadConfigs")
	}

	if val, ok := f.presetsParser.(*presets.FakeParser).CalledGetConfig["laravel"]; !ok || !val {
		t.Error("did not call parser.GetConfig for preset 'laravel'")
	}

	if !f.shell.(*shell.FakeShell).CalledPrintln {
		t.Error("did not call Println")
	}

	expected := "Preset laravel is initializing!"
	output := f.shell.(*shell.FakeShell).OutLines[0]

	if expected != output {
		t.Errorf("Expecting message '%s', got '%s'", expected, output)
	}

	if !f.presetsParser.(*presets.FakeParser).CalledLookUpFiles {
		t.Error("did not call parser.LookUpFiles")
	}

	if !f.presetsParser.(*presets.FakeParser).CalledGetTemplates {
		t.Error("did not call parser.GetTemplates")
	}

	if val, ok := f.presetsParser.(*presets.FakeParser).CalledWriteFiles["laravel"]; !ok || !val {
		t.Error("failed calling parser.WriteFiles for preset 'laravel'")
	}

	if !f.shell.(*shell.FakeShell).CalledSuccess {
		t.Error("did not call Success")
	}

	expected = "Preset laravel initialized!"
	output = fmt.Sprint(f.shell.(*shell.FakeShell).SuccessOutput...)

	if expected != output {
		t.Errorf("Expecting success message '%s', got '%s'", expected, output)
	}
}

func TestInvalidScriptPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"invalid"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.presetsParser.(*presets.FakeParser).CalledExists {
		t.Error("did not call parser.Exists")
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error")
	}

	expected := "Unknown preset invalid"
	output := f.shell.(*shell.FakeShell).Err.Error()

	if expected != output {
		t.Errorf("expecting error '%s', got '%s'", expected, output)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Exit")
	}
}

func TestExistingFilesPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.presetsParser.(*presets.FakeParser).MockExists = true
	f.presetsParser.(*presets.FakeParser).MockFoundFiles = []string{"kool.yml"}
	f.presetsParser.(*presets.FakeParser).MockConfig = map[string]*presets.PresetConfig{
		"laravel": &presets.PresetConfig{},
	}
	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledWarning {
		t.Error("did not call Warning")
	}

	expected := "Some preset files already exist. In case you wanna override them, use --override."
	output := fmt.Sprint(f.shell.(*shell.FakeShell).WarningOutput...)

	if output != expected {
		t.Errorf("expecting message '%s', got '%s'", expected, output)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Exit")
	}
}

func TestOverrideFilesPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.presetsParser.(*presets.FakeParser).MockExists = true
	f.presetsParser.(*presets.FakeParser).MockFoundFiles = []string{"kool.yml"}
	f.presetsParser.(*presets.FakeParser).MockConfig = map[string]*presets.PresetConfig{
		"laravel": &presets.PresetConfig{},
	}

	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"--override", "laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if f.presetsParser.(*presets.FakeParser).CalledLookUpFiles {
		t.Error("unexpected existing files checking")
	}

	if f.shell.(*shell.FakeShell).CalledWarning {
		t.Error("unexpected existing files Warning")
	}

	if f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("unexpected program Exit")
	}

	if !f.shell.(*shell.FakeShell).CalledSuccess {
		t.Error("did not call Success")
	}
}

func TestWriteErrorPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.presetsParser.(*presets.FakeParser).MockExists = true
	f.presetsParser.(*presets.FakeParser).MockPresetKeyContent = map[string]map[string]string{
		"laravel": map[string]string{
			"kool.yml": "kool.yml content",
		},
	}
	f.presetsParser.(*presets.FakeParser).MockError = errors.New("write error")
	f.presetsParser.(*presets.FakeParser).MockConfig = map[string]*presets.PresetConfig{
		"laravel": &presets.PresetConfig{},
	}

	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error")
	}

	expected := "Failed to write preset file : write error"
	output := f.shell.(*shell.FakeShell).Err.Error()

	if output != expected {
		t.Errorf("expecting error '%s', got '%s'", expected, output)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Exit")
	}
}

func TestNoArgsPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()

	mockAnswer := make(map[string]string)
	mockAnswer["What language do you want to use"] = "php"
	mockAnswer["What preset do you want to use"] = "laravel"

	f.promptSelect.(*shell.FakePromptSelect).MockAnswer = mockAnswer
	f.presetsParser.(*presets.FakeParser).MockLanguages = []string{"php"}
	f.presetsParser.(*presets.FakeParser).MockPresets = []string{"laravel"}
	f.presetsParser.(*presets.FakeParser).MockExists = true
	f.presetsParser.(*presets.FakeParser).MockConfig = map[string]*presets.PresetConfig{
		"laravel": &presets.PresetConfig{},
	}

	cmd := NewPresetCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.promptSelect.(*shell.FakePromptSelect).CalledAsk {
		t.Error("did not call Ask on PromptSelect")
	}

	expected := "Preset laravel is initializing!"
	output := f.shell.(*shell.FakeShell).OutLines[0]

	if expected != output {
		t.Errorf("Expecting message '%s', got '%s'", expected, output)
	}
}

func TestFailingLanguageNoArgsPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.presetsParser.(*presets.FakeParser).MockLanguages = []string{"php"}
	f.presetsParser.(*presets.FakeParser).MockPresets = []string{"laravel"}
	f.presetsParser.(*presets.FakeParser).MockConfig = map[string]*presets.PresetConfig{
		"laravel": &presets.PresetConfig{},
	}

	mockError := make(map[string]error)
	mockError["What language do you want to use"] = errors.New("error prompt select language")

	f.promptSelect.(*shell.FakePromptSelect).MockError = mockError

	cmd := NewPresetCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.promptSelect.(*shell.FakePromptSelect).CalledAsk {
		t.Error("did not call Ask on PromptSelect")
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error")
	}

	expected := "error prompt select language"
	output := f.shell.(*shell.FakeShell).Err.Error()

	if output != expected {
		t.Errorf("expecting error '%s', got '%s'", expected, output)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Exit")
	}
}

func TestFailingPresetNoArgsPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.presetsParser.(*presets.FakeParser).MockLanguages = []string{"php"}
	f.presetsParser.(*presets.FakeParser).MockPresets = []string{"laravel"}
	f.presetsParser.(*presets.FakeParser).MockConfig = map[string]*presets.PresetConfig{
		"laravel": &presets.PresetConfig{},
	}

	mockAnswer := make(map[string]string)
	mockAnswer["What language do you want to use"] = "php"

	f.promptSelect.(*shell.FakePromptSelect).MockAnswer = mockAnswer

	mockError := make(map[string]error)
	mockError["What preset do you want to use"] = errors.New("error prompt select preset")

	f.promptSelect.(*shell.FakePromptSelect).MockError = mockError

	cmd := NewPresetCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.promptSelect.(*shell.FakePromptSelect).CalledAsk {
		t.Error("did not call Ask on PromptSelect")
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error")
	}

	expected := "error prompt select preset"
	output := f.shell.(*shell.FakeShell).Err.Error()

	if output != expected {
		t.Errorf("expecting error '%s', got '%s'", expected, output)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Exit")
	}
}

func TestCancellingPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()

	mockError := make(map[string]error)
	mockError["What language do you want to use"] = shell.ErrPromptSelectInterrupted

	f.promptSelect.(*shell.FakePromptSelect).MockError = mockError

	cmd := NewPresetCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledWarning {
		t.Error("did not call Warning")
	}

	expected := "Operation Cancelled\n"
	output := fmt.Sprintln(f.shell.(*shell.FakeShell).WarningOutput...)

	if output != expected {
		t.Errorf("expecting warning '%s', got '%s'", expected, output)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Exit")
	}

	if f.exiter.(*shell.FakeExiter).Code() != 0 {
		t.Error("did not call Exit with code 0")
	}
}

func TestNonTTYPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.term.(*shell.FakeTerminalChecker).MockIsTerminal = false

	cmd := NewPresetCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error")
	}

	err := f.shell.(*shell.FakeShell).Err

	if err == nil {
		t.Error("expecting an error, got none")
	} else if err.Error() != "the input device is not a TTY; for non-tty environments, please specify a preset argument" {
		t.Errorf("expecting error 'the input device is not a TTY; for non-tty environments, please specify a preset argument', got %v", err)
	}
}

func TestCustomDockerComposePresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.presetsParser.(*presets.FakeParser).MockExists = true

	config := &presets.PresetConfig{
		Questions: map[string]presets.PresetConfigQuestion{
			"database": presets.PresetConfigQuestion{
				Message: "What database service do you want to use",
				Options: yaml.MapSlice{
					yaml.MapItem{Key: "mysql", Value: "mysql"},
					yaml.MapItem{Key: "postgresql", Value: "postgresql"},
				},
			},
		},
	}
	f.presetsParser.(*presets.FakeParser).MockConfig = map[string]*presets.PresetConfig{
		"laravel": config,
	}
	f.presetsParser.(*presets.FakeParser).MockPresetKeyContent = map[string]map[string]string{
		"laravel": map[string]string{
			"docker-compose.yml": defaultCompose,
		},
	}
	f.promptSelect.(*shell.FakePromptSelect).MockAnswer = map[string]string{
		"What database service do you want to use": "mysql",
	}
	f.presetsParser.(*presets.FakeParser).MockTemplates = map[string]map[string]string{
		"database": map[string]string{
			"mysql.yml": mysqlTemplate,
		},
	}

	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if val, ok := f.presetsParser.(*presets.FakeParser).CalledGetPresetKeyContent["laravel"]["docker-compose.yml"]; !ok || !val {
		t.Error("failed calling presetsParser.GetPresetKeyContent for preset 'laravel' and key 'docker-compose.yml'")
	}

	if val, ok := f.composeParser.(*compose.FakeParser).CalledLoad[defaultCompose]; !ok || !val {
		t.Error("failed calling compose.Load")
	}

	if val, ok := f.composeParser.(*compose.FakeParser).CalledSetService["database"][mysqlTemplate]; !ok || !val {
		t.Error("failed calling compose.SetService to database mysql service")
	}

	if !f.composeParser.(*compose.FakeParser).CalledString {
		t.Error("failed calling compose.String to database mysql service")
	}
}

func TestCustomDockerNoneOptionComposePresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.presetsParser.(*presets.FakeParser).MockExists = true

	config := &presets.PresetConfig{
		Questions: map[string]presets.PresetConfigQuestion{
			"database": presets.PresetConfigQuestion{
				Message: "What database service do you want to use",
				Options: yaml.MapSlice{
					yaml.MapItem{Key: "mysql", Value: "mysql"},
					yaml.MapItem{Key: "postgresql", Value: "postgresql"},
					yaml.MapItem{Key: "none", Value: "none"},
				},
			},
		},
	}
	f.presetsParser.(*presets.FakeParser).MockConfig = map[string]*presets.PresetConfig{
		"laravel": config,
	}
	f.presetsParser.(*presets.FakeParser).MockPresetKeyContent = map[string]map[string]string{
		"laravel": map[string]string{
			"docker-compose.yml": defaultCompose,
		},
	}
	f.promptSelect.(*shell.FakePromptSelect).MockAnswer = map[string]string{
		"What database service do you want to use": "none",
	}
	f.presetsParser.(*presets.FakeParser).MockTemplates = map[string]map[string]string{
		"database": map[string]string{
			"mysql.yml": mysqlTemplate,
		},
	}

	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if val, ok := f.presetsParser.(*presets.FakeParser).CalledGetPresetKeyContent["laravel"]["docker-compose.yml"]; !ok || !val {
		t.Error("failed calling presetsParser.GetPresetKeyContent for preset 'laravel' and key 'docker-compose.yml'")
	}

	if val, ok := f.composeParser.(*compose.FakeParser).CalledLoad[defaultCompose]; !ok || !val {
		t.Error("failed calling compose.Load")
	}

	if val, ok := f.composeParser.(*compose.FakeParser).CalledRemoveService["database"]; !ok || !val {
		t.Error("failed calling compose.RemoveService to database service")
	}

	if val, ok := f.composeParser.(*compose.FakeParser).CalledRemoveVolume["database"]; !ok || !val {
		t.Error("failed calling compose.RemoveService to database service")
	}

	if _, ok := f.composeParser.(*compose.FakeParser).CalledSetService["database"][mysqlTemplate]; ok {
		t.Error("should not call compose.SetService to database service")
	}

	if !f.composeParser.(*compose.FakeParser).CalledString {
		t.Error("failed calling compose.String to database mysql service")
	}
}

func TestErrorAskForServicePresetCommand(t *testing.T) {
	f := newFakeKoolPreset()

	f.presetsParser.(*presets.FakeParser).MockExists = true

	config := &presets.PresetConfig{
		Questions: map[string]presets.PresetConfigQuestion{
			"database": presets.PresetConfigQuestion{
				Message: "What database service do you want to use",
				Options: yaml.MapSlice{
					yaml.MapItem{Key: "mysql", Value: "mysql"},
					yaml.MapItem{Key: "postgresql", Value: "postgresql"},
				},
			},
		},
	}
	f.presetsParser.(*presets.FakeParser).MockConfig = map[string]*presets.PresetConfig{
		"laravel": config,
	}

	f.promptSelect.(*shell.FakePromptSelect).MockError = map[string]error{
		"What database service do you want to use": errors.New("database question error"),
	}

	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error")
	}

	err := f.shell.(*shell.FakeShell).Err

	if err == nil {
		t.Error("expecting an error, got none")
	} else if err.Error() != "database question error" {
		t.Errorf("expecting error 'database question error', got %v", err)
	}
}

func TestErrorLoadComposePresetCommand(t *testing.T) {
	f := newFakeKoolPreset()

	f.presetsParser.(*presets.FakeParser).MockExists = true

	config := &presets.PresetConfig{
		Questions: map[string]presets.PresetConfigQuestion{
			"database": presets.PresetConfigQuestion{
				Message: "What database service do you want to use",
				Options: yaml.MapSlice{
					yaml.MapItem{Key: "mysql", Value: "mysql"},
					yaml.MapItem{Key: "postgresql", Value: "postgresql"},
				},
			},
		},
	}
	f.presetsParser.(*presets.FakeParser).MockConfig = map[string]*presets.PresetConfig{
		"laravel": config,
	}
	f.presetsParser.(*presets.FakeParser).MockPresetKeyContent = map[string]map[string]string{
		"laravel": map[string]string{
			"docker-compose.yml": defaultCompose,
		},
	}
	f.promptSelect.(*shell.FakePromptSelect).MockAnswer = map[string]string{
		"What database service do you want to use": "mysql",
	}
	f.composeParser.(*compose.FakeParser).MockLoadError = errors.New("compose load error")

	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error")
	}

	err := f.shell.(*shell.FakeShell).Err

	if err == nil {
		t.Error("expecting an error, got none")
	} else if err.Error() != "Failed to write preset file docker-compose.yml: compose load error" {
		t.Errorf("expecting error 'Failed to write preset file docker-compose.yml: compose load error', got %v", err)
	}
}

func TestErrorSetComposeServicePresetCommand(t *testing.T) {
	f := newFakeKoolPreset()

	f.presetsParser.(*presets.FakeParser).MockExists = true

	config := &presets.PresetConfig{
		Questions: map[string]presets.PresetConfigQuestion{
			"database": presets.PresetConfigQuestion{
				Message: "What database service do you want to use",
				Options: yaml.MapSlice{
					yaml.MapItem{Key: "mysql", Value: "mysql"},
					yaml.MapItem{Key: "postgresql", Value: "postgresql"},
				},
			},
		},
	}
	f.presetsParser.(*presets.FakeParser).MockConfig = map[string]*presets.PresetConfig{
		"laravel": config,
	}
	f.presetsParser.(*presets.FakeParser).MockPresetKeyContent = map[string]map[string]string{
		"laravel": map[string]string{
			"docker-compose.yml": defaultCompose,
		},
	}
	f.promptSelect.(*shell.FakePromptSelect).MockAnswer = map[string]string{
		"What database service do you want to use": "mysql",
	}
	f.presetsParser.(*presets.FakeParser).MockTemplates = map[string]map[string]string{
		"database": map[string]string{
			"mysql.yml": mysqlTemplate,
		},
	}

	f.composeParser.(*compose.FakeParser).MockSetServiceError = errors.New("compose set service error")

	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error")
	}

	err := f.shell.(*shell.FakeShell).Err

	if err == nil {
		t.Error("expecting an error, got none")
	} else if err.Error() != "Failed to write preset file docker-compose.yml: compose set service error" {
		t.Errorf("expecting error 'Failed to write preset file docker-compose.yml: compose set service error', got %v", err)
	}
}

func TestErrorComposeStringPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()

	f.presetsParser.(*presets.FakeParser).MockExists = true

	config := &presets.PresetConfig{
		Questions: map[string]presets.PresetConfigQuestion{
			"database": presets.PresetConfigQuestion{
				Message: "What database service do you want to use",
				Options: yaml.MapSlice{
					yaml.MapItem{Key: "mysql", Value: "mysql"},
					yaml.MapItem{Key: "postgresql", Value: "postgresql"},
				},
			},
		},
	}
	f.presetsParser.(*presets.FakeParser).MockConfig = map[string]*presets.PresetConfig{
		"laravel": config,
	}
	f.presetsParser.(*presets.FakeParser).MockPresetKeyContent = map[string]map[string]string{
		"laravel": map[string]string{
			"docker-compose.yml": defaultCompose,
		},
	}
	f.promptSelect.(*shell.FakePromptSelect).MockAnswer = map[string]string{
		"What database service do you want to use": "mysql",
	}
	f.presetsParser.(*presets.FakeParser).MockTemplates = map[string]map[string]string{
		"database": map[string]string{
			"mysql.yml": mysqlTemplate,
		},
	}

	f.composeParser.(*compose.FakeParser).MockStringError = errors.New("compose string error")

	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error")
	}

	err := f.shell.(*shell.FakeShell).Err

	if err == nil {
		t.Error("expecting an error, got none")
	} else if err.Error() != "Failed to write preset file docker-compose.yml: compose string error" {
		t.Errorf("expecting error 'Failed to write preset file docker-compose.yml: compose string error', got %v", err)
	}
}

func TestErrorGetConfigPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()

	f.presetsParser.(*presets.FakeParser).MockExists = true

	f.presetsParser.(*presets.FakeParser).MockGetConfigError = map[string]error{
		"laravel": errors.New("get config error"),
	}

	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error")
	}

	err := f.shell.(*shell.FakeShell).Err

	if err == nil {
		t.Error("expecting an error, got none")
	} else if err.Error() != "error parsing preset config; err: get config error" {
		t.Errorf("expecting error 'error parsing preset config; err: get config error', got %v", err)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Error")
	}
}
