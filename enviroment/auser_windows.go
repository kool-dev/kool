package enviroment

import (
	"os"
)

func InitAuser() {
	os.Setenv("KOOL_ASUSER", "0")
}
