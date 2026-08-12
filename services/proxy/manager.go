package proxy

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/parser"
	"kool-dev/kool/core/shell"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/compose-spec/compose-go/template"
	"golang.org/x/sys/unix"
)

const (
	caddyContainer = "kool-proxy"
	caddyImage     = "caddy:2.10-alpine"
	caddyVolume    = "kool_proxy"
	caddyAdminURL  = "http://127.0.0.1:2019"
	caddyAdminHost = "kool-proxy-admin"
	caddyAdminNet  = "kool_proxy_admin"
	caddyStartCmd  = "exec caddy run --config /etc/caddy/caddy.json"
	defaultNetwork = "kool_global"
)

// Manager controls Kool's local proxy routes.
type Manager interface {
	Prepare([]string) (func(bool), error)
	Remove([]string) error
	RemoveProject(string, []string) error
	Trust() error
}

type route struct {
	Service string
	Listen  int
	Target  int
	Hosts   []string
	HTTPS   bool
}

type config struct {
	Domain  string
	HTTPS   bool
	Network string
	Routes  []route
}

// DefaultManager manages a global Caddy container through its Admin API.
type DefaultManager struct {
	shell           shell.Shell
	env             environment.EnvStorage
	http            *http.Client
	adminURL        string
	projectOverride string
}

// NewManager creates a proxy manager for the current project.
func NewManager(sh shell.Shell, env environment.EnvStorage) Manager {
	return &DefaultManager{shell: sh, env: env, http: &http.Client{Timeout: 2 * time.Second}, adminURL: caddyAdminURL}
}

// Prepare ensures Caddy is running, applies network aliases, and registers routes.
func (m *DefaultManager) Prepare(services []string) (finish func(bool), err error) {
	finish = func(bool) {}
	var cfg *config
	if cfg, err = m.loadConfig(); err != nil || cfg == nil {
		return
	}

	requested := make(map[string]bool, len(services))
	for _, service := range services {
		requested[service] = true
	}
	var routes []route
	for _, route := range cfg.Routes {
		if len(requested) == 0 || requested[route.Service] {
			routes = append(routes, route)
		}
	}
	previousComposeFiles := m.env.Get("COMPOSE_FILE")
	var override string
	if len(routes) > 0 {
		if override, err = m.createAliasOverride(cfg.Network, routes); err != nil {
			return
		}
	}
	cleanupOverride := func() {
		m.env.Set("COMPOSE_FILE", previousComposeFiles)
		if override != "" {
			_ = os.Remove(override)
		}
	}

	var unlock func()
	if unlock, err = m.acquireConfigLock(); err != nil {
		cleanupOverride()
		return func(bool) {}, err
	}
	lockHeld := true
	defer func() {
		if lockHeld {
			unlock()
		}
	}()

	if len(routes) == 0 {
		if _, inspectErr := m.shell.Exec(builder.NewCommand("docker", "inspect", caddyContainer)); inspectErr != nil {
			unlock()
			lockHeld = false
			cleanupOverride()
			return func(bool) {}, nil
		}
	} else if err = m.ensureCaddy(cfg.Network, cfg.Routes); err != nil {
		cleanupOverride()
		return func(bool) {}, err
	}

	var snapshot []byte
	var snapshotExists bool
	if snapshot, snapshotExists, err = m.snapshotApps(); err != nil {
		cleanupOverride()
		return func(bool) {}, err
	}
	rollback := func() { _ = m.restoreApps(snapshot, snapshotExists) }
	if err = m.registerTLSUnlocked(cfg.Routes); err != nil {
		rollback()
		cleanupOverride()
		return func(bool) {}, err
	}
	for _, route := range routes {
		if _, err = m.registerUnlocked(route); err != nil {
			rollback()
			cleanupOverride()
			return func(bool) {}, err
		}
	}
	if err = m.reconcileRoutesUnlocked(cfg.Routes); err != nil {
		rollback()
		cleanupOverride()
		return func(bool) {}, err
	}
	committed, committedExists, snapshotErr := m.snapshotApps()
	if snapshotErr != nil {
		rollback()
		cleanupOverride()
		return func(bool) {}, snapshotErr
	}
	unlock()
	lockHeld = false
	finished := false
	finish = func(success bool) {
		if finished {
			return
		}
		finished = true
		cleanupOverride()
		if success {
			return
		}
		_ = m.withConfigLock(func() error {
			return m.restoreAppsIfUnchanged(committed, committedExists, snapshot, snapshotExists)
		})
	}
	return
}

// Remove deletes all proxy routes belonging to the current project.
func (m *DefaultManager) Remove(services []string) error {
	_, err := m.loadConfig()
	if err != nil {
		return err
	}
	if _, err = m.shell.Exec(builder.NewCommand("docker", "inspect", caddyContainer)); err != nil {
		return nil
	}
	return m.withConfigLock(func() error {
		if err = m.removeProjectRoutesUnlocked(services); err != nil {
			return err
		}
		if len(services) == 0 {
			return m.deleteRoute(m.tlsID())
		}
		return nil
	})
}

// RemoveProject removes proxy routes for an explicit Compose project.
func (m *DefaultManager) RemoveProject(project string, services []string) error {
	manager := *m
	manager.projectOverride = project
	if _, err := manager.shell.Exec(builder.NewCommand("docker", "inspect", caddyContainer)); err != nil {
		return nil
	}
	return manager.withConfigLock(func() error {
		if err := manager.removeProjectRoutesUnlocked(services); err != nil {
			return err
		}
		if len(services) == 0 {
			return manager.deleteRoute(manager.tlsID())
		}
		return nil
	})
}

// Trust installs Caddy's local root CA in the host trust store.
func (m *DefaultManager) Trust() error {
	if _, err := m.shell.Exec(builder.NewCommand("docker", "inspect", caddyContainer)); err != nil {
		return errors.New("kool proxy is not running; run kool start first")
	}
	certificate, err := os.CreateTemp("", "kool-proxy-root-*.crt")
	if err != nil {
		return err
	}
	path := certificate.Name()
	if err = certificate.Close(); err != nil {
		return err
	}
	defer func() { _ = os.Remove(path) }()
	if err = m.shell.Interactive(builder.NewCommand("docker", "cp"), caddyContainer+":/var/lib/caddy/data/caddy/pki/authorities/local/root.crt", path); err != nil {
		return fmt.Errorf("could not export proxy root certificate; start an HTTPS proxy route first: %w", err)
	}

	switch runtime.GOOS {
	case "darwin":
		return m.shell.Interactive(builder.NewCommand("sudo", "security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain"), path)
	case "linux":
		destination := "/usr/local/share/ca-certificates/kool-proxy.crt"
		if err = m.shell.Interactive(builder.NewCommand("sudo", "cp"), path, destination); err != nil {
			return err
		}
		return m.shell.Interactive(builder.NewCommand("sudo", "update-ca-certificates"))
	default:
		return fmt.Errorf("automatic proxy trust is not supported on %s", runtime.GOOS)
	}
}

func (m *DefaultManager) loadConfig() (*config, error) {
	workDir := m.env.Get("PWD")
	var file string
	for _, name := range []string{"kool.yml", "kool.yaml"} {
		candidate := filepath.Join(workDir, name)
		if _, err := os.Stat(candidate); err == nil {
			file = candidate
			break
		}
	}
	if file == "" {
		return nil, nil
	}

	parsed, err := parser.ParseKoolYaml(file)
	if err != nil || parsed.Proxy == nil {
		return nil, err
	}
	domain, err := template.Substitute(parsed.Proxy.Domain, os.LookupEnv)
	domain = strings.TrimSuffix(strings.TrimSpace(domain), ".")
	if err != nil || domain == "" {
		if err == nil {
			err = errors.New("proxy.domain cannot be empty")
		}
		return nil, err
	}
	if strings.Contains(domain, "://") {
		return nil, errors.New("proxy.domain must be a hostname without a URL scheme")
	}
	network := parsed.Proxy.Network
	if network == "" {
		network = m.env.Get("KOOL_GLOBAL_NETWORK")
	}
	if network == "" {
		network = defaultNetwork
	}
	if network, err = template.Substitute(network, os.LookupEnv); err != nil {
		return nil, err
	}
	if network == caddyAdminNet {
		return nil, fmt.Errorf("proxy.network %q is reserved for proxy administration", network)
	}

	cfg := &config{Domain: domain, HTTPS: parsed.Proxy.HTTPS, Network: network}
	for service, routeConfig := range parsed.Proxy.Routes {
		if len(routeConfig.Ports) == 0 {
			return nil, fmt.Errorf("proxy route %s requires at least one port", service)
		}
		if len(routeConfig.Hosts) == 0 {
			routeConfig.Hosts = []string{"@", "*"}
		}
		var hosts []string
		for _, host := range routeConfig.Hosts {
			var resolvedHost string
			if resolvedHost, err = template.Substitute(host, os.LookupEnv); err != nil {
				return nil, err
			}
			resolvedHost = strings.TrimSpace(resolvedHost)
			if resolvedHost == "" {
				return nil, fmt.Errorf("proxy route %s contains an empty host", service)
			}
			hosts = append(hosts, resolvedHost)
		}
		listenPorts := make(map[int]bool)
		for _, portMapping := range routeConfig.Ports {
			var resolved string
			if resolved, err = template.Substitute(portMapping, os.LookupEnv); err != nil {
				return nil, err
			}
			parts := strings.SplitN(resolved, ":", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("proxy route %s must use listen:target", service)
			}
			listen, listenErr := strconv.Atoi(strings.TrimSpace(parts[0]))
			target, targetErr := strconv.Atoi(strings.TrimSpace(parts[1]))
			if listenErr != nil || targetErr != nil || listen < 1 || listen > 65535 || target < 1 || target > 65535 {
				return nil, fmt.Errorf("proxy route %s has invalid ports %q", service, resolved)
			}
			if listen == 2019 {
				return nil, fmt.Errorf("proxy route %s cannot listen on reserved admin port 2019", service)
			}
			if listenPorts[listen] {
				return nil, fmt.Errorf("proxy route %s maps listen port %d more than once", service, listen)
			}
			listenPorts[listen] = true
			cfg.Routes = append(cfg.Routes, route{Service: service, Listen: listen, Target: target, Hosts: hosts, HTTPS: cfg.HTTPS})
		}
	}
	sort.Slice(cfg.Routes, func(i, j int) bool {
		if cfg.Routes[i].Service == cfg.Routes[j].Service {
			if cfg.Routes[i].Listen == cfg.Routes[j].Listen {
				return cfg.Routes[i].Target < cfg.Routes[j].Target
			}
			return cfg.Routes[i].Listen < cfg.Routes[j].Listen
		}
		return cfg.Routes[i].Service < cfg.Routes[j].Service
	})
	bindings := make(map[string]string)
	for _, route := range cfg.Routes {
		for _, host := range m.routeHosts(route) {
			host = strings.ToLower(strings.TrimSuffix(host, "."))
			binding := fmt.Sprintf("%d:%s", route.Listen, host)
			if service := bindings[binding]; service != "" && service != route.Service {
				return nil, fmt.Errorf("proxy routes %s and %s both use %s", service, route.Service, binding)
			}
			bindings[binding] = route.Service
		}
	}
	return cfg, nil
}

func (m *DefaultManager) createAliasOverride(network string, routes []route) (string, error) {
	file, err := os.CreateTemp("", "kool-proxy-*.yml")
	if err != nil {
		return "", err
	}

	var content strings.Builder
	content.WriteString("services:\n")
	written := make(map[string]bool)
	for _, route := range routes {
		if written[route.Service] {
			continue
		}
		fmt.Fprintf(&content, "  %q:\n    ports: !reset []\n    networks:\n      %q:\n        aliases:\n          - %q\n", route.Service, defaultNetwork, m.alias(route.Service))
		written[route.Service] = true
	}
	fmt.Fprintf(&content, "networks:\n  %q:\n    name: %q\n    external: true\n", defaultNetwork, network)
	if _, err = file.WriteString(content.String()); err == nil {
		err = file.Close()
	} else {
		_ = file.Close()
	}
	if err != nil {
		_ = os.Remove(file.Name())
		return "", err
	}

	separator := m.env.Get("COMPOSE_PATH_SEPARATOR")
	if separator == "" {
		separator = string(os.PathListSeparator)
	}
	composeFiles := m.env.Get("COMPOSE_FILE")
	if composeFiles == "" {
		for _, name := range []string{"compose.yaml", "compose.yml", "docker-compose.yml", "docker-compose.yaml"} {
			if _, statErr := os.Stat(filepath.Join(m.env.Get("PWD"), name)); statErr == nil {
				composeFiles = name
				break
			}
		}
	}
	m.env.Set("COMPOSE_FILE", composeFiles+separator+file.Name())
	return file.Name(), nil
}

func (m *DefaultManager) ensureCaddy(network string, routes []route) error {
	inspect := builder.NewCommand("docker", "inspect", "--format", "{{.State.Running}}", caddyContainer)
	running, err := m.shell.Exec(inspect)
	if err == nil {
		compatibility, inspectErr := m.shell.Exec(builder.NewCommand("docker", "inspect", "--format", "{{json .NetworkSettings.Networks}}|{{json .Config.Entrypoint}}|{{json .Config.Cmd}}", caddyContainer))
		if inspectErr != nil {
			return inspectErr
		}
		expectedCommand := `["/bin/sh"]|["-c","` + strings.ReplaceAll(caddyStartCmd, `"`, `\"`) + `"]`
		if !strings.Contains(compatibility, `"`+caddyAdminNet+`"`) || !strings.HasSuffix(compatibility, expectedCommand) {
			if err = m.shell.Interactive(builder.NewCommand("docker", "rm", "--force"), caddyContainer); err != nil {
				return err
			}
			err = errors.New("legacy proxy container removed")
		}
	}
	if err != nil {
		configPath, configErr := m.ensureBaseConfig()
		if configErr != nil {
			return configErr
		}
		if _, networkErr := m.shell.Exec(builder.NewCommand("docker", "network", "inspect", caddyAdminNet)); networkErr != nil {
			if networkErr = m.shell.Interactive(builder.NewCommand("docker", "network", "create"), caddyAdminNet); networkErr != nil {
				return networkErr
			}
		}
		args := []string{"run", "-d", "--name", caddyContainer, "--restart", "unless-stopped", "--network", caddyAdminNet, "--network-alias", caddyAdminHost, "-p", "127.0.0.1:2019:2019"}
		ports := make(map[int]bool)
		for _, route := range routes {
			ports[route.Listen] = true
		}
		var sortedPorts []int
		for port := range ports {
			sortedPorts = append(sortedPorts, port)
		}
		sort.Ints(sortedPorts)
		for _, port := range sortedPorts {
			mapping := fmt.Sprintf("%d:%d", port, port)
			args = append(args, "-p", mapping)
		}
		args = append(args,
			"-v", configPath+":/etc/caddy/caddy.json:ro",
			"-v", caddyVolume+":/var/lib/caddy",
			"-e", "XDG_CONFIG_HOME=/var/lib/caddy/config",
			"-e", "XDG_DATA_HOME=/var/lib/caddy/data",
			"--entrypoint", "/bin/sh",
			caddyImage,
			"-c", caddyStartCmd,
		)
		if err = m.shell.Interactive(builder.NewCommand("docker"), args...); err != nil {
			return err
		}
	} else if running != "true" {
		if err = m.shell.Interactive(builder.NewCommand("docker", "start"), caddyContainer); err != nil {
			return err
		}
	}
	networks, err := m.shell.Exec(builder.NewCommand("docker", "inspect", "--format", "{{json .NetworkSettings.Networks}}", caddyContainer))
	if err != nil {
		return err
	}
	if !strings.Contains(networks, `"`+network+`"`) {
		if err = m.shell.Interactive(builder.NewCommand("docker", "network", "connect", network), caddyContainer); err != nil {
			return err
		}
	}

	for _, route := range routes {
		if _, err = m.shell.Exec(builder.NewCommand("docker", "port", caddyContainer), fmt.Sprintf("%d/tcp", route.Listen)); err != nil {
			return fmt.Errorf("kool proxy does not publish port %d; remove %s and retry", route.Listen, caddyContainer)
		}
	}
	for attempt := 0; attempt < 30; attempt++ {
		response, requestErr := m.request(http.MethodGet, m.adminURL+"/config/", nil)
		if requestErr == nil {
			_ = response.Body.Close()
			if response.StatusCode < 500 {
				return nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errors.New("kool proxy did not become ready")
}

func (m *DefaultManager) ensureBaseConfig() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	directory := filepath.Join(home, ".kool", "proxy")
	if err = os.MkdirAll(directory, 0755); err != nil {
		return "", err
	}
	path := filepath.Join(directory, "caddy.json")
	content := []byte(`{"admin":{"listen":"` + caddyAdminHost + `:2019"},"apps":{"http":{"servers":{}},"tls":{"automation":{"policies":[]}}}}`)
	if err = os.WriteFile(path, content, 0644); err != nil {
		return "", err
	}
	return path, nil
}

func (m *DefaultManager) register(route route) error {
	_, err := m.registerWithResult(route)
	return err
}

func (m *DefaultManager) registerWithResult(route route) (created bool, err error) {
	err = m.withConfigLock(func() error {
		created, err = m.registerUnlocked(route)
		return err
	})
	return
}

func (m *DefaultManager) registerUnlocked(route route) (bool, error) {
	serverID := "kool-" + strconv.Itoa(route.Listen)
	server := caddyServerConfig(route)
	serverBody, _ := json.Marshal(server)
	serverURL := m.adminURL + "/config/apps/http/servers/" + serverID
	response, err := m.request(http.MethodGet, serverURL, nil)
	if err != nil {
		return false, err
	}
	status := response.StatusCode
	serverMissing := status == http.StatusNotFound
	if status >= 400 && !serverMissing {
		defer func() { _ = response.Body.Close() }()
		return false, responseError(response)
	}
	if serverMissing {
		if err = response.Body.Close(); err != nil {
			return false, err
		}
	} else {
		var responseBody []byte
		if responseBody, err = io.ReadAll(response.Body); err != nil {
			_ = response.Body.Close()
			return false, err
		}
		if err = response.Body.Close(); err != nil {
			return false, err
		}
		serverMissing = bytes.Equal(bytes.TrimSpace(responseBody), []byte("null"))
		if !serverMissing {
			if err = m.reconfigureListenerMode(serverURL, responseBody, route); err != nil {
				return false, err
			}
		}
	}
	if serverMissing {
		response, err = m.request(http.MethodPost, serverURL, serverBody)
		if err != nil {
			return false, err
		}
		defer func() { _ = response.Body.Close() }()
		if response.StatusCode >= 400 {
			return false, responseError(response)
		}
	}
	if route.HTTPS {
		if err = m.setTLSPolicy(serverURL); err != nil {
			return false, err
		}
	} else if err = m.disableAutomaticHTTPS(serverURL); err != nil {
		return false, err
	}

	routeConfig := map[string]interface{}{
		"@id":   m.routeID(route),
		"match": []interface{}{map[string]interface{}{"host": m.routeHosts(route)}},
		"handle": []interface{}{map[string]interface{}{
			"@id":       m.projectMarker() + "-" + m.routeID(route),
			"handler":   "reverse_proxy",
			"upstreams": []interface{}{map[string]string{"dial": m.alias(route.Service) + ":" + strconv.Itoa(route.Target)}},
		}},
		"terminal": true,
	}
	return m.upsertRouteUnlocked(serverURL, routeConfig)
}

func (m *DefaultManager) upsertRouteUnlocked(serverURL string, routeConfig map[string]interface{}) (bool, error) {
	routesURL := serverURL + "/routes"
	response, err := m.request(http.MethodGet, routesURL, nil)
	if err != nil {
		return false, err
	}
	if response.StatusCode >= 400 && response.StatusCode != http.StatusNotFound {
		defer func() { _ = response.Body.Close() }()
		return false, responseError(response)
	}

	var routes []json.RawMessage
	if response.StatusCode != http.StatusNotFound {
		if err = json.NewDecoder(response.Body).Decode(&routes); err != nil && err != io.EOF {
			_ = response.Body.Close()
			return false, err
		}
	}
	if err = response.Body.Close(); err != nil {
		return false, err
	}

	routeBody, _ := json.Marshal(routeConfig)
	routeID, _ := routeConfig["@id"].(string)
	replaced := false
	for index, existing := range routes {
		var metadata caddyRouteMetadata
		if json.Unmarshal(existing, &metadata) == nil && metadata.ID == routeID {
			routes[index] = routeBody
			replaced = true
			break
		}
	}
	if !replaced {
		var incoming caddyRouteMetadata
		_ = json.Unmarshal(routeBody, &incoming)
		for _, existing := range routes {
			var metadata caddyRouteMetadata
			if json.Unmarshal(existing, &metadata) != nil || m.routeBelongsToProject(metadata) {
				continue
			}
			if hostSetsOverlap(routeMetadataHosts(metadata), routeMetadataHosts(incoming)) {
				return false, fmt.Errorf("proxy listener host conflict with route %s", metadata.ID)
			}
		}
		routes = append(routes, routeBody)
	}
	routes = orderCaddyRoutes(routes)
	body, _ := json.Marshal(routes)
	response, err = m.request(http.MethodPatch, routesURL, body)
	if err != nil {
		return false, err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode >= 400 {
		return false, responseError(response)
	}
	return !replaced, nil
}

func hostSetsOverlap(first, second []string) bool {
	normalized := make(map[string]bool, len(first))
	for _, host := range first {
		normalized[strings.ToLower(strings.TrimSuffix(host, "."))] = true
	}
	for _, host := range second {
		host = strings.ToLower(strings.TrimSuffix(host, "."))
		if normalized[host] {
			return true
		}
	}
	return false
}

type caddyRouteMetadata struct {
	ID    string `json:"@id"`
	Match []struct {
		Hosts []string `json:"host"`
	} `json:"match"`
	Handle []struct {
		ID        string `json:"@id"`
		Upstreams []struct {
			Dial string `json:"dial"`
		} `json:"upstreams"`
	} `json:"handle"`
}

func (m *DefaultManager) reconcileRoutes(desiredRoutes []route) error {
	return m.withConfigLock(func() error {
		return m.reconcileRoutesUnlocked(desiredRoutes)
	})
}

func (m *DefaultManager) reconcileRoutesUnlocked(desiredRoutes []route) error {
	desired := make(map[string]bool, len(desiredRoutes))
	for _, route := range desiredRoutes {
		desired[m.routeID(route)] = true
	}
	return m.filterProjectRoutesUnlocked(func(route caddyRouteMetadata) bool {
		return !desired[route.ID]
	})
}

func (m *DefaultManager) removeProjectRoutes(services []string) error {
	return m.withConfigLock(func() error {
		return m.removeProjectRoutesUnlocked(services)
	})
}

func (m *DefaultManager) removeProjectRoutesUnlocked(services []string) error {
	aliases := make(map[string]bool, len(services))
	for _, service := range services {
		aliases[m.alias(service)] = true
	}
	return m.filterProjectRoutesUnlocked(func(route caddyRouteMetadata) bool {
		if len(aliases) == 0 {
			return true
		}
		for _, alias := range routeMetadataAliases(route) {
			if aliases[alias] {
				return true
			}
		}
		return false
	})
}

func (m *DefaultManager) filterProjectRoutes(shouldRemove func(caddyRouteMetadata) bool) error {
	return m.withConfigLock(func() error {
		return m.filterProjectRoutesUnlocked(shouldRemove)
	})
}

func (m *DefaultManager) filterProjectRoutesUnlocked(shouldRemove func(caddyRouteMetadata) bool) error {
	serversURL := m.adminURL + "/config/apps/http/servers"
	response, err := m.request(http.MethodGet, serversURL, nil)
	if err != nil {
		return err
	}
	if response.StatusCode == http.StatusNotFound {
		_ = response.Body.Close()
		return nil
	}
	if response.StatusCode >= 400 {
		defer func() { _ = response.Body.Close() }()
		return responseError(response)
	}
	var servers map[string]struct {
		Routes []json.RawMessage `json:"routes"`
	}
	if err = json.NewDecoder(response.Body).Decode(&servers); err != nil {
		_ = response.Body.Close()
		return err
	}
	if err = response.Body.Close(); err != nil {
		return err
	}

	for serverID, server := range servers {
		routes := server.Routes[:0]
		changed := false
		for _, rawRoute := range server.Routes {
			var metadata caddyRouteMetadata
			if json.Unmarshal(rawRoute, &metadata) == nil && m.routeBelongsToProject(metadata) && shouldRemove(metadata) {
				changed = true
				continue
			}
			routes = append(routes, rawRoute)
		}
		if !changed {
			continue
		}
		body, _ := json.Marshal(orderCaddyRoutes(routes))
		routesURL := serversURL + "/" + serverID + "/routes"
		response, err = m.request(http.MethodPatch, routesURL, body)
		if err != nil {
			return err
		}
		if response.StatusCode >= 400 {
			defer func() { _ = response.Body.Close() }()
			return responseError(response)
		}
		if err = response.Body.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (m *DefaultManager) withConfigLock(action func() error) error {
	unlock, err := m.acquireConfigLock()
	if err != nil {
		return err
	}
	defer unlock()
	return action()
}

func (m *DefaultManager) acquireConfigLock() (func(), error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	directory := filepath.Join(home, ".kool", "proxy")
	if err = os.MkdirAll(directory, 0755); err != nil {
		return nil, err
	}
	lock, err := os.OpenFile(filepath.Join(directory, "config.lock"), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err = unix.Flock(int(lock.Fd()), unix.LOCK_EX); err != nil {
		_ = lock.Close()
		return nil, err
	}
	return func() {
		_ = unix.Flock(int(lock.Fd()), unix.LOCK_UN)
		_ = lock.Close()
	}, nil
}

func (m *DefaultManager) routeBelongsToProject(route caddyRouteMetadata) bool {
	prefix := m.projectMarker() + "-"
	for _, handle := range route.Handle {
		if strings.HasPrefix(handle.ID, prefix) {
			return true
		}
	}
	return false
}

func (m *DefaultManager) projectMarker() string {
	digest := sha256.Sum256([]byte(m.project()))
	return fmt.Sprintf("kool-project-%x", digest[:8])
}

func routeMetadataAliases(route caddyRouteMetadata) []string {
	var aliases []string
	for _, handle := range route.Handle {
		for _, upstream := range handle.Upstreams {
			alias := upstream.Dial
			if separator := strings.LastIndex(alias, ":"); separator >= 0 {
				alias = alias[:separator]
			}
			aliases = append(aliases, alias)
		}
	}
	return aliases
}

func orderCaddyRoutes(routes []json.RawMessage) []json.RawMessage {
	ordered := append([]json.RawMessage(nil), routes...)
	for broad := 0; broad < len(ordered); broad++ {
		var broadMetadata caddyRouteMetadata
		if json.Unmarshal(ordered[broad], &broadMetadata) != nil || !strings.HasPrefix(broadMetadata.ID, "kool-") {
			continue
		}
		for specific := broad + 1; specific < len(ordered); specific++ {
			var specificMetadata caddyRouteMetadata
			if json.Unmarshal(ordered[specific], &specificMetadata) != nil || !strings.HasPrefix(specificMetadata.ID, "kool-") {
				continue
			}
			if routeHostsAreMoreSpecific(routeMetadataHosts(specificMetadata), routeMetadataHosts(broadMetadata)) {
				route := ordered[specific]
				copy(ordered[broad+1:specific+1], ordered[broad:specific])
				ordered[broad] = route
				broad--
				break
			}
		}
	}
	return ordered
}

func routeMetadataHosts(route caddyRouteMetadata) []string {
	var hosts []string
	for _, match := range route.Match {
		hosts = append(hosts, match.Hosts...)
	}
	return hosts
}

func routeHostsAreMoreSpecific(hosts, other []string) bool {
	return len(hosts) > 0 && hostPatternsCover(other, hosts) && !hostPatternsCover(hosts, other)
}

func hostPatternsCover(patterns, hosts []string) bool {
	for _, host := range hosts {
		covered := false
		for _, pattern := range patterns {
			if hostPatternCovers(pattern, host) {
				covered = true
				break
			}
		}
		if !covered {
			return false
		}
	}
	return true
}

func hostPatternCovers(pattern, host string) bool {
	pattern = strings.ToLower(strings.TrimSuffix(pattern, "."))
	host = strings.ToLower(strings.TrimSuffix(host, "."))
	if !strings.HasPrefix(pattern, "*.") {
		return pattern == host
	}
	pattern = strings.TrimPrefix(pattern, "*.")
	host = strings.TrimPrefix(host, "*.")
	return host == pattern || strings.HasSuffix(host, "."+pattern)
}

func caddyServerConfig(route route) map[string]interface{} {
	server := map[string]interface{}{"listen": []string{":" + strconv.Itoa(route.Listen)}, "routes": []interface{}{}}
	if route.HTTPS {
		server["tls_connection_policies"] = []interface{}{map[string]interface{}{}}
	} else {
		server["automatic_https"] = map[string]interface{}{"disable": true}
	}
	return server
}

func (m *DefaultManager) reconfigureListenerMode(serverURL string, server []byte, route route) error {
	var config struct {
		AutomaticHTTPS struct {
			Disable bool `json:"disable"`
		} `json:"automatic_https"`
		TLSConnectionPolicies json.RawMessage   `json:"tls_connection_policies"`
		Routes                []json.RawMessage `json:"routes"`
	}
	if err := json.Unmarshal(server, &config); err != nil {
		return fmt.Errorf("could not inspect proxy listener %d: %w", route.Listen, err)
	}
	hasTLS := len(config.TLSConnectionPolicies) > 0 && string(config.TLSConnectionPolicies) != "null"
	modeDiffers := route.HTTPS && config.AutomaticHTTPS.Disable || !route.HTTPS && hasTLS
	if !modeDiffers {
		return nil
	}
	for _, rawRoute := range config.Routes {
		var metadata caddyRouteMetadata
		if json.Unmarshal(rawRoute, &metadata) != nil || !m.routeBelongsToProject(metadata) {
			return fmt.Errorf("proxy listener %d cannot change protocol while another project uses it", route.Listen)
		}
	}

	var replacement map[string]interface{}
	if err := json.Unmarshal(server, &replacement); err != nil {
		return fmt.Errorf("could not inspect proxy listener %d: %w", route.Listen, err)
	}
	delete(replacement, "automatic_https")
	delete(replacement, "tls_connection_policies")
	for key, value := range caddyServerConfig(route) {
		if key != "routes" {
			replacement[key] = value
		}
	}
	body, _ := json.Marshal(replacement)
	response, err := m.request(http.MethodPatch, serverURL, body)
	if err != nil {
		return err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode >= 400 {
		return responseError(response)
	}
	return nil
}

func (m *DefaultManager) disableAutomaticHTTPS(serverURL string) error {
	body, _ := json.Marshal(map[string]interface{}{"disable": true})
	response, err := m.request(http.MethodPost, serverURL+"/automatic_https", body)
	if err != nil {
		return err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode >= 400 {
		return responseError(response)
	}
	return nil
}

func (m *DefaultManager) setTLSPolicy(serverURL string) error {
	url := serverURL + "/tls_connection_policies"
	response, err := m.request(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	method := http.MethodPatch
	if response.StatusCode == http.StatusNotFound {
		method = http.MethodPost
	} else if response.StatusCode >= 400 {
		defer func() { _ = response.Body.Close() }()
		return responseError(response)
	}
	_ = response.Body.Close()

	policies, _ := json.Marshal([]interface{}{map[string]interface{}{}})
	response, err = m.request(method, url, policies)
	if err != nil {
		return err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode >= 400 {
		return responseError(response)
	}
	return nil
}

func (m *DefaultManager) registerTLS(routes []route) error {
	return m.withConfigLock(func() error {
		return m.registerTLSUnlocked(routes)
	})
}

func (m *DefaultManager) registerTLSUnlocked(routes []route) error {
	seen := make(map[string]bool)
	var subjects []string
	for _, route := range routes {
		if !route.HTTPS {
			continue
		}
		for _, host := range m.routeHosts(route) {
			if !seen[host] {
				subjects = append(subjects, host)
				seen[host] = true
			}
		}
	}
	if len(subjects) == 0 {
		return m.deleteRoute(m.tlsID())
	}
	sort.Strings(subjects)

	tlsURL := m.adminURL + "/config/apps/tls"
	response, err := m.request(http.MethodGet, tlsURL, nil)
	if err != nil {
		return err
	}
	status := response.StatusCode
	_ = response.Body.Close()
	if status == http.StatusNotFound {
		body := []byte(`{"automation":{"policies":[]}}`)
		if response, err = m.request(http.MethodPost, tlsURL, body); err != nil {
			return err
		}
		defer func() { _ = response.Body.Close() }()
		if response.StatusCode >= 400 {
			return responseError(response)
		}
	} else if status >= 400 {
		return fmt.Errorf("caddy Admin API returned %s", response.Status)
	}

	if err = m.deleteRoute(m.tlsID()); err != nil {
		return err
	}
	policy := map[string]interface{}{
		"@id":      m.tlsID(),
		"subjects": subjects,
		"issuers":  []interface{}{map[string]string{"module": "internal"}},
	}
	body, _ := json.Marshal(policy)
	response, err = m.request(http.MethodPost, tlsURL+"/automation/policies", body)
	if err != nil {
		return err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode >= 400 {
		return responseError(response)
	}
	return nil
}

func (m *DefaultManager) deleteRoute(id string) error {
	response, err := m.request(http.MethodDelete, m.adminURL+"/id/"+id, nil)
	if err != nil {
		return err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode >= 400 && response.StatusCode != http.StatusNotFound {
		return responseError(response)
	}
	return nil
}

func (m *DefaultManager) snapshotApps() ([]byte, bool, error) {
	response, err := m.request(http.MethodGet, m.adminURL+"/config/apps", nil)
	if err != nil {
		return nil, false, err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode == http.StatusNotFound {
		return nil, false, nil
	}
	if response.StatusCode >= 400 {
		return nil, false, responseError(response)
	}
	body, err := io.ReadAll(response.Body)
	return body, true, err
}

func (m *DefaultManager) restoreApps(snapshot []byte, existed bool) error {
	method := http.MethodPatch
	if !existed {
		method = http.MethodDelete
	}
	response, err := m.request(method, m.adminURL+"/config/apps", snapshot)
	if err != nil {
		return err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode >= 400 && response.StatusCode != http.StatusNotFound {
		return responseError(response)
	}
	return nil
}

func (m *DefaultManager) restoreAppsIfUnchanged(committed []byte, committedExists bool, snapshot []byte, snapshotExists bool) error {
	current, currentExists, err := m.snapshotApps()
	if err != nil || currentExists != committedExists || !bytes.Equal(bytes.TrimSpace(current), bytes.TrimSpace(committed)) {
		return err
	}
	return m.restoreApps(snapshot, snapshotExists)
}

func (m *DefaultManager) request(method, url string, body []byte) (*http.Response, error) {
	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return m.http.Do(request)
}

func (m *DefaultManager) alias(service string) string {
	return m.project() + "-" + service
}

func (m *DefaultManager) routeID(route route) string {
	return fmt.Sprintf("kool-%s-%s-%d-%d", m.project(), route.Service, route.Listen, route.Target)
}

func (m *DefaultManager) tlsID() string {
	return "kool-" + m.project() + "-tls"
}

func (m *DefaultManager) routeHosts(route route) []string {
	base := m.env.Get("KOOL_PROXY_HOST")
	seen := make(map[string]bool)
	var hosts []string
	for _, host := range route.Hosts {
		switch {
		case host == "@":
			host = base
		case host == "*":
			host = "*." + base
		case !strings.Contains(host, "."):
			host += "." + base
		}
		if !seen[host] {
			hosts = append(hosts, host)
			seen[host] = true
		}
	}
	return hosts
}

func (m *DefaultManager) project() string {
	if m.projectOverride != "" {
		return m.projectOverride
	}
	if project := m.env.Get("KOOL_WORKSPACE_PROJECT"); project != "" && m.env.IsTrue("KOOL_WORKSPACE") {
		return project
	}
	if project := m.env.Get("KOOL_WORKSPACE_SOURCE_PROJECT"); project != "" {
		return project
	}
	if project := m.env.Get("COMPOSE_PROJECT_NAME"); project != "" {
		return project
	}
	return m.env.Get("KOOL_NAME")
}

func responseError(response *http.Response) error {
	body, _ := io.ReadAll(response.Body)
	return fmt.Errorf("caddy Admin API returned %s: %s", response.Status, strings.TrimSpace(string(body)))
}
