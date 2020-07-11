package api

import (
	"net/http"
	"os"
)

var (
	apiBaseURL string = "https://kool.dev/api"
)

// SetBaseURL defines the target Kool API URL to be used
// when reaching out endpoints.
func SetBaseURL(url string) {
	apiBaseURL = url
}

func doRequest(request *http.Request) (resp *http.Response, err error) {
	var apiToken string = os.Getenv("KOOL_API_TOKEN")

	if apiToken == "" {
		err = ErrMissingToken
		return
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "Bearer "+apiToken)

	resp, err = http.DefaultClient.Do(request)
	return
}
