package automate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/yamler"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	yaml3 "gopkg.in/yaml.v3"
)

type RetrieveSource func(string) ([]byte, error)

type Executor struct {
	sh            shell.Shell
	getFromSource RetrieveSource
	local         afero.Fs
	prompter      shell.PromptSelect
	multiPrompter shell.PromptMultiSelect
	input         shell.PromptInput

	// promptState is a map of prompt answers
	promptState map[string]string
	// copiedFiles contains files available to a following replace action
	copiedFiles []string
}

func NewExecutor(sh shell.Shell, fn RetrieveSource) *Executor {
	return &Executor{
		sh:            sh,
		getFromSource: fn,
		local:         afero.NewOsFs(),
		prompter:      shell.NewPromptSelect(),
		multiPrompter: shell.NewPromptMultiSelect(),
		input:         shell.NewPromptInput(),
		promptState:   make(map[string]string),
		copiedFiles:   make([]string, 0),
	}
}

func (e *Executor) Do(steps []*ActionSet) (err error) {
	var (
		step   *ActionSet
		action *Action
	)

	for _, step = range steps {
		if step.Name != "" {
			e.sh.Info("⇒ ", step.Name)
		}

		for _, action = range step.Actions {
			switch action.Type() {
			case TypeRecipe:
				// the 'recipe' operation will run a new recipe
				// that is composed by a new array of ActionSet
				if err = e.recipe(action); err != nil {
					return
				}
			case TypeCopy:
				if err = e.copy(action); err != nil {
					return
				}
			case TypeScripts:
				if err = e.scripts(action); err != nil {
					return
				}
			case TypeMerge:
				if err = e.merge(action); err != nil {
					return
				}
			case TypePrompt:
				if err = e.prompt(action); err != nil {
					return
				}
			case TypeInput:
				if err = e.inputValue(action); err != nil {
					return
				}
			case TypeAPI:
				if err = e.api(action); err != nil {
					return
				}
			case TypeReplace:
				if err = e.replace(action); err != nil {
					return
				}
			default:
				err = fmt.Errorf("ops, something is wrong with this preset config (%d)", action.Type())
				return
			}
		}
	}

	return
}

func (e *Executor) copy(action *Action) (err error) {
	var (
		data []byte
		file afero.File
		size int
	)

	// defaults to the same path/file
	if action.Dst == "" {
		action.Dst = action.Src
		e.sh.Println("→ copying", action.Src)
	} else {
		e.sh.Println("→ copying", action.Src, "as", action.Dst)
	}

	if data, err = e.getFromSource(action.Src); err != nil {
		return
	}

	if _, statErr := e.local.Stat(action.Dst); !os.IsNotExist(statErr) {
		renamedFile := fmt.Sprintf("%s.bak.%s", action.Dst, time.Now().Format("20060102"))

		e.sh.Warning(fmt.Sprintf(
			"File %s already exists, overriding. (backup is %s)",
			action.Dst,
			renamedFile,
		))

		if err = e.local.Rename(action.Dst, renamedFile); err != nil {
			return
		}
	}

	if file, err = e.local.OpenFile(action.Dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm); err != nil {
		return
	}

	if size, err = file.Write(data); err != nil {
		return
	}

	if len(data) != size {
		err = fmt.Errorf("failed writing file")
		return
	}

	if err = file.Sync(); err != nil {
		return
	}

	_ = file.Close()
	e.copiedFiles = append(e.copiedFiles, action.Dst)
	return
}

func (e *Executor) replace(action *Action) (err error) {
	var command builder.Command
	if command, err = builder.ParseCommand("replace " + action.Replace); err != nil {
		return
	}
	if len(command.Args()) != 2 {
		return fmt.Errorf("replace action expects a search and replacement value")
	}

	for _, filename := range e.copiedFiles {
		var data []byte
		if data, err = afero.ReadFile(e.local, filename); err != nil {
			return
		}
		data = bytes.ReplaceAll(data, []byte(command.Args()[0]), []byte(command.Args()[1]))
		if err = afero.WriteFile(e.local, filename, data, os.ModePerm); err != nil {
			return
		}
	}
	return
}

func (e *Executor) merge(action *Action) (err error) {
	var (
		data    []byte
		file    afero.File
		merger  = &yamler.DefaultMerger{}
		into    = &yaml3.Node{}
		partial = &yaml3.Node{}
	)

	// defaults to the same path/file
	if action.Dst == "" {
		action.Dst = action.Merge
		e.sh.Println("→ merging", action.Merge)
	} else {
		e.sh.Println("→ merging", action.Merge, "into", action.Dst)
	}

	// partial
	if data, err = e.getFromSource(action.Merge); err != nil {
		return
	}

	if err = yaml3.Unmarshal(data, partial); err != nil {
		return err
	}

	// into
	if file, err = e.local.OpenFile(action.Dst, os.O_RDONLY, os.ModePerm); err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("merge destiny file '%s' does not exist", action.Dst)
		}
		return
	}

	if data, err = io.ReadAll(file); err != nil {
		return
	}

	if err = file.Close(); err != nil {
		return
	}

	if err = yaml3.Unmarshal(data, into); err != nil {
		return err
	}

	if err = merger.Merge(partial, into); err != nil {
		return
	}

	err = new(yamler.DefaultOutputWritter).WriteYAML(action.Dst, into)
	return
}

func (e *Executor) prompt(action *Action) (err error) {
	var (
		optionsList []string
		optionsMap  = make(map[string]*ActionSet)
		opt         *ActionSet
		pick        string
	)

	for _, opt = range action.Options {
		optionsList = append(optionsList, opt.Name)
		optionsMap[opt.Name] = opt
	}

	if e.promptState[action.Ref] != "" {
		// we already got a value for this prompt, so we skip it and reuse the value
		e.sh.Printf("→ Already selected '%s': %s\n", action.Prompt, e.promptState[action.Ref])

		pick = e.promptState[action.Ref]
	} else if pick, err = e.prompter.Ask(action.Prompt, optionsList); err != nil {
		return
	}

	// store selection for later use if needed
	if action.Ref != "" {
		e.promptState[action.Ref] = pick
	}

	err = e.Do([]*ActionSet{optionsMap[pick]})
	return
}

func (e *Executor) inputValue(action *Action) (err error) {
	value, err := e.input.Input(action.Input, action.Default)
	if err != nil {
		return
	}
	if action.Ref != "" {
		e.promptState[action.Ref] = value
	}

	if action.Ref != "" {
		err = os.Setenv(action.Ref, value)
	}
	return
}

func (e *Executor) api(action *Action) (err error) {
	var (
		response *http.Response
		data     map[string]any
	)

	client := &http.Client{Timeout: 15 * time.Second}
	if response, err = client.Get(action.API); err != nil {
		return
	}
	defer func() {
		if closeErr := response.Body.Close(); err == nil {
			err = closeErr
		}
	}()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("API request %s returned HTTP %d", action.API, response.StatusCode)
	}
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return fmt.Errorf("could not decode API response from %s: %w", action.API, err)
	}

	for _, field := range action.Prompts {
		var (
			values       []apiOption
			defaultValue string
		)
		if values, defaultValue, err = apiOptions(data, field); err != nil {
			return fmt.Errorf("could not read API field %q: %w", field.Path, err)
		}
		if len(values) == 0 {
			return fmt.Errorf("API field %q did not contain any options", field.Path)
		}

		options := make([]string, 0, len(values))
		byLabel := make(map[string]string, len(values))
		for _, option := range values {
			options = append(options, option.Label)
			byLabel[option.Label] = option.Value
		}
		if defaultValue != "" {
			defaultLabel := ""
			for _, option := range values {
				if option.Value == defaultValue {
					defaultLabel = option.Label
					break
				}
			}
			if defaultLabel != "" {
				for i, option := range options {
					if option == defaultLabel {
						options = append([]string{defaultLabel}, append(options[:i], options[i+1:]...)...)
						break
					}
				}
			}
		}

		if field.Multiple {
			var selected []string
			if selected, err = e.multiPrompter.AskMany(field.Prompt, options); err != nil {
				return
			}
			ids := make([]string, 0, len(selected))
			for _, label := range selected {
				if value, ok := byLabel[label]; ok {
					ids = append(ids, value)
				}
			}
			err = e.setAPIValue(field.Ref, strings.Join(ids, ","))
			continue
		}

		var selected string
		if selected, err = e.prompter.Ask(field.Prompt, options); err != nil {
			return
		}
		err = e.setAPIValue(field.Ref, byLabel[selected])
	}
	return
}

type apiOption struct {
	Label string
	Value string
}

func apiOptions(data map[string]any, field *APIField) (options []apiOption, defaultValue string, err error) {
	var current any
	if current, err = apiPath(data, field.Path); err != nil {
		return nil, "", err
	}
	if object, ok := current.(map[string]any); ok {
		defaultValue, _ = object["default"].(string)
		if field.Options != "" {
			current, err = apiPath(object, field.Options)
		} else {
			current = object["values"]
		}
	} else if field.Options != "" {
		return nil, "", fmt.Errorf("options path %q requires an object at path %q", field.Options, field.Path)
	}
	if field.Default != "" {
		defaultValue = field.Default
	}
	if err != nil {
		return nil, "", err
	}
	values, ok := current.([]any)
	if !ok {
		return nil, "", fmt.Errorf("path %q has no options", field.Path)
	}
	options, err = flattenAPIOptions(values, "", field.Value, field.Label)
	return
}

func apiPath(root any, path string) (any, error) {
	current := root
	for _, part := range strings.Split(path, ".") {
		if part == "" {
			continue
		}
		value, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("options path %q is not an object", path)
		}
		current, ok = value[part]
		if !ok {
			return nil, fmt.Errorf("options path %q does not exist", path)
		}
	}
	return current, nil
}

func flattenAPIOptions(values []any, group, valueKey, labelKey string) ([]apiOption, error) {
	var options []apiOption
	for _, raw := range values {
		object, ok := raw.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("option is not an object")
		}
		name, _ := object["name"].(string)
		id, _ := object["id"].(string)
		if nested, ok := object["values"].([]any); ok {
			if nestedOptions, err := flattenAPIOptions(nested, name, valueKey, labelKey); err != nil {
				return nil, err
			} else {
				options = append(options, nestedOptions...)
			}
			continue
		}
		if valueKey != "" {
			id, _ = object[valueKey].(string)
		}
		if labelKey != "" {
			name, _ = object[labelKey].(string)
		}
		if id == "" || name == "" {
			continue
		}
		label := name
		if group != "" {
			label = group + ": " + name
		}
		options = append(options, apiOption{Label: label, Value: id})
	}
	return options, nil
}

func (e *Executor) setAPIValue(ref, value string) error {
	if ref == "" {
		return nil
	}
	e.promptState[ref] = value
	return os.Setenv(ref, value)
}

func (e *Executor) recipe(action *Action) (err error) {
	var (
		set  = new(ActionSet)
		data []byte
	)

	if data, err = recipesSource.ReadFile(fmt.Sprintf("recipes/%s.yml", action.Recipe)); err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("recipe '%s' does not exist", action.Recipe)
		}
		return
	}

	if err = yaml.Unmarshal(data, set); err != nil {
		return
	}

	err = e.Do([]*ActionSet{set})
	return
}

func (e *Executor) scripts(action *Action) (err error) {
	var (
		command  builder.Command
		commands []builder.Command
		line     string
	)

	for _, line = range action.Scripts {
		if command, err = builder.ParseCommand(line); err != nil {
			return
		}

		commands = append(commands, command)
	}

	// all commands have parsed succussfully; now execute them
	for _, command = range commands {
		e.sh.Println("→ exec:", command.String())
		if err = e.sh.Interactive(command); err != nil {
			return
		}
	}

	return
}
