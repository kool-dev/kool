package api

var (
	apiBaseURL string = "https://kool.dev/api"

	cliVersion string
)

// SetBaseURL defines the target Kool API URL to be used
// when reaching out endpoints.
func SetBaseURL(url string) {
	apiBaseURL = url
}

// SetCliVersion injects version to this package
func SetCliVersion(v string) {
	cliVersion = v
}
