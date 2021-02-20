package api

import (
	"kool-dev/kool/environment"
	"net/http"
	"testing"
)

func TestNewDefaultDestroyCall(t *testing.T) {
	e := NewDefaultDestroyCall()

	if e.Endpoint.(*DefaultEndpoint).method != "DELETE" {
		t.Errorf("bad method for destroy call")
	}
}

func TestDestroyCall(t *testing.T) {
	e := NewDefaultDestroyCall()
	e.Endpoint.(*DefaultEndpoint).env = environment.NewFakeEnvStorage()
	e.Endpoint.(*DefaultEndpoint).env.Set("KOOL_API_TOKEN", "fake token")

	oldHTTPRequester := httpRequester
	defer func() {
		httpRequester = oldHTTPRequester
	}()
	httpRequester = &fakeHTTP{resp: &http.Response{StatusCode: 200, Body: &fakeIOReaderCloser{
		fakeIOReader: fakeIOReader{data: []byte(`{"environment":{"id":100}}`)},
	}}}

	resp, err := e.Call()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if e.Endpoint.(*DefaultEndpoint).path != "deploy" {
		t.Errorf("bad path: %s", e.Endpoint.(*DefaultEndpoint).path)
	}

	if resp.Environment.ID != 100 {
		t.Errorf("failed parsing proper response: %d", resp.Environment.ID)
	}
}
