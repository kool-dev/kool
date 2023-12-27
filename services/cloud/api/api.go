package api

var (
	apiBaseURL string = "https://kool.dev/api"
)

// SetBaseURL defines the target Kool API URL to be used
// when reaching out endpoints.
func SetBaseURL(url string) {
	apiBaseURL = url
}
