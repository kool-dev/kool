package environment

import (
	"os/user"
	"regexp"
)

var sidExp = regexp.MustCompile(`.*-(?P<uid>\d+)`)

func initUid(envStorage EnvStorage) {
	// under native windows defaults to using
	// root inside containers for kool managed images
	envStorage.Set("UID", uid())
}

func uid() string {
	current, _ := user.Current()
	match := sidExp.FindStringSubmatch(current.Uid)

	results := map[string]string{}
	for i, name := range match {
		results[sidExp.SubexpNames()[i]] = name
	}

	if uid, ok := results["uid"]; ok {
		return uid
	}

	return "0"
}
