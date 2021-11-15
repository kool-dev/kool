//go:build linux
// +build linux

package updater

import "golang.org/x/sys/unix"

func isWriteable(path string) bool {
	return unix.Access(path, unix.W_OK) == nil
}
