# Workspaces and local proxy

Kool detects workspaces managed by [Rift](https://github.com/anomalyco/rift) and linked [Git worktrees](https://git-scm.com/docs/git-worktree). Each workspace runs selected application services in a separate Compose project while databases, caches, and other infrastructure remain in the original project.

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

## Managed proxy

When proxy routes are configured, Kool manages a global `kool-proxy` container using the pinned `caddy:2.10-alpine` image.

Caddy:

- Joins the shared Docker network.
- Publishes each configured listen port.
- Exposes its Admin API only at `127.0.0.1:2019`.
- Does not mount the Docker socket.
- Handles HTTP streaming and WebSocket upgrades, including Vite HMR.

Kool gives proxied services deterministic network aliases such as:

```text
example-app
example-node
example-workspace-task-a-app
example-workspace-task-a-node
```

`kool start` registers routes through Caddy's Admin API. `kool stop` removes only the current source or workspace project's routes.

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

If an existing `kool-proxy` container does not publish a newly configured listen port, remove it and run `kool start` again:

```bash
docker rm -f kool-proxy
kool start
```

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

From the original workspace, `kool start` starts the project normally. From a Rift or linked Git worktree it starts only services listed under `workspaces`, using an internal project name such as `example-workspace-task-a`.

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

This mode requires Docker Compose support for the `!reset` merge tag.
