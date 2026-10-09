# Workspaces and local proxy

Kool detects workspaces managed by [Rift](https://github.com/anomalyco/rift) or linked [Git worktrees](https://git-scm.com/docs/git-worktree). Each workspace runs selected services in a separate Compose project. You choose which services get their own runtime and which infrastructure stays shared in the original project.

For example, a full-stack project can run its app and Vite server separately in each workspace while reusing PostgreSQL and Redis. Files and application containers are isolated; shared data is not. Features that need independent data can run a workspace-local database with the appropriate connection, network, and storage configuration.

## Runnable Node + Vite example

The [kool-worktree-example repository](https://github.com/kool-dev/kool-worktree-example) includes a Node app, Vite with workspace-aware HMR, shared PostgreSQL/Redis, Rift hooks, tests, and an optional private-database configuration.

Use a Kool build containing workspace/proxy support, Docker Engine 28+, and Docker Compose 2.24.4+ for the example's `!reset` / `!override` tags. While the feature is unreleased, the example README explains how to build the feature branch.

```bash
git clone https://github.com/kool-dev/kool-worktree-example.git
cd kool-worktree-example
kool start
# Wait for the app and Vite to be ready, then open http://demo.localhost.

git worktree add ../task-a -b task-a
cd ../task-a
kool start
# Open http://task-a.workspace.demo.localhost once services are ready.
kool run test
kool status
```

Edit `src/main.js` inside `task-a` to see HMR update only that workspace. The demo counter deliberately uses a shared database, so both apps see the same count. Tests run in the current workspace's app; integration tests still use the database you configured. Use separate test data for destructive tests.

Stop before removing a plain Git worktree:

```bash
kool stop
cd ../kool-worktree-example
git worktree remove ../task-a
```

Save or commit changes before removing a worktree. A Git branch name does not determine its URL: the workspace directory name does. Duplicate Git worktree basenames receive a suffix to distinguish their hosts. `kool status` shows the active context and project names.

## Rift lifecycle hooks

Start the original project's infrastructure first, then create a Rift workspace using your Rift CLI or integration. Add this `.rift.toml` to automate its lifecycle:

```toml
version = 1

[[hooks.postcreate]]
run = "kool start"

[[hooks.preremove]]
run = "kool stop"
```

The hooks start selected services after creation and stop them before removal. Inside the Rift workspace, use the same `kool run`, `kool exec`, `kool logs`, and `kool status` commands. Plain Git worktrees do not execute Rift hooks.

## Configuration

Configure workspace services and proxy routes once in `kool.yml`:

```yaml
workspaces:
  - app
  - node

proxy:
  domain: "${APP_DOMAIN:-app.localhost}"
  routes:
    app:
      ports: ["80:80"]
      hosts: ["@", "*"]
    node:
      ports:
        - "3001:3001"

scripts:
  # Existing scripts remain here.
```

Both features are opt-in. Without a non-empty `workspaces` list, Rift and Git worktree detection is skipped and existing Kool commands retain their legacy behavior. Without `proxy`, Kool does not inspect or manage the global Caddy proxy.

A port uses `listen:target`:

```text
"80:8080"
 host port : service container port
```

Every port is mapped across every host in the service route. Host values use these conventions:

```text
@                  proxy.domain
*                  wildcard before proxy.domain
api                api.proxy.domain
admin.example.test complete hostname
```

`hosts` is optional and defaults to `["@", "*"]`, routing both the base domain and wildcard subdomains. Set it explicitly to `["@"]` when wildcard routing is not wanted.

Relative hosts use the workspace base automatically. For example, `api` becomes `api.<workspace>.workspace.<domain>`. A complete hostname remains unchanged.

`proxy.network` optionally changes the shared Docker network. It defaults to `KOOL_GLOBAL_NETWORK`, which defaults to `kool_global`.

Kool injects `KOOL_PROXY_DOMAIN` and the context-specific `KOOL_PROXY_HOST` before Compose starts services. They can be passed into application or Vite configuration:

```yaml
environment:
  VITE_PUBLIC_ORIGIN: "http://${KOOL_PROXY_HOST}:3001"
```

No Kool or proxy labels are required in `docker-compose.yml`.

### Vite assets and HMR

Passing an origin into Compose is not enough by itself: Vite and the application's generated asset URLs must use it. A minimal `vite.config.js` is:

```js
import { defineConfig } from 'vite';

const origin = new URL(process.env.VITE_PUBLIC_ORIGIN);

export default defineConfig({
  server: {
    host: '0.0.0.0',
    port: 3001,
    strictPort: true,
    origin: origin.origin,
    allowedHosts: [origin.hostname],
    // App HTML is served on the default HTTP/HTTPS port in this setup.
    cors: { origin: `${origin.protocol}//${origin.hostname}` },
    hmr: {
      host: origin.hostname,
      protocol: origin.protocol === 'https:' ? 'wss' : 'ws',
      clientPort: Number(origin.port || (origin.protocol === 'https:' ? 443 : 80)),
    },
  },
});
```

Pass `VITE_PUBLIC_ORIGIN` to the app service as well if it generates HTML loading Vite assets. Load both `/@vite/client` and application modules from that origin, not the original project's hostname or `localhost`. For framework plugins, configure their development-server/asset URL equivalent. Adjust the CORS origin if the app uses a non-default port. Caddy forwards WebSocket upgrades; HTTPS origins need `wss` and a trusted local CA.

## Compose services

Workspace services should mount the current project normally and join the external shared network:

```yaml
services:
  app:
    image: example/app
    volumes:
      - .:/app:delegated
    networks:
      - kool_global

  node:
    image: node:22
    working_dir: /app
    volumes:
      - .:/app:delegated
    networks:
      - kool_global

  database:
    image: mysql:8
    networks:
      - kool_global

networks:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"
```

Kool removes published ports from proxied services because Caddy owns the listen ports. It also removes fixed `container_name` values from workspace services. Shared infrastructure must already be running from the original workspace.

Workspace starts use `--no-deps`: `depends_on` does not start unselected services or wait for their readiness. Prepare dependencies/environment files for each workspace, and give applications connection retries or wait until infrastructure is ready.

Use unique shared-network aliases for infrastructure (for example, `demo-database`) and configure the app to connect to them. Generic aliases such as `database` can be ambiguous when unrelated projects share the same global network. Code-dependent queue workers usually belong in `workspaces` too; a shared Redis queue backend does not isolate worker code or jobs.

### Choosing separate data

The `workspaces` list can include a database when a feature needs its own data. Also point the workspace app at that database, keep its connection on a project-local network (or use a unique workspace alias), and use project-scoped storage. Do not reuse a fixed external volume or the original database's connection string. Selecting a service does not automatically rewrite application connections or clone data. The example repository includes a Compose override showing this configuration.

## Managed proxy

When proxy routes are configured, Kool manages a global `kool-proxy` container using the pinned `caddy:2.10-alpine` image.

Caddy:

- Joins the shared Docker network.
- Publishes each configured listen port.
- Exposes its Admin API only at `127.0.0.1:2019`.
- Binds administration only to its dedicated admin-network interface, not application networks.
- Does not mount the Docker socket.
- Handles HTTP streaming and WebSocket upgrades, including Vite HMR.

Kool gives proxied services deterministic network aliases such as:

```text
example-app
example-node
example-workspace-task-a-<path-hash>-app
example-workspace-task-a-<path-hash>-node
```

`kool start` registers routes through Caddy's Admin API. A workspace's `kool stop` removes its own routes. An unqualified stop in the original project first stops its owned workspaces, then the original project; unrelated projects keep their routes.

For `domain: app.localhost`, routes are generated as follows:

```text
Source:
app.localhost             -> app
*.app.localhost           -> app

Workspace task-a:
task-a.workspace.app.localhost   -> task-a app
*.task-a.workspace.app.localhost -> task-a app
```

Routes using different listen ports can use the same hostname. For example, `80:80` can target the application while `3001:3001` targets Vite.

Kool recreates the proxy when new listen ports or a proxy upgrade require it, preserving other projects' routes. Successful starts and route removals save the accepted configuration privately under `~/.kool/proxy/caddy.json`; container and host restarts reload it. Caddy's volume autosave is not used for startup.

## Local HTTPS

Enable Caddy's internal certificate authority with `proxy.https`:

```yaml
proxy:
  domain: "${APP_DOMAIN:-app.localhost}"
  https: true
  routes:
    app:
      ports: ["443:80"]
    node:
      ports: ["3001:3001"]
      hosts: ["@"]
```

All configured listeners use TLS when `https` is enabled. Caddy issues and renews certificates for source and workspace hosts, including wildcards. The `kool_proxy` Docker volume stores autosaved routes under `config/` and CA and certificate data under `data/`.

After starting the proxy, trust its local root CA on macOS or Linux:

```bash
kool start
kool proxy trust
```

Trust installation requires administrator privileges. On WSL, `kool proxy trust` updates the Linux trust store; browsers running on Windows also require importing the root certificate into the Windows certificate store.

HTTPS does not currently add an automatic HTTP-to-HTTPS redirect. Configure only the TLS listen ports clients should use.

## Commands

From the original workspace, `kool start` starts the project normally. From a Rift or linked Git worktree it starts only services listed under `workspaces`, using an internal project name such as `example-workspace-task-a-<path-hash>`. The hash comes from the canonical workspace path; do not hardcode container names.

Existing commands remain transparent:

```bash
kool start
kool stop
kool restart
kool exec app php artisan test
kool run artisan test
kool logs app
kool status
```

From the original project, `kool status` shows the main project and all active workspaces. Inside a workspace, it shows the main project and the current workspace.

`kool stop` inside a workspace removes only that workspace's containers and proxy routes. The original project and other workspaces remain running. Running `kool stop` without service arguments from the original project stops every active workspace before stopping the main project.

Workspace mode requires Docker Compose support for the `!reset` merge tag. The managed proxy also requires Docker 28 or newer to keep its published Admin API on the dedicated admin network using gateway priority.
