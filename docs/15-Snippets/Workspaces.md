# Workspaces and local proxy

Kool runs selected services in separate Compose projects for [Git worktrees](https://git-scm.com/docs/git-worktree) and [Rift](https://github.com/anomalyco/rift) workspaces. You choose which services to isolate and which to share.

For example, each workspace can run its own app and Vite server using shared PostgreSQL and Redis. Independent data requires a workspace-local database with separate connections and storage.

## Runnable Node + Vite example

Try [kool-worktree-example](https://github.com/kool-dev/kool-worktree-example): Node + Vite HMR, shared PostgreSQL/Redis, Rift hooks, tests, and an optional private database.

The example needs Docker Engine 28+, Compose 2.24.4+ (`!reset` / `!override`), and a Kool build with workspace/proxy support. Until released, follow its README to build the feature branch.

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

Edit `src/main.js` in `task-a` to see HMR there. Both apps see the shared database counter. Tests run in the current app container; use separate data for destructive integration tests.

Stop before removing a plain Git worktree:

```bash
kool stop
cd ../kool-worktree-example
git worktree remove ../task-a
```

Save or commit edits before removal. Hostnames follow workspace directory names, with a suffix for duplicate Git worktree basenames. Check project names with `kool status`.

## Rift lifecycle hooks

Start the original project first. This `.rift.toml` starts each Rift workspace after creation and stops it before removal:

```toml
version = 1

[[hooks.postcreate]]
run = "kool start"

[[hooks.preremove]]
run = "kool stop"
```

Use the same Kool commands inside Rift workspaces. Plain Git worktrees do not run these hooks.

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

Both features are opt-in: a non-empty `workspaces` list enables detection; `proxy` enables Caddy management. Without them, existing behavior is unchanged.

A port uses `listen:target`:

```text
"80:8080"
 host port : service container port
```

Each port applies to every host in its route:

```text
@                  proxy.domain
*                  wildcard before proxy.domain
api                api.proxy.domain
admin.example.test complete hostname
```

`hosts` defaults to `["@", "*"]` (base domain and wildcard). Use `["@"]` for the base domain alone.

In a workspace, `api` becomes `api.<workspace>.workspace.<domain>`. Complete hostnames stay unchanged.

`proxy.network` defaults to `KOOL_GLOBAL_NETWORK`, or `kool_global`.

Kool sets `KOOL_PROXY_DOMAIN` and the current context's `KOOL_PROXY_HOST` before Compose starts. Pass them to your app or Vite:

```yaml
environment:
  VITE_PUBLIC_ORIGIN: "http://${KOOL_PROXY_HOST}:3001"
```

No Kool or proxy labels are required in `docker-compose.yml`.

### Vite assets and HMR

Configure Vite to use the public origin passed by Compose:

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

Apps generating HTML must also use `VITE_PUBLIC_ORIGIN` for `/@vite/client` and application modules. Configure framework plugins' equivalent asset URL. Adjust CORS for non-default app ports. Caddy forwards WebSocket upgrades; HTTPS requires `wss` and a trusted local CA.

## Compose services

Mount the current directory and join the external shared network:

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

Caddy owns published proxy ports; Kool removes them from proxied services and clears fixed workspace `container_name` values.

Workspace starts use `--no-deps`. Start shared infrastructure in the original project first; `depends_on` will not start unselected services or wait for them. Prepare each workspace's dependencies/environment files and retry connections until services are ready.

Use project-specific shared aliases such as `demo-database`; generic names can collide across projects. Include code-dependent queue workers in `workspaces` as needed. Sharing Redis does not isolate worker code or jobs.

### Choosing separate data

Add a database to `workspaces` for independent data. Point the app at it over a project-local network or unique alias, and use project-scoped storage. Reusing the source connection string or a fixed external volume still shares data. Kool does not rewrite connections or clone data; the example includes a Compose override for this setup.

## Managed proxy

Kool manages a global `kool-proxy` container using `caddy:2.10-alpine`.

Caddy:

- Joins the shared Docker network.
- Publishes each configured listen port.
- Exposes its Admin API only at `127.0.0.1:2019`.
- Binds administration only to its dedicated admin-network interface, not application networks.
- Does not mount the Docker socket.
- Handles HTTP streaming and WebSocket upgrades, including Vite HMR.

Backend aliases include the project and service:

```text
example-app
example-node
example-workspace-task-a-<path-hash>-app
example-workspace-task-a-<path-hash>-node
```

`kool start` registers routes through Caddy's Admin API; stop removes the stopped projects' routes.

With `domain: app.localhost`:

```text
Source:
app.localhost             -> app
*.app.localhost           -> app

Workspace task-a:
task-a.workspace.app.localhost   -> task-a app
*.task-a.workspace.app.localhost -> task-a app
```

One hostname can route to the app on `80:80` and Vite on `3001:3001`.

New listen ports or proxy upgrades may recreate Caddy while preserving other projects' routes. Accepted configuration is saved privately in `~/.kool/proxy/caddy.json` after starts and route removals, then loaded on restart. Startup does not use Caddy's volume autosave.

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

This enables TLS on every configured listener. Caddy issues and renews source/workspace certificates, including wildcards. The `kool_proxy` volume stores autosaved routes in `config/` and CA/certificates in `data/`.

Trust the local CA on macOS or Linux:

```bash
kool start
kool proxy trust
```

Trust installation needs administrator privileges. On WSL it updates Linux only; import the CA into Windows for Windows browsers.

HTTPS does not add HTTP redirects. Configure the TLS ports clients should use.

## Commands

Start runs the full original project or the selected workspace services. Workspace project names include a canonical-path hash, such as `example-workspace-task-a-<path-hash>`; avoid hardcoding container names.

Use existing commands:

```bash
kool start
kool stop
kool restart
kool exec app php artisan test
kool run artisan test
kool logs app
kool status
```

`kool status` shows the original project plus all active workspaces from the source, or just the current workspace when run inside one.

Workspace stop removes only its containers and routes. An unqualified source-project stop shuts down its workspaces first, leaving unrelated projects running.

Workspace mode requires Compose's `!reset` tag. The proxy requires Docker 28+ for gateway priority to keep its published Admin API on the dedicated admin network.
