package automate

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"kool-dev/kool/core/shell"
)

const fakeMetadata = `{
	"bootVersion": {"default": "4.1.2.RELEASE", "values": [{"id": "4.1.2.RELEASE", "name": "4.1.2"}]},
	"language": {"default": "java", "values": [
		{"id": "java", "name": "Java"},
		{"id": "kotlin", "name": "Kotlin"}
	]},
	"dependencies": {"values": [
		{"name": "Web", "values": [
			{"id": "web", "name": "Spring Web"},
			{"id": "htmx", "name": "htmx", "versionRange": "[4.0.0.RELEASE,4.1.0.M1)"}
		]},
		{"name": "Ops", "values": [{"id": "actuator", "name": "Spring Boot Actuator"}]}
	]}
}`

type fakePromptInput struct {
	called bool
	answer string
}

func (f *fakePromptInput) Input(string, string) (string, error) {
	f.called = true
	return f.answer, nil
}

func newAnswersExecutor(t *testing.T, terminal bool) (*Executor, *shell.FakePromptSelect, *fakePromptInput) {
	t.Helper()

	for _, ref := range []string{"NAME", "LANGUAGE", "DEPENDENCIES"} {
		t.Setenv(ref, "")
	}

	prompter := &shell.FakePromptSelect{MockAnswer: map[string]string{}}
	input := &fakePromptInput{}
	e := NewExecutor(&shell.FakeShell{MockIsTerminal: terminal}, nil)
	e.prompter = prompter
	e.multiPrompter = prompter
	e.input = input

	return e, prompter, input
}

func metadataServer(t *testing.T) string {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(fakeMetadata))
	}))
	t.Cleanup(server.Close)

	return server.URL
}

func apiAction(url string) *Action {
	return &Action{API: url, Prompts: []*APIField{
		{Path: "language", Prompt: "Which language?", Ref: "LANGUAGE"},
		{Path: "dependencies", Prompt: "Which dependencies?", Ref: "DEPENDENCIES", Multiple: true, Compatible: "bootVersion"},
	}}
}

func TestParseAnswers(t *testing.T) {
	answers, err := ParseAnswers([]string{"A=1", "B=x=y", "C="})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if answers["A"] != "1" || answers["B"] != "x=y" || answers["C"] != "" {
		t.Errorf("unexpected answers: %v", answers)
	}

	if _, err = ParseAnswers([]string{"invalid"}); err == nil {
		t.Error("expected error for answer without '='")
	}
	if _, err = ParseAnswers([]string{"=value"}); err == nil {
		t.Error("expected error for answer without ref")
	}
}

func TestProvidedInputValue(t *testing.T) {
	e, _, input := newAnswersExecutor(t, false)
	e.SetAnswers(Answers{"NAME": "bff"}, true)

	if err := e.inputValue(&Action{Input: "Name", Ref: "NAME", Default: "app"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.called {
		t.Error("should not prompt for a provided answer")
	}
	if os.Getenv("NAME") != "bff" || e.promptState["NAME"] != "bff" {
		t.Errorf("unexpected value: %q", os.Getenv("NAME"))
	}
}

func TestProvidedPromptByLabelAndText(t *testing.T) {
	options := []*ActionSet{{Name: "None - do not use a database"}, {Name: "PostgreSQL 16"}}

	for answers, want := range map[string]string{
		"PostgreSQL 16": "PostgreSQL 16",
		"postgresql":    "PostgreSQL 16",
		"none":          "None - do not use a database",
	} {
		e, prompter, _ := newAnswersExecutor(t, false)
		e.SetAnswers(Answers{"DB": answers}, true)

		if err := e.prompt(&Action{Prompt: "Which database?", Ref: "DB", Options: options}); err != nil {
			t.Fatalf("unexpected error for %q: %v", answers, err)
		}
		if prompter.CalledAsk {
			t.Error("should not prompt for a provided answer")
		}
		if e.promptState["DB"] != want {
			t.Errorf("answer %q: expected %q, got %q", answers, want, e.promptState["DB"])
		}
	}

	// prompts without a ref can be answered by their text
	e, _, _ := newAnswersExecutor(t, false)
	e.SetAnswers(Answers{"Which database?": "PostgreSQL 16"}, true)
	if err := e.prompt(&Action{Prompt: "Which database?", Options: options}); err != nil {
		t.Errorf("unexpected error answering by prompt text: %v", err)
	}
}

func TestProvidedAPISelectByValueAndLabel(t *testing.T) {
	url := metadataServer(t)

	for _, language := range []string{"kotlin", "Kotlin", "KOTLIN"} {
		e, prompter, _ := newAnswersExecutor(t, false)
		e.SetAnswers(Answers{"LANGUAGE": language, "DEPENDENCIES": "web"}, true)

		if err := e.api(apiAction(url)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if prompter.CalledAsk {
			t.Error("should not prompt for provided answers")
		}
		if got := os.Getenv("LANGUAGE"); got != "kotlin" {
			t.Errorf("answer %q: expected kotlin, got %q", language, got)
		}
	}
}

func TestProvidedAPIMultiSelect(t *testing.T) {
	e, _, _ := newAnswersExecutor(t, false)
	e.SetAnswers(Answers{"LANGUAGE": "java", "DEPENDENCIES": "Web: Spring Web, spring boot actuator,web"}, true)

	if err := e.api(apiAction(metadataServer(t))); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("DEPENDENCIES"); got != "web,actuator" {
		t.Errorf("expected web,actuator, got %q", got)
	}
}

func TestProvidedAPIInvalidOption(t *testing.T) {
	url := metadataServer(t)

	e, _, _ := newAnswersExecutor(t, false)
	e.SetAnswers(Answers{"LANGUAGE": "cobol"}, true)
	err := e.api(apiAction(url))
	if err == nil || !strings.Contains(err.Error(), `invalid option "cobol"`) || !strings.Contains(err.Error(), `"kotlin" (Kotlin)`) {
		t.Errorf("expected invalid option error listing valid options, got: %v", err)
	}

	e, _, _ = newAnswersExecutor(t, false)
	e.SetAnswers(Answers{"LANGUAGE": "java", "DEPENDENCIES": "web,htmx"}, true)
	err = e.api(apiAction(url))
	if err == nil || !strings.Contains(err.Error(), `"htmx" requires [4.0.0.RELEASE,4.1.0.M1)`) {
		t.Errorf("expected incompatible option error, got: %v", err)
	}
}

func TestMissingAnswerNonInteractive(t *testing.T) {
	for name, setup := range map[string]func(e *Executor){
		"no terminal":    func(e *Executor) { e.SetAnswers(nil, false) },
		"no interaction": func(e *Executor) { e.sh = &shell.FakeShell{MockIsTerminal: true}; e.SetAnswers(Answers{}, true) },
	} {
		e, prompter, input := newAnswersExecutor(t, false)
		setup(e)

		err := e.inputValue(&Action{Input: "Name", Ref: "NAME"})
		if err == nil || !strings.Contains(err.Error(), `no answer provided for "NAME"`) {
			t.Errorf("%s: expected missing answer error, got: %v", name, err)
		}

		err = e.prompt(&Action{Prompt: "Which database?", Options: []*ActionSet{{Name: "None"}}})
		if err == nil || !strings.Contains(err.Error(), `no answer provided for "Which database?"`) {
			t.Errorf("%s: expected missing answer error, got: %v", name, err)
		}

		err = e.api(apiAction(metadataServer(t)))
		if err == nil || !strings.Contains(err.Error(), `no answer provided for "LANGUAGE"`) {
			t.Errorf("%s: expected missing answer error, got: %v", name, err)
		}

		if prompter.CalledAsk || input.called {
			t.Errorf("%s: should not prompt in non-interactive mode", name)
		}
	}
}

func TestInteractivePromptsWithoutAnswers(t *testing.T) {
	e, prompter, input := newAnswersExecutor(t, true)
	input.answer = "typed"
	prompter.MockAnswer["Which language?"] = "Kotlin"
	prompter.MockAnswer["Which dependencies?"] = "Web: Spring Web"
	prompter.MockAnswer["Which database?"] = "PostgreSQL 16"

	if err := e.inputValue(&Action{Input: "Name", Ref: "NAME"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !input.called || os.Getenv("NAME") != "typed" {
		t.Error("expected the input to be prompted")
	}

	if err := e.prompt(&Action{Prompt: "Which database?", Ref: "DB", Options: []*ActionSet{{Name: "None"}, {Name: "PostgreSQL 16"}}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.promptState["DB"] != "PostgreSQL 16" {
		t.Errorf("unexpected prompt selection: %q", e.promptState["DB"])
	}

	if err := e.api(apiAction(metadataServer(t))); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !prompter.CalledAsk || os.Getenv("LANGUAGE") != "kotlin" || os.Getenv("DEPENDENCIES") != "web" {
		t.Errorf("expected API prompts; got LANGUAGE=%q DEPENDENCIES=%q", os.Getenv("LANGUAGE"), os.Getenv("DEPENDENCIES"))
	}
}

func TestVersionInRange(t *testing.T) {
	cases := []struct {
		version, versionRange string
		want                  bool
	}{
		{"4.1.2.RELEASE", "", true},
		{"4.1.2.RELEASE", "4.0.0.RELEASE", true},
		{"3.5.0", "4.0.0.RELEASE", false},
		{"4.1.2.RELEASE", "[4.0.0.RELEASE,4.1.0.M1)", false},
		{"4.0.9.RELEASE", "[4.0.0.RELEASE,4.1.0.M1)", true},
		{"4.1.0-M1", "[4.0.0.RELEASE,4.1.0.M1)", false},
		{"4.1.0.RC1", "[4.0.0.RELEASE,4.1.0.RC1]", true},
		{"4.2.0.BUILD-SNAPSHOT", "[4.1.0.RELEASE,4.2.0.M1)", false},
		{"4.2.0.BUILD-SNAPSHOT", "[4.1.0.RELEASE,4.2.0.RELEASE)", true},
		{"4.0.0.RELEASE", "(4.0.0.RELEASE,4.3.0.M1)", false},
	}

	for _, c := range cases {
		got, err := versionInRange(c.version, c.versionRange)
		if err != nil {
			t.Errorf("%s in %s: unexpected error: %v", c.version, c.versionRange, err)
		}
		if got != c.want {
			t.Errorf("%s in %s: expected %v, got %v", c.version, c.versionRange, c.want, got)
		}
	}

	if _, err := versionInRange("invalid", "4.0.0"); err == nil {
		t.Error("expected error for invalid version")
	}
}

func TestPromptExportsAndValidatesCachedVersion(t *testing.T) {
	t.Setenv("JAVA_VERSION", "")
	e, _, _ := newAnswersExecutor(t, false)
	e.SetAnswers(Answers{"JAVA_VERSION": "21"}, true)
	action := &Action{Prompt: "Java version", Ref: "JAVA_VERSION", Options: []*ActionSet{{Name: "17"}, {Name: "21"}, {Name: "25"}}}
	if err := e.prompt(action); err != nil {
		t.Fatal(err)
	}
	if os.Getenv("JAVA_VERSION") != "21" {
		t.Fatal("selection not exported")
	}
	e.SetAnswers(nil, true)
	if err := e.prompt(action); err != nil {
		t.Fatalf("cached answer not reused: %v", err)
	}
	e.promptState["JAVA_VERSION"] = "7"
	if err := e.prompt(action); err == nil {
		t.Fatal("unsupported cached JDK accepted")
	}
}
