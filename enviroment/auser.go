// +build !windows

package enviroment

import (
	"fmt"
	"os"
)

func InitAuser() {
	if os.Getenv("KOOL_ASUSER") == "" {
		os.Setenv("KOOL_ASUSER", fmt.Sprintf("%d", os.Getuid()))
	}
}
