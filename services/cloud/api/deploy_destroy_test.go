package api

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDeployDestroyCallsDeleteDeploy(t *testing.T) {
	oldHTTPRequester := httpRequester
	oldBaseURL := apiBaseURL
	defer func() {
		httpRequester = oldHTTPRequester
		apiBaseURL = oldBaseURL
	}()

	apiBaseURL = "https://kool.test/api"

	// the API answers with the deleting deploy, its environment included
	fake := &fakeHTTP{resp: &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(`{"id": 12, "environment_id": 34, "status": "deleting",
			"url": "https://app.example.com", "project": "app",
			"environment": {"id": 34, "domain": "app.example.com", "project": {"id": 5, "name": "app"}}}`)),
	}}
	httpRequester = fake

	endpoint := newFakeDefaultEndpoint("DELETE")
	endpoint.env.Set("KOOL_API_TOKEN", "fake token")

	destroy := NewDeployDestroy()
	endpoint.SetPath(destroy.Endpoint.(*DefaultEndpoint).path)
	destroy.Endpoint = endpoint
	destroy.Query().Set("domain", "app.example.com")

	resp, err := destroy.Call()
	if err != nil {
		t.Fatalf("unexpected error destroying: %v", err)
	}

	if fake.calledRequest.Method != http.MethodDelete {
		t.Errorf("bad method: %s", fake.calledRequest.Method)
	}

	if url := fake.calledRequest.URL.String(); url != "https://kool.test/api/deploy?domain=app.example.com" {
		t.Errorf("bad destroy URL: %s", url)
	}

	if resp.Environment.ID != 34 {
		t.Errorf("bad decoded environment ID: %d", resp.Environment.ID)
	}
}
