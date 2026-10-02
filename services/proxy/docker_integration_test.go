package proxy

import (
	"fmt"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Run explicitly with KOOL_PROXY_INTEGRATION=1; never uses the running Kool proxy.
func TestDockerProxyAdminIsolationAndRestart(t *testing.T) {
	if os.Getenv("KOOL_PROXY_INTEGRATION") != "1" {
		t.Skip("set KOOL_PROXY_INTEGRATION=1 to test with disposable Docker resources")
	}
	docker := func(args ...string) string {
		t.Helper()
		output, err := exec.Command("docker", args...).CombinedOutput()
		if err != nil {
			t.Fatalf("docker %v: %v\n%s", args, err, output)
		}
		return strings.TrimSpace(string(output))
	}
	name := fmt.Sprintf("kool-proxy-test-%d-%d", os.Getpid(), time.Now().UnixNano())
	adminNetwork, appNetwork := name+"-admin", name+"-app"
	for _, network := range []string{adminNetwork, appNetwork} {
		docker("network", "create", network)
		t.Cleanup(func() { docker("network", "rm", network) })
	}
	t.Setenv("HOME", t.TempDir())
	env := environment.NewFakeEnvStorage()
	env.Set("COMPOSE_PROJECT_NAME", "integration")
	env.Set("KOOL_PROXY_HOST", "integration.localhost")
	manager := NewManager(&shell.FakeShell{}, env).(*DefaultManager)
	configPath, err := manager.ensureBaseConfig()
	if err != nil {
		t.Fatal(err)
	}
	docker("run", "--detach", "--name", name, "--network", adminNetwork,
		"--network-alias", caddyAdminHost, "--publish", "127.0.0.1::2019",
		"--volume", filepath.Dir(configPath)+":/etc/kool-proxy:ro",
		"--entrypoint", "/bin/sh", caddyImage, "-c", caddyStartCmd)
	t.Cleanup(func() { docker("rm", "--force", name) })
	docker("network", "connect", appNetwork, name)
	manager.adminURL = "http://" + docker("port", name, "2019/tcp")
	manager.http.Transport = proxyRoundTripFunc(func(request *http.Request) (*http.Response, error) {
		// The test uses an ephemeral published port; production always uses 2019.
		request.Host = "127.0.0.1:2019"
		return http.DefaultTransport.RoundTrip(request)
	})
	if err = manager.waitForCaddy(); err != nil {
		t.Fatalf("%v\n%s", err, docker("logs", name))
	}
	appIP := docker("inspect", "--format", fmt.Sprintf(`{{(index .NetworkSettings.Networks %q).IPAddress}}`, appNetwork), name)
	if output, err := exec.Command("docker", "exec", name, "wget", "-T", "1", "-qO-", "http://"+appIP+":2019/config/").CombinedOutput(); err == nil {
		t.Fatalf("Admin API must not listen on the application interface: %s", output)
	} else if !strings.Contains(string(output), "Connection refused") {
		t.Fatalf("expected refused application-interface connection, got %v: %s", err, output)
	}
	proxyRoute := route{Service: "app", Listen: 8080, Target: 8080, Hosts: []string{"@"}}
	if err = manager.register(proxyRoute); err != nil {
		t.Fatal(err)
	}
	apps, exists, err := manager.snapshotApps()
	if err != nil || !exists {
		t.Fatalf("snapshot: exists=%v err=%v", exists, err)
	}
	if err = writeProxyConfig(configPath, apps); err != nil {
		t.Fatal(err)
	}
	docker("restart", name)
	manager.adminURL = "http://" + docker("port", name, "2019/tcp")
	if err = manager.waitForCaddy(); err != nil {
		t.Fatalf("%v\n%s", err, docker("logs", name))
	}
	apps, _, err = manager.snapshotApps()
	if err != nil || !strings.Contains(string(apps), manager.routeID(proxyRoute)) {
		t.Fatalf("saved route must survive restart: %s (error %v)", apps, err)
	}
}
