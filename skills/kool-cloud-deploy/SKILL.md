---
name: kool-cloud-deploy
description: Deploy an app to Kool Cloud (kool.dev) with `kool cloud deploy` - write kool.cloud.yml, a production Dockerfile and .dockerignore, deploy headless, read failures, verify with curl. Use when asked to deploy, ship, host or publish an app on Kool, Kool Cloud or a *.kool.cloud domain, or when the repo has a kool.cloud.yml. Local Docker work (kool start/exec/run) is the kool-cli skill.
---

# Kool Cloud deploy

## Before you start (ask the human; never work around these)
- `kool` CLI installed, Docker running (`docker info`).
- `KOOL_API_TOKEN` set in the environment, or in the git-ignored `.env` that the CLI loads from the working directory. A new environment needs a **project-scoped** token. Never ask for the token in chat, never echo it, never pass `--token` (it shows in `ps` and shell history).
- A domain: `<unique-name>.kool.cloud` has HTTPS at once; a custom domain needs its DNS pointed first (the panel shows the records).

## Files: write them by hand (`kool cloud setup` exits without a TTY)
`kool.cloud.yml` at the repo root (`kool.deploy.yml` is the legacy name), a production `Dockerfile`, a `.dockerignore`, and optionally a git-ignored env file. No docker-compose.yml is needed. Smallest valid config:
```yaml
# kool.cloud.yml
services:
  app:
    build: .        # a directory with a Dockerfile, never a Dockerfile path
    expose: 80      # integer; the port the app listens on, bound to 0.0.0.0
    public: true    # route the domain here; requires expose
```
- `build`: a directory, or `{context: ., dockerfile: Dockerfile.prod, args: {K: "v"}}`. The map form needs `context`, and `dockerfile` is relative to the repo root. `image: name:tag` instead of `build` deploys a prebuilt image.
- `public: true` on one service. A list `[{port: 80}, {port: 3000, path: /ws}]` routes a path prefix to another port, and strips the prefix.
- `environment`: a map (plain env vars) or an env-file name (stored as a Secret). `"{{VAR}}"` is replaced by the environment variable VAR set in the Kool panel. Always quote placeholders; an unknown VAR stays literal.
- `size`: `micro` (the default: 128m CPU, 256Mi RAM), `small`, `medium`, `large`, `xlarge`, `2xlarge`. Each size doubles the one before, trials cap at `medium`, and cron and hooks default to `medium`. `replicas: N` adds pods.
- `persists: {size: 1Gi, paths: [{path: /app/data}]}` gives one disk per service. It makes every deploy stop-then-start (seconds of downtime) and pins the pods to one node. `sync: true` copies the image's copy of the path over the disk at **every pod start**: use it for skeleton folders, never for a database file. Relative paths need `root: /app`.
- `daemons: [{command: [...], size: small}]` are long-running extra containers in the same pod (queue workers), named `daemon0`, `daemon1`. If one crash-loops, the deploy fails.
- `cron: [{schedule: "*/5 * * * *", command: [...]}]` is a Kubernetes CronJob: a fresh pod per run, on the UTC clock.
- `before: [[cmd, arg]]` runs as an init container of **every new pod**. That includes the first deploy (the docs saying otherwise are stale), but also every restart, reschedule and replica. If it fails, the pod never starts and the deploy fails.
- `after: [[cmd, arg]]` runs **once per deploy**, as a Job after the new pods are up. If it fails, the release rolls back.
- Hooks, daemons and cron run with the same image, env and persisted paths. Their command is passed as **arguments to the image's ENTRYPOINT**, so start your app with `CMD`, and keep any `ENTRYPOINT` one that `exec`s its arguments.
- Services reach each other by service name (`cache:6379`).
- **If a docker-compose.yml exists, it is uploaded and merged under kool.cloud.yml.** Every compose service deploys unless kool.cloud.yml sets `deploy: false` on it, and compose keys you don't override (`image`, `command`, `environment`) carry over.

## Rules for every stack
- **Architecture.** The shared Kool Cloud runs linux/amd64, which is the CLI's default build, so don't pass `--platform`. BYOC deploys go to the team's default cluster, or to the one named by `KOOL_CLOUD_CLUSTER=<slug>`. A BYOC cluster may be arm64 (Graviton): pass `--platform linux/arm64` only when the user says its nodes are arm64. A wrong-arch image is accepted, and fails about 5 minutes later with no useful message.
- On Apple Silicon the amd64 build runs under QEMU, which is slow, and Turbopack, Rolldown and esbuild can segfault (`qemu: uncaught target signal 11`). Build arch-independent output (JS bundles, assets) in a `FROM --platform=$BUILDPLATFORM` stage, and install native dependencies (`npm ci`, apt, pecl) in the final stage.
- **Add a `.dockerignore`**: `.git`, `.env*`, the env file, `node_modules`, build output, `*.db`, and `vendor` for PHP. Without it, secrets and host binaries end up in the image.
- **Deploy** with `kool cloud deploy --domain <domain>` from the repo root, in the **foreground**, with your shell tool's longest timeout (20 minutes or more; a deploy usually takes 3 to 8 minutes, longer under QEMU). Never background it, and never end your turn while it runs. Interrupting it before `Deploy started.` fails the deploy; after that line the rollout continues on the server.
- **Never use `--verbose`.** It prints the registry password and every environment variable. Never paste output that holds secrets.
- **Logs and shell:** `kool cloud logs --domain <d> <service>` (add `-c daemon0`, `-f` or `--tail 100`), and `kool cloud exec --domain <d> <service> <cmd>` (don't add `--`). Both see running pods only. They use your local `kubectl`, or else pull the `kooldev/toolkit:full` image.
- **Verify.** `Deploy finished` means the rollout finished, not that the app works. Run `curl -sS -o /dev/null -w '%{http_code}\n' https://<domain>/`, then a real page or API call, checking its content.

## When it fails
| CLI prints | Do |
|---|---|
| `could not find required file 'kool.cloud.yml'`, `no kool.cloud.yml config files found` | Run from the repo root and create the file. |
| `build entry 'X' is not a directory` | Use `build: .` or the map form. |
| `defines a build directory ('X') that does not exist` | Fix the path; it is relative to the repo root. |
| ``is public, but it does not define the `expose` entry`` | Add `expose: <port>` as an integer. |
| `yaml: ...`, `cannot unmarshal` | Fix the YAML: `expose` must be an int and `{{ }}` must be quoted. |
| `missing deploy domain` | Pass `--domain <domain>`. |
| `missing Kool.dev Cloud API token` | Ask the human to export `KOOL_API_TOKEN`. |
| `setup command is not available in non-interactive mode` | Write the files yourself. |
| `file not found, not including on tarball: X` | Create X (the env file) first, or the deploy fails on the server. |
| `401 - Unauthenticated.` | The token is wrong or revoked: ask the human. |
| `402 - Your trial has expired...` | The human must add a payment method. |
| `token > ...project-scoped token` | Ask for a project-scoped token. |
| `token > This token is scoped to environment 'X'` | Deploy to domain X. |
| `environment > Token does not have permission for this environment.` | Another team owns the domain: pick another name. |
| `deploy > Kool Cloud no longer builds images...` | Upgrade the kool CLI. |
| A Docker build error | Fix the Dockerfile; reproduce with `docker build .`. |
| `qemu: uncaught target signal 11`, `exec format error` during the build | Move that step into a `$BUILDPLATFORM` stage. |
| `Cannot connect to the Docker daemon` | Ask the human to start Docker. |
| `deploy failed` | Read the reason; see below. |
| `timeout waiting deploy to finish` | The CLI stopped waiting (15 minutes; 3.9.0 also stops on one failed status check). Check the status and curl the URL before redeploying. |
| `500 - Server Error` from `logs` or `exec` | Nothing is running for that service yet, for example because the first deploy rolled back. |

**Reading why a deploy failed.** Don't use `--verbose`. Take the ID from `Logs available at: https://kool.dev/deploys/<ID>?source=cli` and run this (the `printf` keeps the token out of argv):
```bash
curl -sS -H @<(printf 'Authorization: Bearer %s\n' "$KOOL_API_TOKEN") -H 'Accept: application/json' \
  "https://kool.dev/api/deploy/<ID>/status" | jq -r '.status, .error'
```
Without `jq`, drop the pipe: the JSON is small. Newer CLIs print the reason after `deploy failed:` themselves.
- `context deadline exceeded`, or pods not ready: the new pods were not ready within 5 minutes and the deploy rolled back. A container must accept TCP on `expose` within about 40 seconds of starting. Suspects, in order:
  - a wrong-arch image;
  - the app listening on 127.0.0.1 or on another port;
  - a missing env var or database;
  - OOM at `micro` (raise `size`);
  - a crashing daemon or a failing `before` hook;
  - an `ENTRYPOINT` that ignores its arguments.

  The failed pod is gone, so reproduce locally: `docker build -t smoke . && docker run --rm -p 8080:<expose> smoke`, then `docker rmi smoke`.
- `BackoffLimitExceeded`: an `after` hook failed. Run its command locally.
- `Malformed inline YAML` near `{{`: quote the placeholder.

## Stack references (read the one that matches)
- Laravel: [references/laravel.md](references/laravel.md)
- Static HTML: [references/static.md](references/static.md)
- Vite/React SPA: [references/spa.md](references/spa.md)
- Next.js, with Prisma and SQLite: [references/nextjs.md](references/nextjs.md)
- Node API (Express, Nest): [references/node.md](references/node.md)
- Docs: https://kool.dev/docs/kool-cloud/kool-cloud-yml-reference. Where the docs disagree with this file, trust this file.
