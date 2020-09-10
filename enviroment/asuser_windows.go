package enviroment

import (
	"os"
)

func initAsuser() {
	// under native windows defaults to using
	// root inside containers for kool managed images
	os.Setenv("KOOL_ASUSER", "0")
}
