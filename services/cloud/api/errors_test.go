package api

import (
	"testing"
)

func TestDefinedErrors(t *testing.T) {
	errs := map[string]error{
		"ErrBadAPIServer":       ErrBadAPIServer,
		"ErrDeployFailed":       ErrDeployFailed,
		"ErrUnauthorized":       ErrUnauthorized,
		"ErrPayloadValidation":  ErrPayloadValidation,
		"ErrBadResponseStatus":  ErrBadResponseStatus,
		"ErrUnexpectedResponse": ErrUnexpectedResponse,
		"ErrMissingToken":       ErrMissingToken,
	}

	for e, v := range errs {
		if v == nil {
			t.Errorf("default error not defined: %s", e)
		}
	}
}

func TestApiErr(t *testing.T) {
	err := &ErrAPI{100, "message", nil}

	if err.Error() != "100 - message" {
		t.Errorf("unexpected error message: %s", err.Error())
	}

	err.Errors = make(map[string]interface{})

	if err.Error() != "100 - message (map[])" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}
