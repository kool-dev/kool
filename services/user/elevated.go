package user

import "os"

// CurrentUserIsElevated returns true if the current user
// executing the program has admin privileges (sudo/administrator)
func CurrentUserIsElevated() bool {
	return os.Getuid() == 0
}
