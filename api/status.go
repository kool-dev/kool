package api

import (
	"fmt"
)

// StatusCall interface represents logic for consuming the deploy/status API endpoint
type StatusCall interface {
	Endpoint

	Call() (*StatusResponse, error)
}

// DefaultStatusCall holds data and logic for consuming the "status" endpoint
type DefaultStatusCall struct {
	Endpoint

	deployID string
}

// StatusResponse holds data from the "status" endpoint
type StatusResponse struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

// NewDefaultStatusCall creates a new caller for Deploy API status endpoint
func NewDefaultStatusCall(deployID string) *DefaultStatusCall {
	return &DefaultStatusCall{
		Endpoint: NewDefaultEndpoint("GET"),

		deployID: deployID,
	}
}

// Call performs the request to the endpoint
func (s *DefaultStatusCall) Call() (r *StatusResponse, err error) {
	r = &StatusResponse{}

	s.SetPath(fmt.Sprintf("deploy/%s/status", s.deployID))
	s.SetResponseReceiver(r)

	err = s.Endpoint.DoCall()

	return
}
