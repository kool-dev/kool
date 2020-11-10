package cmd

import (
	"strings"
	"errors"
	"fmt"
	"kool-dev/kool/cmd/presets"
	"kool-dev/kool/cmd/shell"

	"gopkg.in/yaml.v2"
	"github.com/spf13/cobra"
)

// KoolPresetFlags holds the flags for the preset command
type KoolPresetFlags struct {
	Override bool
}

// KoolPreset holds handlers and functions to implement the preset command logic
type KoolPreset struct {
	DefaultKoolService
	Flags        *KoolPresetFlags
	parser       presets.Parser
	terminal     shell.TerminalChecker
	promptSelect shell.PromptSelect
}

// ErrPresetFilesAlreadyExists error for existing presets files
var ErrPresetFilesAlreadyExists = errors.New("some preset files already exist")

func init() {
	var (
		preset    = NewKoolPreset()
		presetCmd = NewPresetCommand(preset)
	)

	rootCmd.AddCommand(presetCmd)
}

// NewKoolPreset creates a new handler for preset logic
func NewKoolPreset() *KoolPreset {
	return &KoolPreset{
		*newDefaultKoolService(),
		&KoolPresetFlags{false},
		&presets.DefaultParser{Presets: presets.GetAll()},
		shell.NewTerminalChecker(),
		shell.NewPromptSelect(),
	}
}

// Execute runs the preset logic with incoming arguments.
func (p *KoolPreset) Execute(args []string) (err error) {
	var (
		fileError, preset, language, database, cache string
		defaultCompose bool
	)

	if len(args) == 0 {
		if !p.IsTerminal() {
			err = fmt.Errorf("the input device is not a TTY; for non-tty environments, please specify a preset argument")
			return
		}

		if language, err = p.promptSelect.Ask("What language do you want to use", p.parser.GetLanguages()); err != nil {
			return
		}

		if preset, err = p.promptSelect.Ask("What preset do you want to use", p.parser.GetPresets(language)); err != nil {
			return
		}

		if askDatabase := p.parser.GetPresetMetaValue(preset, "ask_database"); askDatabase != "" {
			dbOptions := strings.Split(askDatabase, ",")

			if database, err = p.promptSelect.Ask("What database service do you want to use", dbOptions); err != nil {
				return
			}
		}

		if askCache := p.parser.GetPresetMetaValue(preset, "ask_cache"); askCache != "" {
			cacheOptions := strings.Split(askCache, ",")

			if cache, err = p.promptSelect.Ask("What cache service do you want to use", cacheOptions); err != nil {
				return
			}
		}

		defaultCompose = false
	} else {
		preset = args[0]
		defaultCompose = true
	}

	if !p.parser.Exists(preset) {
		err = fmt.Errorf("Unknown preset %s", preset)
		return
	}

	p.Println("Preset", preset, "is initializing!")

	if !p.Flags.Override {
		existingFiles := p.parser.LookUpFiles(preset)
		for _, fileName := range existingFiles {
			p.Warning("Preset file ", fileName, " already exists.")
		}

		if len(existingFiles) > 0 {
			err = ErrPresetFilesAlreadyExists
			return
		}
	}

	files := p.parser.GetPresetContents(preset)

	templates := presets.GetTemplates()

	for fileName, fileContent := range files {
		if fileName == "docker-compose.yml" && !defaultCompose {
			var compose, composeServices, composeVolumes yaml.MapSlice

			compose = append(compose, yaml.MapItem{Key: "version", Value: "3.7"})

			appKey := p.parser.GetPresetMetaValue(preset, "app_template")

			if err = appendYml(&composeServices, "app", templates["app"][appKey]); err != nil {
				err = fmt.Errorf("Failed to write preset file %s: %v", fileName, err)
				return
			}

			if database != "" && database != "none" {
				databaseKey := formatTemplateKey(database)

				if err = appendYml(&composeServices, "database", templates["database"][databaseKey]); err != nil {
					err = fmt.Errorf("Failed to write preset file %s: %v", fileName, err)
					return
				}

				composeVolumes = append(composeVolumes, yaml.MapItem{Key: "db"})
			}

			if cache != "" && cache != "none" {
				cacheKey := formatTemplateKey(cache)

				if err = appendYml(&composeServices, "cache", templates["cache"][cacheKey]); err != nil {
					err = fmt.Errorf("Failed to write preset file %s: %v", fileName, err)
					return
				}

				composeVolumes = append(composeVolumes, yaml.MapItem{Key: "cache"})
			}

			if len(composeServices) > 0 {
				compose = append(compose, yaml.MapItem{Key: "services", Value: composeServices})
			}

			if len(composeVolumes) > 0 {
				compose = append(compose, yaml.MapItem{Key: "volumes", Value: composeVolumes})
			}

			if err = appendYml(&compose, "networks", templates["shared"]["networks.yml"]); err != nil {
				err = fmt.Errorf("Failed to write preset file %s: %v", fileName, err)
				return
			}

			var parsedBytes []byte

			if parsedBytes, err = yaml.Marshal(compose); err != nil {
				err = fmt.Errorf("Failed to write preset file %s: %v", fileName, err)
				return
			}

			fileContent = string(parsedBytes)
		}

		if fileError, err = p.parser.WriteFile(fileName, fileContent); err != nil {
			err = fmt.Errorf("Failed to write preset file %s: %v", fileError, err)
			return
		}
	}

	p.Success("Preset ", preset, " initialized!")
	return
}

// NewPresetCommand initializes new kool preset command
func NewPresetCommand(preset *KoolPreset) (presetCmd *cobra.Command) {
	presetCmd = &cobra.Command{
		Use:   "preset [PRESET]",
		Short: "Initialize kool preset in the current working directory. If no preset argument is specified you will be prompted to pick among the existing options.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			preset.SetWriter(cmd.OutOrStdout())

			if err := preset.Execute(args); err != nil {
				if err.Error() == ErrPresetFilesAlreadyExists.Error() {
					preset.Warning("Some preset files already exist. In case you wanna override them, use --override.")
					preset.Exit(2)
				} else if err.Error() == shell.ErrPromptSelectInterrupted.Error() {
					preset.Warning("Operation Cancelled")
					preset.Exit(0)
				} else {
					preset.Error(err)
					preset.Exit(1)
				}
			}
		},
	}

	presetCmd.Flags().BoolVarP(&preset.Flags.Override, "override", "", false, "Force replace local existing files with the preset files")
	return
}

func parseYml(data string) (yaml.MapSlice, error) {
	parsed := yaml.MapSlice{}

	if err := yaml.Unmarshal([]byte(data), &parsed); err != nil {
		return nil, err
	}

	return parsed, nil
}

func formatTemplateKey(key string) (formattedKey string) {
	formattedKey = strings.ReplaceAll(key, " ", "")
	formattedKey = strings.ReplaceAll(formattedKey, ".", "")
	formattedKey = strings.ToLower(formattedKey) + ".yml"
	return
}

func appendYml(services *yaml.MapSlice, name string, content string) (err error) {
	var template yaml.MapSlice

	if template, err = parseYml(content); err != nil {
		return
	}

	*services = append(*services, yaml.MapItem{Key: name, Value: template})
	return
}
