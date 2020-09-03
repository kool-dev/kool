package enviroment

import (
	"os"
)

func InitAsuser() {
	os.Setenv("KOOL_ASUSER", "0")
}
