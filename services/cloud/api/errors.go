package api

import (
	"errors"
	"fmt"
	"strings"
)

// ErrBadAPIServer represents some issue in the API side
var ErrBadAPIServer error

// ErrMissingToken is the lack of the access token
var ErrMissingToken error

// ErrDeployFailed is returned when checking the status of a failed deploy
var ErrDeployFailed error

// ErrUnauthorized unauthorized; please check your KOOL_API_TOKEN
var ErrUnauthorized error

// failed validating deploy payload
var ErrPayloadValidation error

// ErrBadResponseStatus unexpected return status
var ErrBadResponseStatus error

// ErrUnexpectedResponse bad API response; please ask for support
var ErrUnexpectedResponse error

// ErrAPI reprents a default error returned from the API
type ErrAPI struct {
	Status int

	Message string                 `json:"message"`
	Errors  map[string]interface{} `json:"errors"`
}

// Error returns the string representation for the error
func (e *ErrAPI) Error() string {
	if e.Errors != nil {
		s := []string{}
		for k, e := range e.Errors {
			s = append(s, fmt.Sprintf("\t%s > %v", k, e.([]interface{})[0]))
		}
		return fmt.Sprintf("\n%d - %s\n\n%s\n", e.Status, e.Message, strings.Join(s, "\n"))
	}
	return fmt.Sprintf("\n%d - %s\n", e.Status, e.Message)
}

func init() {
	ErrBadAPIServer = errors.New("bad API server response")
	ErrDeployFailed = errors.New("deploy process has failed")
	ErrUnauthorized = errors.New("unauthorized; please check your KOOL_API_TOKEN")
	ErrPayloadValidation = errors.New("failed validating deploy payload")
	ErrBadResponseStatus = errors.New("unexpected return status")
	ErrUnexpectedResponse = errors.New("bad API response; please ask for support")
	ErrMissingToken = errors.New("missing KOOL_API_TOKEN")
}
