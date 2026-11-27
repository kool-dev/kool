package redact

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// fakeRegistryToken looks like the base64 registry password returned by /deploy/create
var fakeRegistryToken = strings.Repeat("eyJwYXlsb2FkIjoiQUFBQUJCQkJDQ0NDIn0+/Zm9vYmFy", 40) + "=="

func TestIsSensitiveKey(t *testing.T) {
	for _, key := range []string{
		"password", "DB_PASSWORD", "login", "token", "KOOL_API_TOKEN", "accessToken",
		"client_secret", "APP_KEY", "AWS_ACCESS_KEY_ID", "apiKey", "DB_PASS",
		"SENTRY_DSN", "Authorization", "--token", "aws_secret_access_key",
	} {
		if !IsSensitiveKey(key) {
			t.Errorf("expected %q to be a sensitive key", key)
		}
	}

	for _, key := range []string{
		"id", "status", "image_prefix", "image_repository", "image_tag", "url",
		"logs_url", "name", "PWD", "region", "DB_HOST", "platform",
	} {
		if IsSensitiveKey(key) {
			t.Errorf("expected %q not to be a sensitive key", key)
		}
	}
}

func TestJSONRedactsDeployCreateCredentials(t *testing.T) {
	raw := []byte(`{"deploy":{"id":44417,"url":"https://app.example.com","environment":{"name":"app","env":{"APP_ENV":"production","DB_PASSWORD":"hunter2","SOME_VALUE":"Zx9kQ2LmP4rT7vWb1NcY8sHd3FgJ6aKe0uIo"}}},` +
		`"stuff":{"image_prefix":"123456789012.dkr.ecr.us-east-1.amazonaws.com","image_repository":"kool-cloud/app","image_tag":"20261127"},` +
		`"docker":{"login":"AWS","password":"` + fakeRegistryToken + `"},"logs_url":"https://kool.dev/deploys/44417?a=1&b=2"}`)

	redacted := JSON(raw)

	for _, secret := range []string{fakeRegistryToken[:40], `"AWS"`, "hunter2", "Zx9kQ2LmP4rT7vWb1NcY8sHd3FgJ6aKe0uIo"} {
		if strings.Contains(redacted, secret) {
			t.Errorf("expected %q to be redacted; got %s", secret, redacted)
		}
	}

	var doc struct {
		Deploy struct {
			ID          int    `json:"id"`
			URL         string `json:"url"`
			Environment struct {
				Env map[string]string `json:"env"`
			} `json:"environment"`
		} `json:"deploy"`
		Stuff  map[string]string `json:"stuff"`
		Docker map[string]string `json:"docker"`
		Logs   string            `json:"logs_url"`
	}

	if err := json.Unmarshal([]byte(redacted), &doc); err != nil {
		t.Fatalf("redacted JSON must still be valid JSON: %v (%s)", err, redacted)
	}

	if doc.Docker["login"] != Placeholder || doc.Docker["password"] != Placeholder {
		t.Errorf("expected docker credentials to be redacted; got %v", doc.Docker)
	}

	if env := doc.Deploy.Environment.Env; env["DB_PASSWORD"] != Placeholder || env["SOME_VALUE"] != Placeholder || env["APP_ENV"] != "production" {
		t.Errorf("unexpected redaction of environment variables: %v", env)
	}

	// what helps debugging stays as it was
	if doc.Deploy.ID != 44417 || doc.Deploy.URL != "https://app.example.com" || doc.Logs != "https://kool.dev/deploys/44417?a=1&b=2" {
		t.Errorf("unexpected change to non-credential fields: %s", redacted)
	}

	if expected := map[string]string{"image_prefix": "123456789012.dkr.ecr.us-east-1.amazonaws.com", "image_repository": "kool-cloud/app", "image_tag": "20261127"}; !reflect.DeepEqual(doc.Stuff, expected) {
		t.Errorf("unexpected change to the image config: %v", doc.Stuff)
	}
}

func TestJSONKeepsBlankCredentials(t *testing.T) {
	if redacted := JSON([]byte(`{"token":"","password":null,"has_login":true}`)); redacted != `{"has_login":true,"password":null,"token":""}` {
		t.Errorf("expected blank credentials to be kept as they are; got %s", redacted)
	}
}

func TestJSONRedactsExecToken(t *testing.T) {
	redacted := JSON([]byte(`{"server":"https://k8s.example.com:6443","namespace":"dev-app","path":"pod/app-1","token":"short-k8s-token","ca.crt":"-----BEGIN CERTIFICATE-----\nMIIBdzCCAR2gAwIBAgIBADAKBggqhkjOPQQDAjAjMSEwHwYDVQQDDBhrM3Mtc2Vy\n-----END CERTIFICATE-----"}`))

	if strings.Contains(redacted, "short-k8s-token") || strings.Contains(redacted, "MIIBdzCCAR2gAwIBAgIBADAKBggqhkjOPQQDAjAjMSEwHwYDVQQDDBhrM3Mtc2Vy") {
		t.Errorf("expected exec credentials to be redacted; got %s", redacted)
	}

	if !strings.Contains(redacted, `"namespace":"dev-app"`) || !strings.Contains(redacted, `"server":"https://k8s.example.com:6443"`) {
		t.Errorf("unexpected change to non-credential fields: %s", redacted)
	}
}

func TestJSONFallsBackToText(t *testing.T) {
	for _, raw := range []string{
		`<html>password=hunter2</html>`,
		`{"password":"hunter2"} trailing`,
		`{"password":"hunter2"`,
	} {
		if redacted := JSON([]byte(raw)); strings.Contains(redacted, "hunter2") {
			t.Errorf("expected password to be redacted from %q; got %q", raw, redacted)
		}
	}
}

func TestText(t *testing.T) {
	tests := []struct {
		in, expected string
	}{
		{"Authorization: Bearer abc.def-123", "Authorization: [REDACTED] [REDACTED]"},
		{"got bearer abc123", "got bearer [REDACTED]"},
		{"mysql://root:hunter2@db:3306/app", "mysql://root:[REDACTED]@db:3306/app"},
		{"https://kool.dev/api/deploy/status?id=12&token=abc", "https://kool.dev/api/deploy/status?id=12&token=[REDACTED]"},
		{"DB_PASSWORD=hunter2 DB_HOST=db", "DB_PASSWORD=[REDACTED] DB_HOST=db"},
		{`{"secret": "hunter2", "id": 1}`, `{"secret": "[REDACTED]", "id": 1}`},
		{"login with " + fakeRegistryToken + " now", "login with [REDACTED] now"},
		{"key 1|Zx9kQ2LmP4rT7vWb1NcY8sHd3FgJ6aKe0uIoPq", "key 1|[REDACTED]"},
		{"hex 0123456789abcdef0123456789abcdef0123", "hex [REDACTED]"},
		// not credentials
		{"Going to call: https://kool.dev/api/deploy/create?", "Going to call: https://kool.dev/api/deploy/create?"},
		{"/Users/someone/Projects/Kool2026/apps/node-20/src", "/Users/someone/Projects/Kool2026/apps/node-20/src"},
		{"dev-my-app-example-com-2026-11-27-production", "dev-my-app-example-com-2026-11-27-production"},
		{"KOOL_DEPLOY_DOMAIN_EXTRAS_FOR_SOME_REASON_123", "KOOL_DEPLOY_DOMAIN_EXTRAS_FOR_SOME_REASON_123"},
		{"short Zx9kQ2LmP4rT", "short Zx9kQ2LmP4rT"},
		{"linux/amd64", "linux/amd64"},
	}

	for _, tt := range tests {
		if got := Text(tt.in); got != tt.expected {
			t.Errorf("Text(%q): expected %q; got %q", tt.in, tt.expected, got)
		}
	}
}

func TestArgs(t *testing.T) {
	tests := []struct {
		in, expected []string
	}{
		{
			[]string{"login", "-u", "AWS", "--password-stdin", "123456789012.dkr.ecr.us-east-1.amazonaws.com"},
			[]string{"login", "-u", "[REDACTED]", "--password-stdin", "123456789012.dkr.ecr.us-east-1.amazonaws.com"},
		},
		{
			[]string{"login", "--username=robot", "-p", "hunter2", "registry.kool.dev"},
			[]string{"login", "--username=[REDACTED]", "-p", "[REDACTED]", "registry.kool.dev"},
		},
		{
			[]string{"--server", "https://k8s.example.com", "--token", "short-k8s-token", "--namespace", "dev-app", "exec", "-i", "pod/app", "--", "bash"},
			[]string{"--server", "https://k8s.example.com", "--token", "[REDACTED]", "--namespace", "dev-app", "exec", "-i", "pod/app", "--", "bash"},
		},
		{
			[]string{"cloud", "deploy", "--token=abc", "--domain", "app.example.com"},
			[]string{"cloud", "deploy", "--token=[REDACTED]", "--domain", "app.example.com"},
		},
		{
			[]string{"build", "-t", "registry/app:web-1", "--platform", "linux/amd64", "--build-arg", "NPM_TOKEN=hunter2", "--build-arg", "NODE_ENV=production", "--build-arg=DB_PASSWORD=hunter2", "/Users/someone/app"},
			[]string{"build", "-t", "registry/app:web-1", "--platform", "linux/amd64", "--build-arg", "NPM_TOKEN=[REDACTED]", "--build-arg", "NODE_ENV=production", "--build-arg=DB_PASSWORD=[REDACTED]", "/Users/someone/app"},
		},
		{
			// -p and -u are only credentials in a login command
			[]string{"run", "-p", "8080:80", "-u", "1000", "-e", "SECRET=" + fakeRegistryToken, "image", fakeRegistryToken},
			[]string{"run", "-p", "8080:80", "-u", "1000", "-e", "SECRET=[REDACTED]", "image", "[REDACTED]"},
		},
	}

	for _, tt := range tests {
		if got := Args(tt.in); !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("Args(%v):\nexpected %v\n     got %v", tt.in, tt.expected, got)
		}
	}
}

func TestArgsDoesNotChangeInput(t *testing.T) {
	in := []string{"--token", "abc"}
	_ = Args(in)

	if in[1] != "abc" {
		t.Error("Args must not change the given slice")
	}
}
