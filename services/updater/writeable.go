//go:build !linux

package updater

import "fmt"

func isWriteable(path string) bool {
	fmt.Println("called unimplemented updater.isWriteable")

	return false
}
