package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Endpoint holds common data and logic for making API calls
type Endpoint struct {
	method, url string
	response    interface{}
}

func newEndpoint(method string) *Endpoint {
	return &Endpoint{method: method}
}

// SetURL sets the URL to be called
func (e *Endpoint) SetURL(url string) {
	e.url = url
}

// SetResponseReceiver sets the points to use for parsing the response
func (e *Endpoint) SetResponseReceiver(r interface{}) {
	e.response = r
}

// Call performs the request against the API
func (e *Endpoint) Call() (err error) {
	var (
		request *http.Request
		resp    *http.Response
		raw     []byte
	)

	if request, err = http.NewRequest(
		"GET",
		e.url,
		nil,
	); err != nil {
		return
	}

	if resp, err = doRequest(request); err != nil {
		return
	}

	request = nil
	defer resp.Body.Close()

	if raw, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	if err = json.Unmarshal(raw, e.response); err != nil {
		err = fmt.Errorf("%v (parse error: %v", ErrUnexpectedResponse, err)
	}

	return
}
