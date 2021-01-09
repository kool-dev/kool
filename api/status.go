package api

import (
	"fmt"
)

// StatusCall holds data and logic for consuming the "exec" endpoint
type StatusCall struct {
	Endpoint

	deployID string
}

// StatusResponse holds data from the "exec" endpoint
type StatusResponse struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

// NewStatusCall creates a new caller for Deploy API exec endpoint
func NewStatusCall(deployID string) *StatusCall {
	return &StatusCall{
		Endpoint: *newEndpoint("GET"),

		deployID: deployID,
	}
}

// Call performs the request to the endpoint
func (s *StatusCall) Call() (r *StatusResponse, err error) {
	r = &StatusResponse{}

	s.Endpoint.SetURL(fmt.Sprintf("%s/deploy/%s/status", apiBaseURL, s.deployID))
	s.Endpoint.SetResponseReceiver(r)

	err = s.Endpoint.Call()

	return
}
