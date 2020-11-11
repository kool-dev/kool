package cmd

import (
	"errors"
	"fmt"
	"kool-dev/kool/cmd/presets"
	"kool-dev/kool/cmd/shell"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
		useDefaultCompose                               bool
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
	} else {
		preset = args[0]
	}

	if !p.parser.Exists(preset) {
		err = fmt.Errorf("Unknown preset %s", preset)
		return
	}

	useDefaultCompose = true

	if dbOptionsStr := p.parser.GetPresetKeyContent(preset, "preset_database_options"); dbOptionsStr != "" && p.IsTerminal() {
		useDefaultCompose = false
		dbOptions := strings.Split(dbOptionsStr, ",")

		if database, err = p.promptSelect.Ask("What database service do you want to use", dbOptions); err != nil {
			return
		}
	}

	if cacheOptionsStr := p.parser.GetPresetKeyContent(preset, "preset_cache_options"); cacheOptionsStr != "" && p.IsTerminal() {
		useDefaultCompose = false
		cacheOptions := strings.Split(cacheOptionsStr, ",")

		if cache, err = p.promptSelect.Ask("What cache service do you want to use", cacheOptions); err != nil {
			return
		}
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

	presetKeys := p.parser.GetPresetKeys(preset)

	templates := presets.GetTemplates()

	for _, presetKey := range presetKeys {
		if strings.HasPrefix(presetKey, "preset_") {
			continue
		}

		var content string

		if presetKey == "docker-compose.yml" && !useDefaultCompose {
			var compose yaml.MapSlice

			defaultCompose := p.parser.GetPresetKeyContent(preset, presetKey)

			if compose, err = parseYml(defaultCompose); err != nil {
				err = fmt.Errorf("Failed to write preset file %s: %v", presetKey, err)
				return
			}

			if database != "" {
				databaseKey := formatTemplateKey(database)

				if compose, err = replaceComposeService(compose, "database", templates["database"][databaseKey]); err != nil {
					err = fmt.Errorf("Failed to write preset file %s: %v", presetKey, err)
					return
				}
			}

			if cache != "" {
				cacheKey := formatTemplateKey(cache)

				if compose, err = replaceComposeService(compose, "cache", templates["cache"][cacheKey]); err != nil {
					err = fmt.Errorf("Failed to write preset file %s: %v", presetKey, err)
					return
				}
			}

			var parsedBytes []byte

			if parsedBytes, err = yaml.Marshal(compose); err != nil {
				err = fmt.Errorf("Failed to write preset file %s: %v", presetKey, err)
				return
			}

			content = string(parsedBytes)
		} else {
			content = p.parser.GetPresetKeyContent(preset, presetKey)
		}

		if fileError, err = p.parser.WriteFile(presetKey, content); err != nil {
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

func replaceComposeService(compose yaml.MapSlice, name string, content string) (yaml.MapSlice, error) {
	var err error
	for sectionKey, section := range compose {
		if section.Key == "services" {
			for serviceKey, service := range section.Value.(yaml.MapSlice) {
				if service.Key == name {
					var template yaml.MapSlice

					if template, err = parseYml(content); err != nil {
						return compose, err
					}

					compose[sectionKey].Value.(yaml.MapSlice)[serviceKey].Value = template
					return compose, nil
				}
			}
		}
	}

	return compose, nil
}
