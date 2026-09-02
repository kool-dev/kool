package automate

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func parseAction(data string, t *testing.T) *Action {
	a := new(Action)
	if err := yaml.Unmarshal([]byte(data), a); err != nil {
		t.Fatalf("error trying to parse action: %v - [action YML: %v]", err, string(data))
	}
	return a
}

func TestParseActionAdd(t *testing.T) {
	t.Run("Basic add parse", func(t *testing.T) {
		a := parseAction("recipe: 'foo'", t)

		if a.Recipe != "foo" {
			t.Errorf("failed parsing ActionAdd - expected Recipe=foo: %v", a)
		}

		if a.Type() != TypeRecipe {
			t.Errorf("failed parsing ActionAdd type; got: %v - %+v", a.Type(), a)
		}
	})
}

func TestParseActionCopy(t *testing.T) {
	t.Run("Basic copy parse", func(t *testing.T) {
		a := parseAction("copy: 'foo'", t)

		if a.Src != "foo" || a.Dst != "" {
			t.Errorf("failed parsing ActionCopy - expected Src=foo Dst=<empty>: %v", a)
		}

		if a.Type() != TypeCopy {
			t.Errorf("failed parsing ActionCopy type; got: %v - %+v", a.Type(), a)
		}
	})
}

func TestParseActionScripts(t *testing.T) {
	t.Run("Parse scripts basic", func(t *testing.T) {
		a := parseAction("scripts: [ 'foo', 'bar' ]", t)

		if len(a.Scripts) != 2 || a.Scripts[0] != "foo" || a.Scripts[1] != "bar" {
			t.Errorf("failed parsing ActionScripts - expected foo, bar: %v", a)
		}

		if a.Type() != TypeScripts {
			t.Errorf("failed parsing ActionScripts type; got: %v - %+v", a.Type(), a)
		}
	})
}

func TestParseActionPrompt(t *testing.T) {
	t.Run("Parse prompt basic", func(t *testing.T) {
		a := parseAction("prompt: 'foo?'\ndefault: 'bar'", t)

		if a.Prompt != "foo?" || a.Default != "bar" || len(a.Options) > 0 {
			t.Errorf("failed parsing ActionPrompt - expected foo, bar: %v", a)
		}

		if a.Type() != TypePrompt {
			t.Errorf("failed parsing ActionPrompt type; got: %v - %+v", a.Type(), a)
		}
	})
}

func TestParseActionInput(t *testing.T) {
	a := parseAction("input: 'Project name'\nref: 'PROJECT_NAME'\ndefault: 'application'", t)

	if a.Input != "Project name" || a.Ref != "PROJECT_NAME" || a.Default != "application" {
		t.Errorf("failed parsing ActionInput: %v", a)
	}
	if a.Type() != TypeInput {
		t.Errorf("failed parsing ActionInput type; got: %v", a.Type())
	}
}

func TestParseActionAPI(t *testing.T) {
	a := parseAction("api: 'https://example.test/options'\nprompts:\n  - path: 'data'\n    options: 'items'\n    value: 'key'\n    label: 'title'\n    multiple: true", t)

	if a.API == "" || len(a.Prompts) != 1 || a.Prompts[0].Options != "items" || !a.Prompts[0].Multiple {
		t.Errorf("failed parsing ActionAPI: %v", a)
	}
	if a.Type() != TypeAPI {
		t.Errorf("failed parsing ActionAPI type; got: %v", a.Type())
	}
}

func TestParseActionReplace(t *testing.T) {
	a := parseAction("replace: 'JAVA_VERSION $JAVA_VERSION'", t)

	if a.Replace != "JAVA_VERSION $JAVA_VERSION" {
		t.Errorf("failed parsing ActionReplace: %v", a)
	}
	if a.Type() != TypeReplace {
		t.Errorf("failed parsing ActionReplace type; got: %v", a.Type())
	}
}

func TestAPIOptions(t *testing.T) {
	data := map[string]any{
		"data": map[string]any{
			"default": "two",
			"items": []any{
				map[string]any{"key": "one", "title": "One"},
				map[string]any{"key": "two", "title": "Two"},
			},
		},
	}
	options, defaultValue, err := apiOptions(data, &APIField{
		Path:    "data",
		Options: "items",
		Value:   "key",
		Label:   "title",
	})
	if err != nil || len(options) != 2 || options[1].Value != "two" || defaultValue != "two" {
		t.Errorf("unexpected API options: options=%v default=%q err=%v", options, defaultValue, err)
	}
}

func TestParseActionMerge(t *testing.T) {
	t.Run("Parse merge basic", func(t *testing.T) {
		a := parseAction("merge: 'foo'", t)

		if a.Merge != "foo" {
			t.Errorf("failed parsing ActionMerge - expected foo: %v", a)
		}

		if a.Type() != TypeMerge {
			t.Errorf("failed parsing ActionMerge type; got: %v - %+v", a.Type(), a)
		}
	})
}

func TestParseActionSets(t *testing.T) {
	t.Run("Parse no steps", func(t *testing.T) {
		s := new(ActionSet)
		if err := yaml.Unmarshal([]byte("name: foo\nactions:\n"), s); err != nil {
			t.Errorf("unexpected error parsing ActionSet: %v", err)
		}

		if s.Name != "foo" || len(s.Actions) != 0 {
			t.Errorf("failed parsing ActionSet; expected foo/0, got: %v", s)
		}
	})

	t.Run("Parse mutiple actions", func(t *testing.T) {
		s := new(ActionSet)
		if err := yaml.Unmarshal([]byte("actions:\n  - recipe: bar\n  - copy: file"), s); err != nil {
			t.Errorf("unexpected error parsing ActionSet: %v", err)
		}

		if len(s.Actions) != 2 || s.Actions[0].Type() != TypeRecipe || s.Actions[1].Type() != TypeCopy {
			t.Errorf("failed parsing ActionSet; expected 2/add/copy, got: %v", s)
		}
	})
}

func TestFullStepParse(t *testing.T) {
	config := `name: 'Select the desired setup'
actions:
  # Defines which app service to use (PHP version)
  - prompt: 'Which app service do you want to use'
    default: 'PHP 8.0'
    options:
      - name: 'PHP 8.0'
        actions:
          - recipe: php-8
      - name: 'PHP 7.4'
        actions:
          - recipe: php-7.4
`
	s := new(ActionSet)
	if err := yaml.Unmarshal([]byte(config), s); err != nil {
		t.Errorf("error parsing full ActionSet: %v", err)
	}

	if s.Name != "Select the desired setup" || len(s.Actions) != 1 || s.Actions[0].Type() != TypePrompt {
		t.Errorf("bad parse ActionSet; expected: select/1/prompt; got: %v", s)
	} else if s.Actions[0].Prompt != "Which app service do you want to use" || s.Actions[0].Default != "PHP 8.0" {
		t.Errorf("bad full ActionSet action[0] prompt/default: %v", s.Actions[0])
	} else if len(s.Actions[0].Options) != 2 {
		t.Errorf("bad full ActionSet action[0].Options: %v", s.Actions[0].Options)
	} else if s.Actions[0].Options[1].Name != "PHP 7.4" || len(s.Actions[0].Options[1].Actions) != 1 || s.Actions[0].Options[1].Actions[0].Type() != TypeRecipe {
		t.Errorf("bad full actionSet action[0].Options[1]; expected 7.4/1/add; got: %v", s.Actions[0].Options[1])
	}
}
