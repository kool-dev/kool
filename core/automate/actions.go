package automate

type ActionType uint

const (
	TypeUnknown ActionType = iota
	TypeCopy
	TypeScripts
	TypePrompt
	TypeAdd
)

// ActionStep represents a set of single actions or a question
type ActionStep struct {
	Name    string    `yaml:"name"`
	Actions []*Action `yaml:"actions"`
}

// Action is a union kind of type that holds
// one specific action within it; used for parsing
type Action struct {
	// add
	Recipe string `yaml:"add"`
	// copy
	Src string `yaml:"copy"`
	Dst string `yaml:"dst"`
	// scripts
	Scripts []string `yaml:"scripts"`
	// prompt
	Prompt  string        `yaml:"prompt"`
	Default string        `yaml:"default"`
	Options []*ActionStep `yaml:"options"`
}

// Type tells the actual implementation of this action
func (a *Action) Type() ActionType {
	if a.Scripts != nil {
		return TypeScripts
	}

	if a.Recipe != "" {
		return TypeAdd
	}

	if a.Src != "" {
		return TypeCopy
	}

	if a.Prompt != "" {
		return TypePrompt
	}

	return TypeUnknown
}
