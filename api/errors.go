package api

import (
	"errors"
	"fmt"
)

// ErrBadAPIServer represents some issue in the API side
var ErrBadAPIServer error

// ErrMissingToken is the lack of the access token
var ErrMissingToken error

// ErrDeployFailed is returned when checking the status of a failed deploy
var ErrDeployFailed error

// ErrUnauthorized unauthorized; please check your KOOL_API_TOKEN
var ErrUnauthorized error

// ErrPayloadValidation something went wrong validating the payload
var ErrPayloadValidation error

// ErrBadResponseStatus unexpected return status
var ErrBadResponseStatus error

// ErrUnexpectedResponse bad API response; please ask for support
var ErrUnexpectedResponse error

// ErrAPI reprents a default error returned from the API
type ErrAPI struct {
	// {"message":"The given data was invalid.","errors":{"domain":["The domain field is required."]}}
	Status  int
	Message string `json:"message"`

	Errors map[string]interface{} `json:"errors"`
}

// Error returns the string representation for the error
func (e *ErrAPI) Error() string {
	return fmt.Sprintf("%d - %s (%v)", e.Status, e.Message, e.Errors)
}

func init() {
	ErrBadAPIServer = errors.New("bad API server response")
	ErrDeployFailed = errors.New("deploy process has failed")
	ErrUnauthorized = errors.New("unauthorized; please check your KOOL_API_TOKEN")
	ErrPayloadValidation = errors.New("something went wrong validating the payload")
	ErrBadResponseStatus = errors.New("unexpected return status")
	ErrUnexpectedResponse = errors.New("bad API response; please ask for support")
	ErrMissingToken = errors.New("missing KOOL_API_TOKEN")
}
