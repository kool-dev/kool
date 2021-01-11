package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"kool-dev/kool/environment"
	"net/http"
	"net/url"
	"strings"
)

// Endpoint interface encapsulates the behaviour necessary for consuming
// an API endpoint
type Endpoint interface {
	SetPath(string)
	SetResponseReceiver(interface{})
	DoCall() error
	Body() url.Values
	Query() url.Values
	SetRawBody(io.Reader)
	SetContentType(string)
	StatusCode() int
}

// DefaultEndpoint holds common data and logic for making API calls
type DefaultEndpoint struct {
	method, path string
	contentType  string
	response     interface{}
	query, body  url.Values
	rawBody      io.Reader
	env          environment.EnvStorage
	statusCode   int
}

func newDefaultEndpoint(method string) *DefaultEndpoint {
	return &DefaultEndpoint{
		method: method,
		query:  url.Values{},
		env:    environment.NewEnvStorage(),
	}
}

// SetPath sets the URL path to be called
func (e *DefaultEndpoint) SetPath(path string) {
	e.path = path
}

// SetRawBody sets the request body to POST
func (e *DefaultEndpoint) SetRawBody(rawBody io.Reader) {
	e.rawBody = rawBody
}

// SetContentType sets the body content type
func (e *DefaultEndpoint) SetContentType(contentType string) {
	e.contentType = contentType
}

// Body sets the URL path to be called
func (e *DefaultEndpoint) Body() url.Values {
	if e.body == nil {
		e.body = url.Values{}
	}

	return e.body
}

// StatusCode returns the latest HTTP response status code
func (e *DefaultEndpoint) StatusCode() int {
	return e.statusCode
}

// Query exposes the query string builder for setting parameters
func (e *DefaultEndpoint) Query() url.Values {
	return e.query
}

// SetResponseReceiver sets the points to use for parsing the response
func (e *DefaultEndpoint) SetResponseReceiver(r interface{}) {
	e.response = r
}

// DoCall performs the actual request against the API
func (e *DefaultEndpoint) DoCall() (err error) {
	var (
		request *http.Request
		resp    *http.Response
		raw     []byte
		body    io.Reader
	)

	if e.method == "POST" {
		if e.rawBody != nil {
			body = e.rawBody
		} else if e.body != nil {
			body = strings.NewReader(e.body.Encode())
			e.contentType = "application/x-www-form-urlencoded"
		}
	}

	reqURL := fmt.Sprintf("%s/%s", apiBaseURL, e.path)
	if request, err = http.NewRequest(e.method, reqURL, body); err != nil {
		return
	}

	if e.contentType != "" {
		request.Header.Add("Content-type", e.contentType)
	}

	if resp, err = e.doRequest(request); err != nil {
		return
	}

	request = nil
	defer resp.Body.Close()

	e.statusCode = resp.StatusCode

	if raw, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	if err = json.Unmarshal(raw, e.response); err != nil {
		err = fmt.Errorf("%v (parse error: %v", ErrUnexpectedResponse, err)
	}

	return
}

func (e *DefaultEndpoint) doRequest(request *http.Request) (resp *http.Response, err error) {
	var apiToken string = e.env.Get("KOOL_API_TOKEN")

	if apiToken == "" {
		err = ErrMissingToken
		return
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "Bearer "+apiToken)

	resp, err = http.DefaultClient.Do(request)

	return
}
