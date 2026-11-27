package api

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestDeployStatusResponseDecodesPanelShape(t *testing.T) {
	raw := []byte(`{
		"id": 42,
		"status": "failed",
		"type": "local",
		"url": "https://app.example.com",
		"error": "Error building images",
		"started_at": "2026-11-27T10:00:00+00:00",
		"completed_at": "2026-11-27T10:05:00+00:00",
		"services": [{"name": "app", "type": "web", "size": "small"}],
		"wordpress_credentials": null
	}`)

	resp := &DeployStatusResponse{}
	if err := json.Unmarshal(raw, resp); err != nil {
		t.Fatalf("unexpected error decoding status: %v", err)
	}

	if resp.ID != 42 || resp.Status != "failed" || resp.Type != "local" || resp.Url != "https://app.example.com" {
		t.Errorf("bad decoded status: %+v", resp)
	}

	if resp.Error != "Error building images" || resp.Hint != "" || resp.Detail != "" {
		t.Errorf("bad decoded error fields: %+v", resp)
	}

	if resp.StartedAt != "2026-11-27T10:00:00+00:00" || resp.CompletedAt != "2026-11-27T10:05:00+00:00" {
		t.Errorf("bad decoded timestamps: %+v", resp)
	}

	if len(resp.Services) != 1 || resp.Services[0] != (DeployStatusService{"app", "web", "small"}) {
		t.Errorf("bad decoded services: %+v", resp.Services)
	}

	// while running, the panel sends nulls for the fields it does not know yet
	resp = &DeployStatusResponse{}
	if err := json.Unmarshal([]byte(`{"id": 42, "status": "deploying", "url": null, "error": null, "completed_at": null, "services": []}`), resp); err != nil {
		t.Fatalf("unexpected error decoding status with nulls: %v", err)
	}

	if resp.Status != "deploying" || resp.Url != "" || resp.Error != "" || resp.CompletedAt != "" {
		t.Errorf("bad decoded running status: %+v", resp)
	}
}

// TestDeployStatusDecodesEveryPanelShape goes through DoCall with every shape
// the deploy status endpoint can answer with: the released panel's (without
// hint and detail) and the panel's that explains failures (with them).
func TestDeployStatusDecodesEveryPanelShape(t *testing.T) {
	const (
		startedAt   = "2026-11-27T10:00:00+00:00"
		completedAt = "2026-11-27T10:05:00+00:00"
	)

	cases := []struct {
		name     string
		raw      string
		expected DeployStatusResponse
	}{
		{
			name: "running, nothing known yet",
			raw: `{"id": 7, "status": "building", "type": "deploy", "url": "https://app.example.com",
				"error": null, "started_at": "` + startedAt + `", "completed_at": null,
				"services": [], "wordpress_credentials": null}`,
			expected: DeployStatusResponse{ID: 7, Status: "building", Type: "deploy", Url: "https://app.example.com",
				StartedAt: startedAt, Services: []DeployStatusService{}},
		},
		{
			name: "legacy deploy without a type",
			raw: `{"id": 7, "status": "deploying", "type": null, "url": "https://app.example.com",
				"error": null, "started_at": "` + startedAt + `", "completed_at": null,
				"services": [], "wordpress_credentials": null}`,
			expected: DeployStatusResponse{ID: 7, Status: "deploying", Url: "https://app.example.com",
				StartedAt: startedAt, Services: []DeployStatusService{}},
		},
		{
			name: "success, with a service inheriting the default size",
			raw: `{"id": 7, "status": "success", "type": "deploy", "url": "https://app.example.com",
				"error": null, "started_at": "` + startedAt + `", "completed_at": "` + completedAt + `",
				"services": [{"name": "app", "type": "service", "size": "2xlarge"}, {"name": "worker", "type": "daemon", "size": null}],
				"wordpress_credentials": null}`,
			expected: DeployStatusResponse{ID: 7, Status: "success", Type: "deploy", Url: "https://app.example.com",
				StartedAt: startedAt, CompletedAt: completedAt,
				Services: []DeployStatusService{{"app", "service", "2xlarge"}, {"worker", "daemon", ""}}},
		},
		{
			name: "WordPress success, with its credentials",
			raw: `{"id": 7, "status": "success", "type": "deploy", "url": "https://wp.example.com",
				"error": null, "started_at": "` + startedAt + `", "completed_at": "` + completedAt + `",
				"services": [{"name": "wordpress", "type": "service", "size": "small"}],
				"wordpress_credentials": {"url": "https://wp.example.com/wp-admin", "username": "admin", "password": "not-a-real-one"}}`,
			expected: DeployStatusResponse{ID: 7, Status: "success", Type: "deploy", Url: "https://wp.example.com",
				StartedAt: startedAt, CompletedAt: completedAt,
				Services: []DeployStatusService{{"wordpress", "service", "small"}}},
		},
		{
			name: "failed, with no error found in the log",
			raw: `{"id": 7, "status": "failed", "type": "deploy", "url": "https://app.example.com",
				"error": null, "started_at": "` + startedAt + `", "completed_at": "` + completedAt + `",
				"services": [], "wordpress_credentials": null}`,
			expected: DeployStatusResponse{ID: 7, Status: "failed", Type: "deploy", Url: "https://app.example.com",
				StartedAt: startedAt, CompletedAt: completedAt, Services: []DeployStatusService{}},
		},
		{
			name: "failed, with the last error line",
			raw: `{"id": 7, "status": "failed", "type": "deploy", "url": "https://app.example.com",
				"error": "UPGRADE FAILED: context deadline exceeded", "started_at": "` + startedAt + `", "completed_at": "` + completedAt + `",
				"services": [], "wordpress_credentials": null}`,
			expected: DeployStatusResponse{ID: 7, Status: "failed", Type: "deploy", Url: "https://app.example.com",
				Error: "UPGRADE FAILED: context deadline exceeded", StartedAt: startedAt, CompletedAt: completedAt,
				Services: []DeployStatusService{}},
		},
		{
			name: "timed out, with the last error line",
			raw: `{"id": 7, "status": "timedout", "type": "deploy", "url": "https://app.example.com",
				"error": "Timeout building images for deploy.", "started_at": "` + startedAt + `", "completed_at": "` + completedAt + `",
				"services": [], "wordpress_credentials": null}`,
			expected: DeployStatusResponse{ID: 7, Status: "timedout", Type: "deploy", Url: "https://app.example.com",
				Error: "Timeout building images for deploy.", StartedAt: startedAt, CompletedAt: completedAt,
				Services: []DeployStatusService{}},
		},
		{
			name:     "without a services list",
			raw:      `{"id": 7, "status": "creating", "started_at": "` + startedAt + `"}`,
			expected: DeployStatusResponse{ID: 7, Status: "creating", StartedAt: startedAt},
		},
		{
			name: "failed, explained with error, hint and detail",
			raw: `{"id": 7, "status": "failed", "type": "deploy", "url": "https://app.example.com",
				"error": "UPGRADE FAILED: pre-upgrade hooks failed: 1 error occurred:",
				"hint": "A release hook job failed: check its output.",
				"detail": "job app-before0 failed: BackoffLimitExceeded",
				"started_at": "` + startedAt + `", "completed_at": "` + completedAt + `",
				"services": [{"name": "app", "type": "service", "size": "medium"}], "wordpress_credentials": null}`,
			expected: DeployStatusResponse{ID: 7, Status: "failed", Type: "deploy", Url: "https://app.example.com",
				Error:     "UPGRADE FAILED: pre-upgrade hooks failed: 1 error occurred:",
				Hint:      "A release hook job failed: check its output.",
				Detail:    "job app-before0 failed: BackoffLimitExceeded",
				StartedAt: startedAt, CompletedAt: completedAt,
				Services: []DeployStatusService{{"app", "service", "medium"}}},
		},
		{
			name: "failed, explained with an error only",
			raw: `{"id": 7, "status": "failed", "type": "deploy", "url": "https://app.example.com",
				"error": "Error building images: exit status 1", "hint": null, "detail": null,
				"started_at": "` + startedAt + `", "completed_at": "` + completedAt + `",
				"services": [], "wordpress_credentials": null}`,
			expected: DeployStatusResponse{ID: 7, Status: "failed", Type: "deploy", Url: "https://app.example.com",
				Error: "Error building images: exit status 1", StartedAt: startedAt, CompletedAt: completedAt,
				Services: []DeployStatusService{}},
		},
		{
			name: "success, with no failure to explain",
			raw: `{"id": 7, "status": "success", "type": "deploy", "url": "https://app.example.com",
				"error": null, "hint": null, "detail": null,
				"started_at": "` + startedAt + `", "completed_at": "` + completedAt + `",
				"services": [], "wordpress_credentials": null}`,
			expected: DeployStatusResponse{ID: 7, Status: "success", Type: "deploy", Url: "https://app.example.com",
				StartedAt: startedAt, CompletedAt: completedAt, Services: []DeployStatusService{}},
		},
	}

	oldHTTPRequester := httpRequester
	defer func() { httpRequester = oldHTTPRequester }()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			httpRequester = &fakeHTTP{resp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(tc.raw)),
			}}

			endpoint := newFakeDefaultEndpoint("GET")
			endpoint.env.Set("KOOL_API_TOKEN", "fake token")

			status := &DeployStatus{Endpoint: endpoint}
			status.SetPath("deploy/7/status")

			resp, err := status.Run()
			if err != nil {
				t.Fatalf("unexpected error decoding status: %v", err)
			}

			if !reflect.DeepEqual(*resp, tc.expected) {
				t.Errorf("bad decoded status:\n got: %+v\nwant: %+v", *resp, tc.expected)
			}
		})
	}
}
