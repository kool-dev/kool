package proxy

import (
	"encoding/json"
	"io"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	workDir := t.TempDir()
	content := `proxy:
  domain: app.localhost
  network: shared
  routes:
    app:
      ports: ["80:8080", "8080:8081"]
      hosts: ["@", "*"]
    node:
      ports: ["3001:3001"]
`
	if err := os.WriteFile(filepath.Join(workDir, "kool.yml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	env := environment.NewFakeEnvStorage()
	env.Set("PWD", workDir)
	manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)

	config, err := manager.loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if config.Domain != "app.localhost" || config.Network != "shared" || len(config.Routes) != 3 {
		t.Fatalf("unexpected proxy config: %#v", config)
	}
	if config.Routes[0].Service != "app" || config.Routes[0].Listen != 80 || config.Routes[0].Target != 8080 || len(config.Routes[0].Hosts) != 2 {
		t.Errorf("unexpected app route: %#v", config.Routes[0])
	}
	if config.Routes[2].Service != "node" || len(config.Routes[2].Hosts) != 2 || config.Routes[2].Hosts[0] != "@" || config.Routes[2].Hosts[1] != "*" {
		t.Errorf("expected node route to default to base and wildcard hosts, got %#v", config.Routes[2])
	}
}

func TestRouteHosts(t *testing.T) {
	env := environment.NewFakeEnvStorage()
	env.Set("KOOL_PROXY_HOST", "task-a.workspace.app.localhost")
	manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)

	hosts := manager.routeHosts(route{Hosts: []string{"@", "*", "api", "admin.example.test"}})
	expected := []string{
		"task-a.workspace.app.localhost",
		"*.task-a.workspace.app.localhost",
		"api.task-a.workspace.app.localhost",
		"admin.example.test",
	}
	if strings.Join(hosts, ",") != strings.Join(expected, ",") {
		t.Errorf("expected hosts %v, got %v", expected, hosts)
	}
}

func TestLoadConfigRejectsInvalidRoute(t *testing.T) {
	workDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(workDir, "kool.yml"), []byte("proxy:\n  domain: app.localhost\n  routes:\n    app:\n      ports: [invalid]\n      hosts: ['@']\n"), 0644); err != nil {
		t.Fatal(err)
	}
	env := environment.NewFakeEnvStorage()
	env.Set("PWD", workDir)
	manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)

	if _, err := manager.loadConfig(); err == nil || !strings.Contains(err.Error(), "listen:target") {
		t.Fatalf("expected route format error, got %v", err)
	}
}

func TestLoadConfigRejectsEquivalentResolvedHosts(t *testing.T) {
	workDir := t.TempDir()
	content := "proxy:\n  domain: app.localhost\n  routes:\n    app:\n      ports: ['80:80']\n      hosts: ['@']\n    node:\n      ports: ['80:3001']\n      hosts: ['app.localhost']\n"
	if err := os.WriteFile(filepath.Join(workDir, "kool.yml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	env := environment.NewFakeEnvStorage()
	env.Set("PWD", workDir)
	env.Set("KOOL_PROXY_HOST", "app.localhost")
	manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)

	if _, err := manager.loadConfig(); err == nil || !strings.Contains(err.Error(), "both use 80:app.localhost") {
		t.Fatalf("expected equivalent host conflict, got %v", err)
	}
}

func TestCreateAliasOverride(t *testing.T) {
	workDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(workDir, "compose.yml"), []byte("services: {}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	env := environment.NewFakeEnvStorage()
	env.Set("PWD", workDir)
	env.Set("KOOL_NAME", "example")
	manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)

	file, err := manager.createAliasOverride("kool_global", []route{{Service: "app", Listen: 80, Target: 80}, {Service: "app", Listen: 8080, Target: 3000}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(file) })
	content, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "example-app") || !strings.Contains(string(content), "kool_global") {
		t.Errorf("unexpected alias override:\n%s", content)
	}
	if strings.Count(string(content), `"app":`) != 1 {
		t.Errorf("expected one service override for multiple routes, got:\n%s", content)
	}
	if !strings.Contains(env.Get("COMPOSE_FILE"), file) {
		t.Errorf("expected COMPOSE_FILE to include %s, got %s", file, env.Get("COMPOSE_FILE"))
	}
}

func TestCreateAliasOverrideMapsCustomExternalNetwork(t *testing.T) {
	workDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(workDir, "compose.yml"), []byte("services: {}\nnetworks:\n  kool_global:\n    external: true\n"), 0644); err != nil {
		t.Fatal(err)
	}
	env := environment.NewFakeEnvStorage()
	env.Set("PWD", workDir)
	env.Set("KOOL_NAME", "example")
	manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)

	file, err := manager.createAliasOverride("my_network", []route{{Service: "app", Listen: 80, Target: 80}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(file) })
	content, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{`"kool_global":`, `name: "my_network"`, "external: true"} {
		if !strings.Contains(string(content), expected) {
			t.Errorf("expected override to contain %q, got:\n%s", expected, content)
		}
	}
	if strings.Contains(string(content), `"my_network":`) {
		t.Errorf("custom network name must not be used as an undeclared Compose key:\n%s", content)
	}
}

func TestPrepareRestoresComposeFile(t *testing.T) {
	workDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(workDir, "kool.yml"), []byte("proxy:\n  domain: app.localhost\n  routes:\n    app:\n      ports: ['80:80']\n"), 0644); err != nil {
		t.Fatal(err)
	}
	env := environment.NewFakeEnvStorage()
	env.Set("PWD", workDir)
	env.Set("COMPOSE_FILE", "compose.yml:compose.dev.yml")
	manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)

	cleanup, err := manager.Prepare([]string{"app"})
	if err == nil {
		t.Fatal("expected fake shell to fail while ensuring Caddy")
	}
	cleanup()
	if value := env.Get("COMPOSE_FILE"); value != "compose.yml:compose.dev.yml" {
		t.Fatalf("expected COMPOSE_FILE to be restored, got %q", value)
	}
}

func TestPrepareReconcilesWhenProxyRoutesBecomeEmpty(t *testing.T) {
	state, server := newCaddyRouteState(t)
	workDir := t.TempDir()
	configPath := filepath.Join(workDir, "kool.yml")
	if err := os.WriteFile(configPath, []byte("proxy:\n  domain: app.localhost\n  routes: {}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	env := environment.NewFakeEnvStorage()
	env.Set("PWD", workDir)
	env.Set("COMPOSE_PROJECT_NAME", "example")
	env.Set("KOOL_PROXY_HOST", "app.localhost")
	manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)
	manager.adminURL = server.URL
	mustRegister(t, manager, route{Service: "app", Listen: 80, Target: 80, Hosts: []string{"@", "*"}})

	cleanup, err := manager.Prepare(nil)
	if err != nil {
		t.Fatal(err)
	}
	cleanup()
	state.requireRoutes(t, "kool-80", nil)
}

func TestRegisterWorkspaceRoute(t *testing.T) {
	serverCreated := false
	var events []string
	var registered map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch {
		case request.Method == http.MethodGet && strings.HasSuffix(request.URL.Path, "/routes"):
			_, _ = response.Write([]byte(`[]`))
		case request.Method == http.MethodGet:
			if !serverCreated {
				_, _ = response.Write([]byte("null"))
				return
			}
			_, _ = response.Write([]byte(`{"listen":[":80"],"routes":[]}`))
		case request.Method == http.MethodPost && strings.HasSuffix(request.URL.Path, "/kool-80"):
			serverCreated = true
			events = append(events, "server")
			return
		case request.Method == http.MethodDelete:
			response.WriteHeader(http.StatusNotFound)
			return
		case request.Method == http.MethodPatch && strings.HasSuffix(request.URL.Path, "/routes"):
			if !serverCreated {
				http.Error(response, "server was not created", http.StatusInternalServerError)
				return
			}
			events = append(events, "route")
			body, _ := io.ReadAll(request.Body)
			var routes []map[string]interface{}
			_ = json.Unmarshal(body, &routes)
			registered = routes[0]
			return
		}
		response.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	env := environment.NewFakeEnvStorage()
	env.Set("KOOL_WORKSPACE", "true")
	env.Set("KOOL_WORKSPACE_NAME", "task-a")
	env.Set("KOOL_WORKSPACE_PROJECT", "example-workspace-task-a")
	manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)
	manager.adminURL = server.URL

	env.Set("KOOL_PROXY_HOST", "task-a.workspace.app.localhost")
	if err := manager.register(route{Service: "app", Listen: 80, Target: 8080, Hosts: []string{"@", "*"}}); err != nil {
		t.Fatal(err)
	}
	if strings.Join(events, ",") != "server,route" {
		t.Errorf("expected server creation before route registration, got %v", events)
	}
	encoded, _ := json.Marshal(registered)
	result := string(encoded)
	for _, expected := range []string{"task-a.workspace.app.localhost", "*.task-a.workspace.app.localhost", "example-workspace-task-a-app:8080"} {
		if !strings.Contains(result, expected) {
			t.Errorf("expected registered route to contain %q, got %s", expected, result)
		}
	}
}

func TestCaddyHTTPServerDisablesAutomaticHTTPS(t *testing.T) {
	server := caddyServerConfig(route{Listen: 3001})
	automaticHTTPS, ok := server["automatic_https"].(map[string]interface{})
	if !ok || automaticHTTPS["disable"] != true {
		t.Errorf("expected HTTP server to disable automatic HTTPS, got %#v", server)
	}
	if _, exists := server["tls_connection_policies"]; exists {
		t.Errorf("did not expect HTTP server to configure TLS, got %#v", server)
	}
}

func TestCaddyHTTPSServerPreservesTLSConfiguration(t *testing.T) {
	server := caddyServerConfig(route{Listen: 3001, HTTPS: true})
	if _, exists := server["tls_connection_policies"]; !exists {
		t.Errorf("expected HTTPS server to configure TLS, got %#v", server)
	}
	if _, exists := server["automatic_https"]; exists {
		t.Errorf("did not expect HTTPS server to disable automatic HTTPS, got %#v", server)
	}
}

func TestRegisterRouteWithExistingServer(t *testing.T) {
	serverPostCalled := false
	routePostCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch {
		case request.Method == http.MethodGet && strings.HasSuffix(request.URL.Path, "/routes"):
			_, _ = response.Write([]byte(`[]`))
		case request.Method == http.MethodGet:
			_, _ = response.Write([]byte(`{"listen":[":80"],"routes":[]}`))
		case request.Method == http.MethodPost && strings.HasSuffix(request.URL.Path, "/kool-80"):
			serverPostCalled = true
		case request.Method == http.MethodPatch && strings.HasSuffix(request.URL.Path, "/routes"):
			routePostCalled = true
		}
	}))
	defer server.Close()

	env := environment.NewFakeEnvStorage()
	env.Set("KOOL_NAME", "example")
	env.Set("KOOL_PROXY_HOST", "app.localhost")
	manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)
	manager.adminURL = server.URL

	if err := manager.register(route{Service: "app", Listen: 80, Target: 8080, Hosts: []string{"@"}}); err != nil {
		t.Fatal(err)
	}
	if serverPostCalled {
		t.Error("did not expect an existing server to be recreated")
	}
	if !routePostCalled {
		t.Error("expected route to be registered on the existing server")
	}
}

func TestRegisterRouteReportsWhetherRouteWasCreated(t *testing.T) {
	state, server := newCaddyRouteState(t)
	manager := testRouteManager(server.URL, "example", "app.localhost")
	proxyRoute := route{Service: "app", Listen: 80, Target: 80, Hosts: []string{"@"}}

	created, err := manager.registerWithResult(proxyRoute)
	if err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Error("expected first registration to report a newly created route")
	}
	created, err = manager.registerWithResult(proxyRoute)
	if err != nil {
		t.Fatal(err)
	}
	if created {
		t.Error("expected repeated registration to report an existing route")
	}
	state.requireRoutes(t, "kool-80", []string{"kool-example-app-80-80"})
}

func TestRegisterRouteRejectsHostClaimedByAnotherProject(t *testing.T) {
	_, server := newCaddyRouteState(t)
	first := testRouteManager(server.URL, "first", "app.localhost")
	second := testRouteManager(server.URL, "second", "app.localhost")
	proxyRoute := route{Service: "app", Listen: 80, Target: 80, Hosts: []string{"@"}}
	mustRegister(t, first, proxyRoute)

	err := second.register(proxyRoute)
	if err == nil || !strings.Contains(err.Error(), "host conflict") {
		t.Fatalf("expected cross-project host conflict, got %v", err)
	}
}

func TestRegisterRouteChangesListenerMode(t *testing.T) {
	tests := []struct {
		name   string
		server string
		https  bool
	}{
		{name: "HTTPS on HTTP listener", server: `{"listen":[":80"],"automatic_https":{"disable":true},"routes":[]}`, https: true},
		{name: "HTTP on HTTPS listener", server: `{"listen":[":80"],"tls_connection_policies":[{}],"routes":[]}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			listenerPatched := false
			server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				switch {
				case request.Method == http.MethodGet && strings.HasSuffix(request.URL.Path, "/routes"):
					_, _ = response.Write([]byte(`[]`))
				case request.Method == http.MethodGet && strings.HasSuffix(request.URL.Path, "/kool-80"):
					_, _ = response.Write([]byte(test.server))
				case request.Method == http.MethodPatch && strings.HasSuffix(request.URL.Path, "/kool-80"):
					listenerPatched = true
				default:
					response.WriteHeader(http.StatusOK)
				}
			}))
			defer server.Close()

			env := environment.NewFakeEnvStorage()
			manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)
			manager.adminURL = server.URL
			if err := manager.register(route{Service: "app", Listen: 80, Target: 8080, HTTPS: test.https}); err != nil {
				t.Fatal(err)
			}
			if !listenerPatched {
				t.Error("expected existing listener protocol to be replaced")
			}
		})
	}
}

func TestRegisterRouteRejectsListenerModeChangeUsedByAnotherProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		_, _ = response.Write([]byte(`{"listen":[":80"],"automatic_https":{"disable":true},"routes":[{"@id":"kool-other-app-80-80","handle":[{"@id":"kool-project-other-kool-other-app-80-80"}]}]}`))
	}))
	defer server.Close()

	env := environment.NewFakeEnvStorage()
	env.Set("KOOL_NAME", "example")
	manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)
	manager.adminURL = server.URL
	err := manager.register(route{Service: "app", Listen: 80, Target: 8080, HTTPS: true})
	if err == nil || !strings.Contains(err.Error(), "another project uses it") {
		t.Fatalf("expected listener ownership conflict, got %v", err)
	}
}

func TestRegisterRouteCreatesServerAfterNotFound(t *testing.T) {
	serverCreated := false
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch {
		case request.Method == http.MethodGet && strings.HasSuffix(request.URL.Path, "/routes") && serverCreated:
			_, _ = response.Write([]byte(`[]`))
		case request.Method == http.MethodGet:
			response.WriteHeader(http.StatusNotFound)
		case request.Method == http.MethodPost && strings.HasSuffix(request.URL.Path, "/kool-80"):
			serverCreated = true
		}
	}))
	defer server.Close()

	env := environment.NewFakeEnvStorage()
	env.Set("KOOL_NAME", "example")
	env.Set("KOOL_PROXY_HOST", "app.localhost")
	manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)
	manager.adminURL = server.URL

	if err := manager.register(route{Service: "app", Listen: 80, Target: 8080, Hosts: []string{"@"}}); err != nil {
		t.Fatal(err)
	}
	if !serverCreated {
		t.Error("expected a missing server to be created after a 404")
	}
}

func TestRegisterTLS(t *testing.T) {
	var policy map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch {
		case request.Method == http.MethodGet:
			_, _ = response.Write([]byte(`{"automation":{"policies":[]}}`))
		case request.Method == http.MethodDelete:
			response.WriteHeader(http.StatusNotFound)
		case request.Method == http.MethodPost && strings.HasSuffix(request.URL.Path, "/policies"):
			body, _ := io.ReadAll(request.Body)
			_ = json.Unmarshal(body, &policy)
		default:
			response.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	env := environment.NewFakeEnvStorage()
	env.Set("KOOL_NAME", "example")
	env.Set("KOOL_PROXY_HOST", "app.localhost")
	manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)
	manager.adminURL = server.URL

	if err := manager.registerTLS([]route{
		{Service: "app", Listen: 443, Target: 80, Hosts: []string{"@", "*"}, HTTPS: true},
		{Service: "node", Listen: 3001, Target: 3001, Hosts: []string{"@"}, HTTPS: true},
	}); err != nil {
		t.Fatal(err)
	}
	encoded, _ := json.Marshal(policy)
	result := string(encoded)
	for _, expected := range []string{"kool-example-tls", "app.localhost", "*.app.localhost", "internal"} {
		if !strings.Contains(result, expected) {
			t.Errorf("expected TLS policy to contain %q, got %s", expected, result)
		}
	}
}

func TestRegisterTLSRemovesPolicyWhenHTTPSIsDisabled(t *testing.T) {
	deleted := ""
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodDelete {
			deleted = request.URL.Path
		}
		response.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	env := environment.NewFakeEnvStorage()
	env.Set("KOOL_NAME", "example")
	manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)
	manager.adminURL = server.URL
	if err := manager.registerTLS([]route{{Service: "app", Listen: 80, Target: 80}}); err != nil {
		t.Fatal(err)
	}
	if deleted != "/id/kool-example-tls" {
		t.Fatalf("expected stale project TLS policy to be deleted, got %q", deleted)
	}
}

func TestRestoreAppsReplacesPreviousConfiguration(t *testing.T) {
	var method string
	var restored []byte
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		method = request.Method
		restored, _ = io.ReadAll(request.Body)
		response.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	manager := NewManager(&shell.FakeShell{}, environment.NewFakeEnvStorage()).(*DefaultManager)
	manager.adminURL = server.URL
	snapshot := []byte(`{"http":{"servers":{"kool-80":{"listen":[":80"]}}}}`)
	if err := manager.restoreApps(snapshot, true); err != nil {
		t.Fatal(err)
	}
	if method != http.MethodPatch || string(restored) != string(snapshot) {
		t.Fatalf("expected previous apps configuration to be restored, got %s %s", method, restored)
	}
}

func TestRestoreAppsDeletesConfigurationWhenPreviouslyMissing(t *testing.T) {
	method := ""
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		method = request.Method
		response.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	manager := NewManager(&shell.FakeShell{}, environment.NewFakeEnvStorage()).(*DefaultManager)
	manager.adminURL = server.URL
	if err := manager.restoreApps(nil, false); err != nil {
		t.Fatal(err)
	}
	if method != http.MethodDelete {
		t.Fatalf("expected newly created apps configuration to be deleted, got %s", method)
	}
}

func TestWorkspaceRoutePrecedence(t *testing.T) {
	proxyRoute := route{Service: "app", Listen: 80, Target: 80, Hosts: []string{"@", "*"}}
	viteRoute := route{Service: "node", Listen: 3001, Target: 3001, Hosts: []string{"@", "*"}}

	t.Run("source started before workspace", func(t *testing.T) {
		state, server := newCaddyRouteState(t)
		source := testRouteManager(server.URL, "exlink", "exlink.localhost")
		workspace := testRouteManager(server.URL, "exlink-workspace-task-a", "task-a.workspace.exlink.localhost")

		mustRegister(t, source, proxyRoute)
		mustRegister(t, workspace, proxyRoute)

		state.requireRoutes(t, "kool-80", []string{
			"kool-exlink-workspace-task-a-app-80-80",
			"kool-exlink-app-80-80",
		})
		state.requireUpstream(t, "kool-80", "kool-exlink-workspace-task-a-app-80-80", "exlink-workspace-task-a-app:80")
	})

	t.Run("workspace started before and routes re-registered", func(t *testing.T) {
		state, server := newCaddyRouteState(t)
		source := testRouteManager(server.URL, "exlink", "exlink.localhost")
		workspace := testRouteManager(server.URL, "exlink-workspace-task-a", "task-a.workspace.exlink.localhost")

		mustRegister(t, workspace, proxyRoute)
		mustRegister(t, source, proxyRoute)
		mustRegister(t, source, proxyRoute)
		mustRegister(t, workspace, proxyRoute)

		state.requireRoutes(t, "kool-80", []string{
			"kool-exlink-workspace-task-a-app-80-80",
			"kool-exlink-app-80-80",
		})
	})

	t.Run("multiple workspaces and unrelated project retain order", func(t *testing.T) {
		state, server := newCaddyRouteState(t)
		source := testRouteManager(server.URL, "exlink", "exlink.localhost")
		other := testRouteManager(server.URL, "exlink-api", "other.localhost")
		workspaceA := testRouteManager(server.URL, "exlink-workspace-task-a", "task-a.workspace.exlink.localhost")
		workspaceB := testRouteManager(server.URL, "exlink-workspace-task-b", "task-b.workspace.exlink.localhost")

		mustRegister(t, source, proxyRoute)
		mustRegister(t, other, proxyRoute)
		mustRegister(t, workspaceA, proxyRoute)
		mustRegister(t, workspaceB, proxyRoute)
		mustRegister(t, workspaceA, proxyRoute)

		state.requireRoutes(t, "kool-80", []string{
			"kool-exlink-workspace-task-a-app-80-80",
			"kool-exlink-workspace-task-b-app-80-80",
			"kool-exlink-app-80-80",
			"kool-exlink-api-app-80-80",
		})
	})

	t.Run("app and Vite listeners use the same precedence", func(t *testing.T) {
		state, server := newCaddyRouteState(t)
		source := testRouteManager(server.URL, "exlink", "exlink.localhost")
		workspace := testRouteManager(server.URL, "exlink-workspace-task-a", "task-a.workspace.exlink.localhost")

		for _, manager := range []*DefaultManager{source, workspace} {
			mustRegister(t, manager, proxyRoute)
			mustRegister(t, manager, viteRoute)
		}

		state.requireRoutes(t, "kool-80", []string{
			"kool-exlink-workspace-task-a-app-80-80",
			"kool-exlink-app-80-80",
		})
		state.requireRoutes(t, "kool-3001", []string{
			"kool-exlink-workspace-task-a-node-3001-3001",
			"kool-exlink-node-3001-3001",
		})
	})

	t.Run("stopping one workspace removes only its routes", func(t *testing.T) {
		state, server := newCaddyRouteState(t)
		source := testRouteManager(server.URL, "exlink", "exlink.localhost")
		workspaceA := testRouteManager(server.URL, "exlink-workspace-task-a", "task-a.workspace.exlink.localhost")
		workspaceB := testRouteManager(server.URL, "exlink-workspace-task-b", "task-b.workspace.exlink.localhost")

		for _, manager := range []*DefaultManager{source, workspaceA, workspaceB} {
			mustRegister(t, manager, proxyRoute)
			mustRegister(t, manager, viteRoute)
		}
		if err := workspaceA.removeProjectRoutes(nil); err != nil {
			t.Fatal(err)
		}

		state.requireRoutes(t, "kool-80", []string{
			"kool-exlink-workspace-task-b-app-80-80",
			"kool-exlink-app-80-80",
		})
		state.requireRoutes(t, "kool-3001", []string{
			"kool-exlink-workspace-task-b-node-3001-3001",
			"kool-exlink-node-3001-3001",
		})
		if err := source.removeProjectRoutes(nil); err != nil {
			t.Fatal(err)
		}
		state.requireRoutes(t, "kool-80", []string{"kool-exlink-workspace-task-b-app-80-80"})
		state.requireRoutes(t, "kool-3001", []string{"kool-exlink-workspace-task-b-node-3001-3001"})
	})

	t.Run("changed configuration removes stale routes on every listener", func(t *testing.T) {
		state, server := newCaddyRouteState(t)
		source := testRouteManager(server.URL, "exlink", "exlink.localhost")
		oldApp := route{Service: "app", Listen: 80, Target: 80, Hosts: []string{"@", "*"}}
		oldVite := route{Service: "node", Listen: 3001, Target: 3001, Hosts: []string{"@", "*"}}
		newApp := route{Service: "app", Listen: 8080, Target: 8080, Hosts: []string{"@", "*"}}

		mustRegister(t, source, oldApp)
		mustRegister(t, source, oldVite)
		if err := source.reconcileRoutes([]route{newApp}); err != nil {
			t.Fatal(err)
		}

		state.requireRoutes(t, "kool-80", nil)
		state.requireRoutes(t, "kool-3001", nil)
	})
}

type caddyRouteState struct {
	routes map[string][]json.RawMessage
}

func newCaddyRouteState(t *testing.T) (*caddyRouteState, *httptest.Server) {
	t.Helper()
	state := &caddyRouteState{routes: make(map[string][]json.RawMessage)}
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		const serverPrefix = "/config/apps/http/servers/"
		if request.Method == http.MethodGet && request.URL.Path == strings.TrimSuffix(serverPrefix, "/") {
			servers := make(map[string]map[string][]json.RawMessage, len(state.routes))
			for serverID, routes := range state.routes {
				servers[serverID] = map[string][]json.RawMessage{"routes": routes}
			}
			_ = json.NewEncoder(response).Encode(servers)
			return
		}
		switch {
		case request.Method == http.MethodDelete && strings.HasPrefix(request.URL.Path, "/id/"):
			id := strings.TrimPrefix(request.URL.Path, "/id/")
			for serverID, routes := range state.routes {
				for index, route := range routes {
					if caddyRouteID(route) == id {
						state.routes[serverID] = append(routes[:index], routes[index+1:]...)
						response.WriteHeader(http.StatusOK)
						return
					}
				}
			}
			response.WriteHeader(http.StatusNotFound)
			return
		case !strings.HasPrefix(request.URL.Path, serverPrefix):
			response.WriteHeader(http.StatusOK)
			return
		}

		path := strings.TrimPrefix(request.URL.Path, serverPrefix)
		parts := strings.Split(path, "/")
		serverID := parts[0]
		isRoutes := len(parts) == 2 && parts[1] == "routes"
		switch {
		case request.Method == http.MethodGet && isRoutes:
			routes, exists := state.routes[serverID]
			if !exists {
				response.WriteHeader(http.StatusNotFound)
				return
			}
			_ = json.NewEncoder(response).Encode(routes)
		case request.Method == http.MethodGet:
			routes, exists := state.routes[serverID]
			if !exists {
				_, _ = response.Write([]byte("null"))
				return
			}
			_ = json.NewEncoder(response).Encode(map[string]interface{}{"routes": routes})
		case request.Method == http.MethodPost && len(parts) == 1:
			state.routes[serverID] = []json.RawMessage{}
		case request.Method == http.MethodPatch && isRoutes:
			var routes []json.RawMessage
			if err := json.NewDecoder(request.Body).Decode(&routes); err != nil {
				http.Error(response, err.Error(), http.StatusBadRequest)
				return
			}
			state.routes[serverID] = routes
		default:
			response.WriteHeader(http.StatusOK)
		}
	}))
	t.Cleanup(server.Close)
	return state, server
}

func testRouteManager(adminURL, project, host string) *DefaultManager {
	env := environment.NewFakeEnvStorage()
	env.Set("KOOL_PROXY_HOST", host)
	env.Set("COMPOSE_PROJECT_NAME", project)
	manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)
	manager.adminURL = adminURL
	return manager
}

func mustRegister(t *testing.T, manager *DefaultManager, route route) {
	t.Helper()
	if err := manager.register(route); err != nil {
		t.Fatal(err)
	}
}

func (state *caddyRouteState) requireRoutes(t *testing.T, serverID string, expected []string) {
	t.Helper()
	var actual []string
	for _, route := range state.routes[serverID] {
		actual = append(actual, caddyRouteID(route))
	}
	if strings.Join(actual, ",") != strings.Join(expected, ",") {
		t.Fatalf("expected %s routes %v, got %v", serverID, expected, actual)
	}
}

func (state *caddyRouteState) requireUpstream(t *testing.T, serverID, routeID, expected string) {
	t.Helper()
	for _, rawRoute := range state.routes[serverID] {
		if caddyRouteID(rawRoute) != routeID {
			continue
		}
		var routeConfig struct {
			Handle []struct {
				Upstreams []struct {
					Dial string `json:"dial"`
				} `json:"upstreams"`
			} `json:"handle"`
		}
		if json.Unmarshal(rawRoute, &routeConfig) != nil || len(routeConfig.Handle) == 0 || len(routeConfig.Handle[0].Upstreams) == 0 {
			t.Fatalf("route %s has no upstream", routeID)
		}
		if actual := routeConfig.Handle[0].Upstreams[0].Dial; actual != expected {
			t.Fatalf("expected route %s upstream %s, got %s", routeID, expected, actual)
		}
		return
	}
	t.Fatalf("route %s was not found", routeID)
}

func caddyRouteID(route json.RawMessage) string {
	var metadata caddyRouteMetadata
	_ = json.Unmarshal(route, &metadata)
	return metadata.ID
}
