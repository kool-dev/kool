package parser

import (
	"fmt"
	"io"
	"kool-dev/kool/core/builder"
	"os"
	"strings"

	"github.com/agnivade/levenshtein"
	"gopkg.in/yaml.v3"
)

// SimilarThreshold represents the minimal Levenshteindistance of two
// script names for them to be considered similarss
const SimilarThreshold int = 2

type yamlMarshalFnType func(interface{}) ([]byte, error)

// KoolYaml holds the structure for parsing the custom commands file
type KoolYaml struct {
	Scripts       map[string]interface{}  `yaml:"scripts"`
	ScriptDetails map[string]ScriptDetail `yaml:"-"`
}

// ScriptDetail describes a kool.yml script with context
type ScriptDetail struct {
	Name     string   `json:"name"`
	Comments []string `json:"comments"`
	Commands []string `json:"commands"`
}

// KoolYamlParser holds logic for handling kool yaml
type KoolYamlParser interface {
	Parse(string) error
	HasScript(string) bool
	ParseCommands(string) ([]builder.Command, error)
	SetScript(string, []string)
	String() (string, error)
}

var yamlMarshalFn yamlMarshalFnType = yaml.Marshal

// ParseKoolYaml decodes the target kool.yml into its
// the expected KoolYaml representation.
func ParseKoolYaml(filePath string) (parsed *KoolYaml, err error) {
	var (
		file *os.File
		raw  []byte
	)

	if file, err = os.OpenFile(filePath, os.O_RDONLY, os.ModePerm); err != nil {
		return
	}

	defer file.Close()

	if raw, err = io.ReadAll(file); err != nil {
		return
	}

	parsed = new(KoolYaml)
	if err = yaml.Unmarshal(raw, parsed); err != nil {
		return
	}

	return
}

// ParseKoolYamlWithDetails decodes the target kool.yml and includes script details.
func ParseKoolYamlWithDetails(filePath string) (parsed *KoolYaml, err error) {
	var (
		file *os.File
		raw  []byte
		root yaml.Node
	)

	if file, err = os.OpenFile(filePath, os.O_RDONLY, os.ModePerm); err != nil {
		return
	}

	defer file.Close()

	if raw, err = io.ReadAll(file); err != nil {
		return
	}

	parsed = new(KoolYaml)
	if err = yaml.Unmarshal(raw, parsed); err != nil {
		return
	}

	if err = yaml.Unmarshal(raw, &root); err != nil {
		return
	}

	parsed.ScriptDetails = parseScriptDetails(&root)
	return
}

// Parse decodes the target kool.yml
func (y *KoolYaml) Parse(filePath string) (err error) {
	var parsed *KoolYaml
	if parsed, err = ParseKoolYamlWithDetails(filePath); err != nil {
		return
	}

	y.Scripts = parsed.Scripts
	y.ScriptDetails = parsed.ScriptDetails
	return
}

// HasScript tells if the given script exists on this parsed YAML.
func (y *KoolYaml) HasScript(script string) (has bool) {
	if y.Scripts != nil {
		_, has = y.Scripts[script]
	}
	return
}

// GetSimilars checks for scripts with similar name.
func (y *KoolYaml) GetSimilars(script string) (has bool, similars []string) {
	var name string
	for name = range y.Scripts {
		if levenshtein.ComputeDistance(name, script) < SimilarThreshold {
			has = true
			similars = append(similars, name)
		}
	}
	return
}

// ParseCommands parsed the given script from kool.yml file into a list
// of commands parsed.
func (y *KoolYaml) ParseCommands(script string) (commands []builder.Command, err error) {
	var (
		isSingle bool
		isList   bool
		line     string
		lines    []interface{}
		linesStr []string
		command  *builder.DefaultCommand
	)

	if line, isSingle = y.Scripts[script].(string); isSingle {
		if command, err = builder.ParseCommand(line); err != nil {
			return
		}

		commands = append(commands, command)
	} else if linesStr, isList = y.Scripts[script].([]string); isList {
		for _, line := range linesStr {
			if command, err = builder.ParseCommand(line); err != nil {
				return
			}

			commands = append(commands, command)
		}
	} else if lines, isList = y.Scripts[script].([]interface{}); isList {
		for _, i := range lines {
			var lineStr string
			if lineStr, isSingle = i.(string); !isSingle {
				err = fmt.Errorf("failed parsing script '%s': expected string or array of strings", script)
				return
			}
			if command, err = builder.ParseCommand(lineStr); err != nil {
				return
			}

			commands = append(commands, command)
		}
	} else {
		err = fmt.Errorf("failed parsing script '%s': expected string or array of strings", script)
	}
	return
}

// SetScript set script into kool yaml
func (y *KoolYaml) SetScript(key string, commands []string) {
	if len(commands) == 0 {
		return
	}

	if y.Scripts == nil {
		y.Scripts = make(map[string]interface{})
	}

	if y.ScriptDetails == nil {
		y.ScriptDetails = make(map[string]ScriptDetail)
	}

	currentDetail := y.ScriptDetails[key]
	currentDetail.Name = key
	currentDetail.Commands = append([]string{}, commands...)
	if currentDetail.Comments == nil {
		currentDetail.Comments = []string{}
	}
	y.ScriptDetails[key] = currentDetail

	if len(commands) == 1 {
		y.Scripts[key] = commands[0]
	} else {
		var scripts []interface{}

		for _, c := range commands {
			scripts = append(scripts, c)
		}
		y.Scripts[key] = scripts
	}
}

// String returns docker compose as string
func (y *KoolYaml) String() (content string, err error) {
	var parsedBytes []byte

	if parsedBytes, err = yamlMarshalFn(y); err != nil {
		return
	}

	content = string(parsedBytes)
	return
}

func parseScriptDetails(root *yaml.Node) map[string]ScriptDetail {
	result := make(map[string]ScriptDetail)
	if root == nil {
		return result
	}

	scriptsNode := findScriptsNode(root)
	if scriptsNode == nil || scriptsNode.Kind != yaml.MappingNode {
		return result
	}

	for i := 0; i+1 < len(scriptsNode.Content); i += 2 {
		keyNode := scriptsNode.Content[i]
		valueNode := scriptsNode.Content[i+1]
		name := keyNode.Value
		detail := ScriptDetail{
			Name:     name,
			Comments: collectComments(keyNode, valueNode),
			Commands: parseCommandsNode(valueNode),
		}
		if detail.Comments == nil {
			detail.Comments = []string{}
		}
		if detail.Commands == nil {
			detail.Commands = []string{}
		}
		result[name] = detail
	}

	return result
}

func findScriptsNode(root *yaml.Node) *yaml.Node {
	current := root
	if current.Kind == yaml.DocumentNode && len(current.Content) > 0 {
		current = current.Content[0]
	}

	if current.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i+1 < len(current.Content); i += 2 {
		keyNode := current.Content[i]
		valueNode := current.Content[i+1]
		if keyNode.Value == "scripts" {
			return valueNode
		}
	}

	return nil
}

func parseCommandsNode(node *yaml.Node) []string {
	if node == nil {
		return []string{}
	}

	switch node.Kind {
	case yaml.ScalarNode:
		return []string{node.Value}
	case yaml.SequenceNode:
		commands := make([]string, 0, len(node.Content))
		for _, item := range node.Content {
			if item.Kind == yaml.ScalarNode {
				commands = append(commands, item.Value)
			}
		}
		return commands
	default:
		return []string{}
	}
}

func collectComments(nodes ...*yaml.Node) []string {
	var comments []string
	for _, node := range nodes {
		if node == nil {
			continue
		}
		comments = appendCommentLines(comments, node.HeadComment)
		comments = appendCommentLines(comments, node.LineComment)
	}

	return comments
}

func appendCommentLines(comments []string, raw string) []string {
	if raw == "" {
		return comments
	}

	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "#")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !containsString(comments, line) {
			comments = append(comments, line)
		}
	}

	return comments
}

func containsString(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}

	return false
}
