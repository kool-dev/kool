package parser

import "fmt"

// ErrMultipleDefinedScript happens when the script asked for is
// found within more than one kool.yml files targeted.
var ErrMultipleDefinedScript = fmt.Errorf("script was found in more than one kool.yml file")

// IsMultipleDefinedScriptError tells whether the given error is parser.ErrMultipleDefinedScript
func IsMultipleDefinedScriptError(err error) bool {
	return err.Error() == ErrMultipleDefinedScript.Error()
}
