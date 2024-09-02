//go:build !linux
// +build !linux

package updater

func isWriteable(_ string) bool {
	return false
}
