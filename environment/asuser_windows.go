package environment

import (
	"os"
)

func initAsuser(envStorage EnvStorage) {
	// under native windows defaults to using
	// root inside containers for kool managed images
	envStorage.Set("KOOL_ASUSER", "0")
}
