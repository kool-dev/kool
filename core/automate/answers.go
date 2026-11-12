package automate

import (
	"fmt"
	"strings"
)

// Answers holds values provided beforehand for the questions asked by
// presets and recipes, keyed by the action ref (or, for questions
// without a ref, by the question text itself).
type Answers map[string]string

// ParseAnswers parses a list of KEY=value strings into Answers.
func ParseAnswers(list []string) (answers Answers, err error) {
	answers = make(Answers, len(list))

	for _, item := range list {
		key, value, found := strings.Cut(item, "=")
		key = strings.TrimSpace(key)

		if !found || key == "" {
			err = fmt.Errorf("invalid answer %q: expected the format REF=value", item)
			return
		}

		answers[key] = value
	}

	return
}

// lookup finds the answer for a question by its ref or text
func (a Answers) lookup(ref, question string) (value string, found bool) {
	if ref != "" {
		if value, found = a[ref]; found {
			return
		}
	}

	if question != "" {
		value, found = a[question]
	}

	return
}

// questionKey returns how a question can be referenced when answering it
func questionKey(ref, question string) string {
	if ref != "" {
		return ref
	}

	return question
}

// matchOption resolves an answer to one of the available options. The answer
// may be the option value or its label (case insensitive, with or without the
// "Group: " prefix for grouped options), or an unambiguous label prefix.
func matchOption(answer string, options []apiOption) (option apiOption, err error) {
	answer = strings.TrimSpace(answer)

	for _, option = range options {
		if option.Value == answer {
			return
		}
	}

	var matches []apiOption

	for _, candidate := range options {
		_, name, grouped := strings.Cut(candidate.Label, ": ")

		if strings.EqualFold(candidate.Value, answer) ||
			strings.EqualFold(candidate.Label, answer) ||
			(grouped && strings.EqualFold(name, answer)) {
			matches = append(matches, candidate)
		}
	}

	if len(matches) == 0 && answer != "" {
		for _, candidate := range options {
			if strings.HasPrefix(strings.ToLower(candidate.Label), strings.ToLower(answer)) {
				matches = append(matches, candidate)
			}
		}
	}

	switch len(matches) {
	case 1:
		option = matches[0]
	case 0:
		err = fmt.Errorf("invalid option %q; valid options are: %s", answer, describeOptions(options))
	default:
		err = fmt.Errorf("ambiguous option %q matches: %s", answer, describeOptions(matches))
	}

	return
}

// matchOptions resolves a comma separated list of answers into options
func matchOptions(answer string, options []apiOption) (selected []apiOption, err error) {
	seen := make(map[string]bool)

	for _, item := range strings.Split(answer, ",") {
		if strings.TrimSpace(item) == "" {
			continue
		}

		var option apiOption
		if option, err = matchOption(item, options); err != nil {
			return
		}

		if !seen[option.Value] {
			seen[option.Value] = true
			selected = append(selected, option)
		}
	}

	return
}

func describeOptions(options []apiOption) string {
	described := make([]string, 0, len(options))

	for _, option := range options {
		if option.Value == option.Label {
			described = append(described, fmt.Sprintf("%q", option.Value))
		} else {
			described = append(described, fmt.Sprintf("%q (%s)", option.Value, option.Label))
		}
	}

	return strings.Join(described, ", ")
}
