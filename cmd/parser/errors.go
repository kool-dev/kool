package parser

import (
	"errors"
	"fmt"
	"strings"
)

// ErrMultipleDefinedScript happens when the script asked for is
// found within more than one kool.yml files targeted.
var ErrMultipleDefinedScript = fmt.Errorf("script was found in more than one kool.yml file")

// IsMultipleDefinedScriptError tells whether the given error is parser.ErrMultipleDefinedScript
func IsMultipleDefinedScriptError(err error) bool {
	return err.Error() == ErrMultipleDefinedScript.Error()
}

// ErrKoolYmlNotFound means there was no kool.yml file in the targeted folders
var ErrKoolYmlNotFound = errors.New("could not find any kool.yml file")

// ErrPossibleTypo implements error interface and can be used
// to determine specific situations of not-found scripts but
// where similar names exist, indicating a possible typo
type ErrPossibleTypo struct {
	similars []string
}

// Error formats a default string with the similar names found
func (e *ErrPossibleTypo) Error() string {
	if len(e.similars) == 1 {
		return fmt.Sprintf("did you mean '%s'?", e.similars[0])
	}

	return fmt.Sprintf("did you mean one of ['%s']?", strings.Join(e.similars, "', '"))
}

// Similars get the possible similar scripts
func (e *ErrPossibleTypo) Similars() []string {
	return e.similars
}

// SetSimilars set the possible similar scripts
func (e *ErrPossibleTypo) SetSimilars(similars []string) {
	e.similars = similars
}

// IsPossibleTypoError tells whether the given error is parser.ErrPossibleTypo
func IsPossibleTypoError(err error) (assert bool) {
	_, assert = err.(*ErrPossibleTypo)
	return
}
