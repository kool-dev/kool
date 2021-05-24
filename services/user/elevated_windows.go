package user

import "os"

// CurrentUserIsElevated returns true if the current user
// executing the program has admin privileges (sudo/administrator)
// Reference: https://gist.github.com/jerblack/d0eb182cc5a1c1d92d92a4c4fcc416c6
func CurrentUserIsElevated() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}
