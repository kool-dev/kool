package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/shlex"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// RunFlags holds the flags for the start command
type RunFlags struct {
}

// KoolYaml holds the structure for parsing the custom commands file
type KoolYaml struct {
	Scripts map[string]interface{} `yaml:"scripts"`
}

var runCmd = &cobra.Command{
	Use:   "run [script]",
	Short: "Runs a custom command defined at kool.yaml",
	Args:  cobra.MinimumNArgs(1),
	Run:   runRun,
}

var runFlags = &RunFlags{}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runRun(cmd *cobra.Command, args []string) {
	var (
		err    error
		script string = args[0]
	)

	commands := parseCustomCommandsScript(script)

	if len(args) > 1 && len(commands) > 1 {
		fmt.Println("Error: you cannot pass in extra arguments to multiple commands scripts.")
		os.Exit(2)
	}

	for _, exec := range commands {
		var execArgs = exec[1:]
		if len(commands) == 1 {
			// single command - forward extra args
			execArgs = append(execArgs, args[1:]...)
		}
		err = shellInteractive(exec[0], execArgs...)

		if err != nil {
			execError("", err)
			os.Exit(1)
		}
	}
}

func parseCustomCommandsScript(script string) (parsedCommands [][]string) {
	var (
		err      error
		fileName string
		file     *os.File
		yml      []byte
	)

	if _, err = os.Stat("kool.yml"); !os.IsNotExist(err) {
		fileName = "kool.yml"
	} else {
		if _, err = os.Stat("kool.yaml"); !os.IsNotExist(err) {
			fileName = "kool.yaml"
		}
	}

	if fileName == "" {
		fmt.Println("Could not find kool.yml in the current working directory.")
		os.Exit(2)
	}

	file, err = os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)

	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	defer file.Close()

	yml, err = ioutil.ReadAll(file)

	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	var parsed *KoolYaml = new(KoolYaml)
	err = yaml.Unmarshal(yml, parsed)
	yml = nil

	if err != nil {
		fmt.Println("Failed to parse", fileName, ":", err)
	}

	if _, found := parsed.Scripts[script]; !found {
		fmt.Println("Could not find script", script, "within", fileName)
		os.Exit(2)
	}

	if singleCommand, isSingleString := parsed.Scripts[script].(string); isSingleString {
		parsedCommands = append(parsedCommands, parseCustomCommand(singleCommand))
	} else if commands, isList := parsed.Scripts[script].([]interface{}); isList {
		for _, line := range commands {
			parsedCommands = append(parsedCommands, parseCustomCommand(line.(string)))
		}
	} else {
		fmt.Println("Could not parse script with key", script, ": it must be either a single command or an array of commands. Please refer to the documentation.")
		os.Exit(2)
	}

	return
}

func parseCustomCommand(line string) (parsed []string) {
	var err error

	parsed, err = shlex.Split(line)

	if err != nil {
		fmt.Println("Failed parsing custom command:", line, err)
		os.Exit(1)
	}

	for i := range parsed {
		for _, env := range os.Environ() {
			envPair := strings.SplitN(env, "=", 2)
			parsed[i] = strings.ReplaceAll(parsed[i], "$"+envPair[0], envPair[1])
		}
	}

	return
}
