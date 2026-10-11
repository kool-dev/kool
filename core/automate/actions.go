package automate

type ActionType uint

const (
	TypeUnknown ActionType = iota
	TypeCopy
	TypeScripts
	TypePrompt
	TypeRecipe
	TypeMerge
	TypeInput
	TypeAPI
	TypeReplace
)

// ActionSet represents a set of single actions or a question
type ActionSet struct {
	Name    string    `yaml:"name"`
	Actions []*Action `yaml:"actions"`
}

// Action is a union kind of type that holds
// one specific action within it; used for parsing
type Action struct {
	// ref
	Ref string `yaml:"ref"`

	// recipe
	Recipe string `yaml:"recipe"`
	// merge
	Merge string `yaml:"merge"`
	// copy
	Src string `yaml:"copy"`
	Dst string `yaml:"dst"`
	// scripts
	Scripts []string `yaml:"scripts"`
	// prompt
	Prompt  string       `yaml:"prompt"`
	Input   string       `yaml:"input"`
	Default string       `yaml:"default"`
	Options []*ActionSet `yaml:"options"`
	// api
	API     string      `yaml:"api"`
	Prompts []*APIField `yaml:"prompts"`
	// replace
	Replace string `yaml:"replace"`
}

// APIField describes a prompt backed by a field in an API response.
type APIField struct {
	Path     string `yaml:"path"`
	Options  string `yaml:"options"`
	Value    string `yaml:"value"`
	Label    string `yaml:"label"`
	Default  string `yaml:"default"`
	Prompt   string `yaml:"prompt"`
	Ref      string `yaml:"ref"`
	Multiple bool   `yaml:"multiple"`
}

// Type tells the actual implementation of this action
func (a *Action) Type() ActionType {
	if a.Scripts != nil {
		return TypeScripts
	}

	if a.Recipe != "" {
		return TypeRecipe
	}

	if a.Src != "" {
		return TypeCopy
	}

	if a.Prompt != "" {
		return TypePrompt
	}

	if a.Input != "" {
		return TypeInput
	}

	if a.Merge != "" {
		return TypeMerge
	}

	if a.API != "" {
		return TypeAPI
	}

	if a.Replace != "" {
		return TypeReplace
	}

	return TypeUnknown
}
