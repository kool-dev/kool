package api

var (
	// apiBaseURL string = "https://kool.dev/api"
	apiBaseURL string = "http://kool.localhost/api"
)

// SetBaseURL defines the target Kool API URL to be used
// when reaching out endpoints.
func SetBaseURL(url string) {
	apiBaseURL = url
}
