package api

import (
	"kool-dev/kool/core/environment"
	"net/http"
	"testing"
)

func TestNewDefaultStatusCall(t *testing.T) {
	s := NewDefaultStatusCall("foo")

	if s.Endpoint.(*DefaultEndpoint).method != "GET" {
		t.Errorf("bad method for status call: %s", s.Endpoint.(*DefaultEndpoint).method)
	}

	if s.deployID != "foo" {
		t.Errorf("failure setting deployID: %s", s.deployID)
	}
}

func TestStatusCall(t *testing.T) {
	e := NewDefaultStatusCall("foo")
	e.Endpoint.(*DefaultEndpoint).env = environment.NewFakeEnvStorage()
	e.Endpoint.(*DefaultEndpoint).env.Set("KOOL_API_TOKEN", "fake token")

	oldHTTPRequester := httpRequester
	defer func() {
		httpRequester = oldHTTPRequester
	}()
	httpRequester = &fakeHTTP{resp: &http.Response{StatusCode: 200, Body: &fakeIOReaderCloser{
		fakeIOReader: fakeIOReader{data: []byte(`{"status":"foo","url":"bar"}`)},
	}}}

	resp, err := e.Call()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if e.Endpoint.(*DefaultEndpoint).path != "deploy/foo/status" {
		t.Errorf("bad path: %s", e.Endpoint.(*DefaultEndpoint).path)
	}

	if resp.Status != "foo" {
		t.Errorf("failed parsing proper response: Status %s", resp.Status)
	}
	if resp.URL != "bar" {
		t.Errorf("failed parsing proper response: URL %s", resp.URL)
	}
}
