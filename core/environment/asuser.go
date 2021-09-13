//go:build !windows
// +build !windows

package environment

func initAsuser(envStorage EnvStorage) {
	if envStorage.Get("KOOL_ASUSER") == "" {
		envStorage.Set("KOOL_ASUSER", uid())
	}
}
