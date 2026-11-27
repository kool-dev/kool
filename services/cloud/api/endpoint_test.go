package api

import (
	"bytes"
	"errors"
	"io"
	"kool-dev/kool/core/environment"
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
	e := NewDefaultEndpoint(method)

	if e.method != method {
		t.Error("unexpected method")
	}

	if _, ok := e.env.(*environment.DefaultEnvStorage); !ok {
		t.Error("unexpected environment storage")
	}

	var path = "path"
	e.SetPath(path)

	if e.path != path {
		t.Error("failed to SetPath")
	}

	var contentType = "type"
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

	calledRequest *http.Request
}

func (f *fakeHTTP) Do(r *http.Request) (*http.Response, error) {
	f.calledRequest = r
	return f.resp, f.err
}

type fakeIOReader struct {
	data []byte
	err  error
}

func (r *fakeIOReader) Read(p []byte) (n int, err error) {
	copy(p, r.data)
	n = len(r.data)
	err = r.err
	if err == nil {
		err = io.EOF
	}
	return
}

type fakeIOReaderCloser struct {
	fakeIOReader

	closeErr    error
	calledClose bool
}

func (c *fakeIOReaderCloser) Close() error {
	c.calledClose = true
	return c.closeErr
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

func TestDoCallGet(t *testing.T) {
	e := newFakeDefaultEndpoint("GET")
	e.env.Set("KOOL_API_TOKEN", "fake token")

	oldHTTPRequester := httpRequester
	defer func() {
		httpRequester = oldHTTPRequester
	}()

	httpErr := errors.New("fake http error")
	httpRequester = &fakeHTTP{err: httpErr}

	apiBaseURL = "base-url"
	e.SetContentType("content-type")
	e.SetPath("path")
	e.Query().Set("foo", "bar")

	if err := e.DoCall(); !errors.Is(err, httpErr) {
		t.Errorf("unexpected error returned from DoCall: %v", err)
	}

	httpRequester.(*fakeHTTP).err = nil
	httpRequester.(*fakeHTTP).resp = &http.Response{StatusCode: 200, Body: &fakeIOReaderCloser{
		fakeIOReader: fakeIOReader{data: []byte("test bad response")},
	}}

	e.SetResponseReceiver(struct{}{})

	if err := e.DoCall(); err == nil || !strings.Contains(err.Error(), "parse error") {
		t.Errorf("unexpected error return from DoCall; expected parse error got: %v", err)
	}

	if e.StatusCode() != 200 {
		t.Errorf("bad status code %d", e.StatusCode())
	}

	httpRequester.(*fakeHTTP).resp = &http.Response{StatusCode: 400, Body: &fakeIOReaderCloser{
		fakeIOReader: fakeIOReader{data: []byte(`still bad response`)},
	}}

	e.SetResponseReceiver(struct{}{})

	if err := e.DoCall(); err == nil || !strings.Contains(err.Error(), "parse error") {
		t.Errorf("unexpected error return from DoCall; expected parse error got: %v", err)
	}

	if e.StatusCode() != 400 {
		t.Errorf("bad status code %d", e.StatusCode())
	}

	httpRequester.(*fakeHTTP).resp = &http.Response{StatusCode: 403, Body: &fakeIOReaderCloser{
		fakeIOReader: fakeIOReader{data: []byte(`{"message":"err-message"}`)},
	}}

	e.SetResponseReceiver(struct{}{})

	if err := e.DoCall(); err == nil || !strings.Contains(err.Error(), "err-message") {
		t.Errorf("unexpected error return from DoCall; expected err-message got: %v", err)
	}

	if e.StatusCode() != 403 {
		t.Errorf("bad status code %d", e.StatusCode())
	}

	e.env.Set("KOOL_VERBOSE", "1")

	httpRequester.(*fakeHTTP).resp = &http.Response{StatusCode: 200, Body: &fakeIOReaderCloser{
		fakeIOReader: fakeIOReader{data: []byte(`{"foo":"bar"}`)},
	}}

	resp := &struct{ Foo string }{}
	httpRequester.(*fakeHTTP).resp.Body.(*fakeIOReaderCloser).calledClose = false
	e.SetResponseReceiver(resp)

	if err := e.DoCall(); err != nil {
		t.Errorf("unexpected error return from DoCall: %v", err)
	}

	if e.StatusCode() != 200 {
		t.Errorf("bad status code %d", e.StatusCode())
	}
	if resp.Foo != "bar" {
		t.Errorf("response did not get parsed properly %v", resp)
	}
	if !httpRequester.(*fakeHTTP).resp.Body.(*fakeIOReaderCloser).calledClose {
		t.Errorf("response was not closed")
	}

	e = newFakeDefaultEndpoint(" ")
	if err := e.DoCall(); err == nil || !strings.Contains(err.Error(), "invalid method") {
		t.Errorf("unexpected error return from DoCall: %v", err)
	}

	e = newFakeDefaultEndpoint("GET")
	e.env.Set("KOOL_API_TOKEN", "fake token")
	httpRequester.(*fakeHTTP).resp = &http.Response{StatusCode: 200, Body: &fakeIOReaderCloser{
		fakeIOReader: fakeIOReader{err: errors.New("read body error")},
	}}
	if err := e.DoCall(); err == nil || !strings.Contains(err.Error(), "read body error") {
		t.Errorf("unexpected error return from DoCall: %v", err)
	}
}

func TestDoCallPost(t *testing.T) {
	e := newFakeDefaultEndpoint("POST")
	e.env.Set("KOOL_API_TOKEN", "fake token")

	oldHTTPRequester := httpRequester
	defer func() {
		httpRequester = oldHTTPRequester
	}()

	httpRequester = &fakeHTTP{resp: &http.Response{StatusCode: 200, Body: &fakeIOReaderCloser{
		fakeIOReader: fakeIOReader{data: []byte(`"response"`)},
	}}}

	e.SetRawBody(&fakeIOReader{})

	var resp string
	e.SetResponseReceiver(&resp)

	if err := e.DoCall(); err != nil || resp != "response" {
		t.Errorf("unexpected error; unexpected response: %v - %s", err, resp)
	}

	e.SetRawBody(nil)
	e.Body().Set("foo", "bar")

	if err := e.DoCall(); err != nil || resp != "response" {
		t.Errorf("unexpected error; unexpected response: %v - %s", err, resp)
	}

	if e.contentType != "application/x-www-form-urlencoded" {
		t.Errorf("bad contentType: %s", e.contentType)
	}
}

func TestDoCallSetsClientHeaders(t *testing.T) {
	e := newFakeDefaultEndpoint("GET")
	e.env.Set("KOOL_API_TOKEN", "fake token")
	e.env.Set("CLAUDECODE", "1")

	oldHTTPRequester := httpRequester
	oldCliVersion := cliVersion
	defer func() {
		httpRequester = oldHTTPRequester
		cliVersion = oldCliVersion
	}()

	SetCliVersion("1.2.3")
	fake := &fakeHTTP{resp: &http.Response{StatusCode: 200, Body: &fakeIOReaderCloser{
		fakeIOReader: fakeIOReader{data: []byte(`{}`)},
	}}}
	httpRequester = fake

	var resp map[string]interface{}
	e.SetResponseReceiver(&resp)

	if err := e.DoCall(); err != nil {
		t.Fatalf("unexpected error from DoCall: %v", err)
	}

	if v := fake.calledRequest.Header.Get("x-kool-cli-version"); v != "1.2.3" {
		t.Errorf("bad x-kool-cli-version header: %q", v)
	}

	if v := fake.calledRequest.Header.Get("x-kool-client"); v != "claude-code" {
		t.Errorf("bad x-kool-client header: %q", v)
	}
}

func TestDoCallVerboseRedactsCredentials(t *testing.T) {
	var (
		registryToken = strings.Repeat("eyJwYXlsb2FkIjoiQUFBQUJCQkJDQ0NDIn0+/Zm9vYmFy", 40) + "=="
		apiToken      = "123|Zx9kQ2LmP4rT7vWb1NcY8sHd3FgJ6aKe0uIoPq"
		verbose       = new(bytes.Buffer)
	)

	e := newFakeDefaultEndpoint("POST")
	e.env.Set("KOOL_API_TOKEN", apiToken)
	e.env.Set("KOOL_VERBOSE", "1")
	e.SetPath("deploy/create")
	e.Query().Set("token", apiToken)

	oldHTTPRequester, oldVerboseOutput := httpRequester, verboseOutput
	defer func() {
		httpRequester, verboseOutput = oldHTTPRequester, oldVerboseOutput
	}()

	verboseOutput = verbose
	httpRequester = &fakeHTTP{resp: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(
		`{"deploy":{"id":44417,"status":"building"},` +
			`"stuff":{"image_prefix":"123456789012.dkr.ecr.us-east-1.amazonaws.com"},` +
			`"docker":{"login":"AWS","password":"` + registryToken + `"}}`,
	))}}

	var resp DeployCreateResponse
	e.SetResponseReceiver(&resp)

	if err := e.DoCall(); err != nil {
		t.Fatalf("unexpected error from DoCall: %v", err)
	}

	// the CLI itself still gets the credentials it needs
	if resp.Docker.Login != "AWS" || resp.Docker.Password != registryToken {
		t.Error("the response must not be redacted, only what --verbose prints")
	}

	output := verbose.String()

	for _, secret := range []string{registryToken[:40], `"AWS"`, apiToken, "Zx9kQ2LmP4rT7vWb1NcY8sHd3FgJ6aKe0uIoPq"} {
		if strings.Contains(output, secret) {
			t.Errorf("--verbose printed a credential (%q): %s", secret, output)
		}
	}

	for _, expected := range []string{
		"Going to call: " + apiBaseURL + "/deploy/create?token=[REDACTED]",
		`"docker":{"login":"[REDACTED]","password":"[REDACTED]"}`,
		`"deploy":{"id":44417,"status":"building"}`,
		"123456789012.dkr.ecr.us-east-1.amazonaws.com",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected --verbose to print %q; got %s", expected, output)
		}
	}
}
