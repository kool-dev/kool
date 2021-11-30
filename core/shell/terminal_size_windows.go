//go:build !windows
// +build !windows

package shell

// GetTerminalWidth checks if input is a terminal
func GetTerminalWidth(tty interface{}) (width int, err error) {
	return standardTermWidth, nil
}
