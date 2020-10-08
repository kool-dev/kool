// +build !windows

package environment

import (
	"fmt"
	"os"
)

func initAsuser(envStorage EnvStorage) {
	if envStorage.Get("KOOL_ASUSER") == "" {
		envStorage.Set("KOOL_ASUSER", fmt.Sprintf("%d", os.Getuid()))
	}
}
