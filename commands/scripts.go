package commands

import (
	"encoding/json"
	"errors"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/parser"
	"path"

	"github.com/spf13/cobra"
)

// KoolScriptsFlags holds the flags for the scripts command
type KoolScriptsFlags struct {
	JSON bool
}

// KoolScripts holds handlers and functions to implement the scripts command logic
type KoolScripts struct {
	DefaultKoolService
	Flags  *KoolScriptsFlags
	parser parser.Parser
	env    environment.EnvStorage
}

func AddKoolScripts(root *cobra.Command) {
	var (
		scripts    = NewKoolScripts()
		scriptsCmd = NewScriptsCommand(scripts)
	)

	root.AddCommand(scriptsCmd)
}

// NewKoolScripts creates a new handler for scripts logic
func NewKoolScripts() *KoolScripts {
	return &KoolScripts{
		*newDefaultKoolService(),
		&KoolScriptsFlags{},
		parser.NewParser(),
		environment.NewEnvStorage(),
	}
}

// Execute runs the scripts logic with incoming arguments.
func (s *KoolScripts) Execute(args []string) (err error) {
	var filter string
	if len(args) > 0 {
		filter = args[0]
	}

	cwdErr := s.parser.AddLookupPath(s.env.Get("PWD"))
	homeErr := s.parser.AddLookupPath(path.Join(s.env.Get("HOME"), "kool"))

	if isKoolYmlNotFound(cwdErr) && isKoolYmlNotFound(homeErr) {
		if s.Flags.JSON {
			return s.printJSON([]parser.ScriptDetail{})
		}
		s.Shell().Warning("No kool.yml found in current directory or ~/kool.")
		return nil
	}

	if err = firstLookupError(cwdErr, homeErr); err != nil {
		return
	}

	if s.Flags.JSON {
		var details []parser.ScriptDetail
		if details, err = s.parser.ParseAvailableScriptsDetails(filter); err != nil {
			return
		}
		return s.printJSON(details)
	}

	var scripts []string
	if scripts, err = s.parser.ParseAvailableScripts(filter); err != nil {
		return
	}

	if len(scripts) == 0 {
		if filter == "" {
			s.Shell().Warning("No scripts found.")
		} else {
			s.Shell().Warning("No scripts found with prefix:", filter)
		}
		return nil
	}

	s.Shell().Info("Available scripts:")
	for _, script := range scripts {
		s.Shell().Println("  " + script)
	}

	return
}

// NewScriptsCommand initializes new kool scripts command
func NewScriptsCommand(scripts *KoolScripts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scripts [FILTER]",
		Short: "List scripts defined in kool.yml",
		Long: `List the scripts defined in kool.yml or kool.yaml in the current
working directory and in ~/kool. Use the optional FILTER to show only scripts
that start with a given prefix.`,
		Args:                  cobra.MaximumNArgs(1),
		RunE:                  DefaultCommandRunFunction(scripts),
		DisableFlagsInUseLine: true,
	}

	cmd.Flags().BoolVar(&scripts.Flags.JSON, "json", false, "Output scripts as JSON")

	return cmd
}

func isKoolYmlNotFound(err error) bool {
	return errors.Is(err, parser.ErrKoolYmlNotFound)
}

func firstLookupError(cwdErr, homeErr error) error {
	if cwdErr != nil && !isKoolYmlNotFound(cwdErr) {
		return cwdErr
	}

	if homeErr != nil && !isKoolYmlNotFound(homeErr) {
		return homeErr
	}

	return nil
}

func (s *KoolScripts) printJSON(details []parser.ScriptDetail) (err error) {
	if details == nil {
		details = []parser.ScriptDetail{}
	}

	for i := range details {
		if details[i].Comments == nil {
			details[i].Comments = []string{}
		}
		if details[i].Commands == nil {
			details[i].Commands = []string{}
		}
	}

	var payload []byte
	if payload, err = json.Marshal(details); err != nil {
		return
	}

	s.Shell().Println(string(payload))
	return nil
}
