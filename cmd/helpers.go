package cmd

import "os"

func dockerComposeDefaultArgs() []string {
	return []string{"-p", os.Getenv("KOOL_NAME")}
}
