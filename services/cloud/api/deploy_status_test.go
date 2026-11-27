package api

import (
	"encoding/json"
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
