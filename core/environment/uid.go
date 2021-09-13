//go:build !windows
// +build !windows

package environment

import (
	"fmt"
	"os"
)

func initUid(envStorage EnvStorage) {
	if envStorage.Get("UID") == "" {
		envStorage.Set("UID", uid())
	}
}

func uid() string {
	return fmt.Sprintf("%d", os.Getuid())
}
