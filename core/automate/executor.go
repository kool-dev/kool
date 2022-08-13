package automate

import (
	"fmt"
	"io"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/yamler"
	"os"
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
}

func NewExecutor(sh shell.Shell, fn RetrieveSource) *Executor {
	return &Executor{
		sh:            sh,
		getFromSource: fn,
		local:         afero.NewOsFs(),
		prompter:      shell.NewPromptSelect(),
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

	if pick, err = e.prompter.Ask(action.Prompt, optionsList); err != nil {
		return
	}

	err = e.Do([]*ActionSet{optionsMap[pick]})
	return
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
