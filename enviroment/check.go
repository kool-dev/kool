package enviroment

import "os"

// IsTrue checks whether the given environment variable is
// to what would be a boolean value of true (either 1 or "true")
func IsTrue(env string) bool {
	verbose := os.Getenv(env)
	return verbose == "1" || verbose == "true"
}
