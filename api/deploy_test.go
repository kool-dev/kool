package api

import (
	"errors"
	"kool-dev/kool/environment"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func mockHTTPRequester(status int, body string, err error) {
	httpRequester = &fakeHTTP{resp: &http.Response{StatusCode: status, Body: &fakeIOReaderCloser{
		fakeIOReader: fakeIOReader{data: []byte(body), err: err},
	}}}
}

func TestNewDeploy(t *testing.T) {
	var tarball = "tarball"
	d := NewDeploy(tarball)

	if _, ok := d.env.(*environment.DefaultEnvStorage); !ok {
		t.Error("unexpected default environment.EnvStorage")
	}

	if _, ok := d.Endpoint.(*DefaultEndpoint); !ok {
		t.Error("unexpected default Endpoint")
	}

	if tarball != d.tarballPath {
		t.Error("failed setting tarballPath")
	}

	var id = "id"
	d.id = id

	if id != d.GetID() {
		t.Error("failed setting id")
	}

	var url = "url"
	d.Status = &StatusResponse{Status: "success", URL: url}

	if !d.IsSuccessful() {
		t.Error("failed asserting success")
	}

	if url != d.GetURL() {
		t.Error("failed getting URL")
	}
}

func TestSendFile(t *testing.T) {
	tarball := filepath.Join(t.TempDir(), "test.tgz")
	_ = os.WriteFile(tarball, []byte("test"), os.ModePerm)

	d := NewDeploy(tarball)

	d.Endpoint.(*DefaultEndpoint).env = environment.NewFakeEnvStorage()
	d.Endpoint.(*DefaultEndpoint).env.Set("KOOL_API_TOKEN", "fake token")
	d.Endpoint.(*DefaultEndpoint).env.Set("KOOL_DEPLOY_DOMAIN", "foo")
	d.Endpoint.(*DefaultEndpoint).env.Set("KOOL_DEPLOY_DOMAIN_EXTRAS", "bar")
	d.Endpoint.(*DefaultEndpoint).env.Set("KOOL_DEPLOY_WWW_REDIRECT", "zim")

	oldHTTPRequester := httpRequester
	defer func() {
		httpRequester = oldHTTPRequester
	}()
	mockHTTPRequester(200, `{"id":100}`, nil)

	if err := d.SendFile(); err != nil {
		t.Errorf("unexpected error from SendFile: %v", err)
	}
	if d.Endpoint.(*DefaultEndpoint).path != "deploy/create" {
		t.Errorf("unexpected path: %s", d.Endpoint.(*DefaultEndpoint).path)
	}
	if d.id != "100" {
		t.Errorf("unexpected id: %s", d.id)
	}

	mockHTTPRequester(200, `{"id":0}`, nil)
	if err := d.SendFile(); err == nil || !strings.Contains(err.Error(), "unexpected API response") {
		t.Errorf("unexpected error from SendFile: %v", err)
	}

	mockHTTPRequester(401, `{"id":0}`, nil)
	if err := d.SendFile(); err == nil || !errors.Is(err, ErrUnauthorized) {
		t.Errorf("unexpected error from SendFile (ErrUnauthorized): %v", err)
	}
	mockHTTPRequester(422, `{"id":0}`, nil)
	if err := d.SendFile(); err == nil || !errors.Is(err, ErrPayloadValidation) {
		t.Errorf("unexpected error from SendFile (ErrPayloadValidation): %v", err)
	}
	mockHTTPRequester(500, `{"id":0}`, nil)
	if err := d.SendFile(); err == nil || !errors.Is(err, ErrBadResponseStatus) {
		t.Errorf("unexpected error from SendFile (ErrBadResponseStatus): %v", err)
	}
}

func TestFetchLatestStatus(t *testing.T) {
	d := NewDeploy("tarball")
	d.id = "100"

	d.Endpoint.(*DefaultEndpoint).env.Set("KOOL_API_TOKEN", "fake token")

	oldHTTPRequester := httpRequester
	defer func() {
		httpRequester = oldHTTPRequester
	}()

	mockHTTPRequester(200, `{"status":"foo"}`, nil)
	if err := d.FetchLatestStatus(); err != nil {
		t.Errorf("unexpected error from FetchLatestStatus: %v", err)
	}

	mockHTTPRequester(200, `{"status":"failed"}`, nil)
	if err := d.FetchLatestStatus(); !errors.Is(err, ErrDeployFailed) {
		t.Errorf("unexpected error from FetchLatestStatus (ErrDeployFailed): %v", err)
	}
	mockHTTPRequester(500, `bad response`, nil)
	if err := d.FetchLatestStatus(); !strings.Contains(err.Error(), "bad API response") {
		t.Errorf("unexpected error from FetchLatestStatus (bad API response): %v", err)
	}
}
