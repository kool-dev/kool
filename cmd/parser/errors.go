package parser

import (
	"errors"
	"fmt"
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
