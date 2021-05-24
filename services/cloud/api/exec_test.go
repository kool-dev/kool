package api

import (
	"kool-dev/kool/core/environment"
	"net/http"
	"testing"
)

func TestNewDefaultExecCall(t *testing.T) {
	e := NewDefaultExecCall()
	if e.Endpoint.(*DefaultEndpoint).method != "POST" {
		t.Errorf("bad method: %v", e.Endpoint.(*DefaultEndpoint).method)
	}
}

func TestExecCall(t *testing.T) {
	e := NewDefaultExecCall()
	e.Endpoint.(*DefaultEndpoint).env = environment.NewFakeEnvStorage()
	e.Endpoint.(*DefaultEndpoint).env.Set("KOOL_API_TOKEN", "fake token")

	oldHTTPRequester := httpRequester
	defer func() {
		httpRequester = oldHTTPRequester
	}()
	httpRequester = &fakeHTTP{resp: &http.Response{StatusCode: 200, Body: &fakeIOReaderCloser{
		fakeIOReader: fakeIOReader{data: []byte(`{"server":"server","namespace":"ns","ca.crt":"ca"}`)},
	}}}

	resp, err := e.Call()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if e.Endpoint.(*DefaultEndpoint).path != "deploy/exec" {
		t.Errorf("bad path: %s", e.Endpoint.(*DefaultEndpoint).path)
	}

	if resp.Server != "server" || resp.Namespace != "ns" || resp.CA != "ca" {
		t.Errorf("failed parsing proper response: %+v", resp)
	}
}
