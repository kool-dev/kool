package api

import (
	"errors"
	"kool-dev/kool/environment"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func newFakeDefaultEndpoint(method string) *DefaultEndpoint {
	return &DefaultEndpoint{
		method: method,
		query:  url.Values{},
		env:    environment.NewFakeEnvStorage(),
	}
}

func TestNewDefaultEndpoint(t *testing.T) {
	var method = "GET"
	e := newDefaultEndpoint(method)

	if e.method != method {
		t.Error("unexpected method")
	}

	if _, ok := e.env.(*environment.DefaultEnvStorage); !ok {
		t.Error("unexpected environment storage")
	}

	var path string = "path"
	e.SetPath(path)

	if e.path != path {
		t.Error("failed to SetPath")
	}

	var contentType string = "type"
	e.SetContentType(contentType)

	if e.contentType != contentType {
		t.Error("failed to SetContentType")
	}

	if e.rawBody != nil {
		t.Error("unexpected default rawBody")
	}

	e.SetRawBody(new(fakeIOReader))

	if _, ok := e.rawBody.(*fakeIOReader); !ok {
		t.Error("failed setting rawBody")
	}

	if e.body != nil {
		t.Error("unexpected non-null default body")
	}

	e.Body().Set("foo", "bar")
	e.Body().Set("foo2", "bar2")

	if e.body == nil {
		t.Error("body must not be null after access")
	}

	if e.StatusCode() != e.statusCode || e.statusCode != 0 {
		t.Error("unexpected default statusCode")
	}

	e.Query().Add("foo", "qbar")

	if e.query.Get("foo") != "qbar" {
		t.Error("failed to write query")
	}

	if e.response != nil {
		t.Error("unexpected default response receiver")
	}

	e.SetResponseReceiver("receiver")

	if resp, ok := e.response.(string); !ok || resp != "receiver" {
		t.Error("failed SetResponseReceiver")
	}
}

type fakeHTTP struct {
	err  error
	resp *http.Response
}

func (f *fakeHTTP) Do(r *http.Request) (*http.Response, error) {
	return f.resp, f.err
}

func TestDoRequest(t *testing.T) {
	e := newFakeDefaultEndpoint("GET")

	request, _ := http.NewRequest("GET", "http://example.com", nil)

	if _, err := e.doRequest(request); err == nil || err.Error() != ErrMissingToken.Error() {
		t.Errorf("was expecting ErrMisstingToken; got %v", err)
	}

	e.env.Set("KOOL_API_TOKEN", "fake token")

	oldHTTPRequester := httpRequester
	defer func() {
		httpRequester = oldHTTPRequester
	}()

	httpRequester = &fakeHTTP{err: errors.New("fake http error")}

	if _, err := e.doRequest(request); err == nil || err.Error() != "fake http error" {
		t.Errorf("unexpected error returned; %v", err)
	}
	if request.Header.Get("Accept") != "application/json" {
		t.Error("failed setting Accept header")
	}

	if !strings.Contains(request.Header.Get("Authorization"), "fake token") {
		t.Error("failed setting Authorization header")
	}
}
