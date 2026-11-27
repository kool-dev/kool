<!-- kool:begin (written by kool; edits between these markers are replaced on update) -->
## Kool: local Docker and Kool Cloud deploys

**Local dev** (defined in docker-compose.yml and kool.yml):
- `kool start` / `kool stop` / `kool status`
- `kool exec <service> <cmd>`
- `kool logs -f <service>`
- `kool run <script>`; list scripts with `kool run --json`

**Deploy to Kool Cloud.** The full guide is the `kool-cloud-deploy` skill (github.com/kool-dev/kool, `skills/kool-cloud-deploy`). Docs: https://kool.dev/docs/kool-cloud/kool-cloud-yml-reference

- **Files.** Write them by hand: `kool cloud setup` needs a TTY.
  - `kool.cloud.yml` at the repo root, plus a production `Dockerfile` and a `.dockerignore` (`.git`, `.env*`, env files, `node_modules`, build output, `*.db`).
  - No docker-compose.yml is needed. If one exists, its services deploy too, unless kool.cloud.yml sets `deploy: false` on them.
  - Minimal config: `services: {app: {build: ., expose: <port>, public: true}}`. `build` is a directory, not a Dockerfile path. `expose` is an integer and is required by `public`. The app must listen on 0.0.0.0:<port>.
- **Token.** It comes from `KOOL_API_TOKEN` in the environment or the git-ignored `.env`. Never pass `--token`, never print the token, never ask for it in chat.
- **Deploy.** Run `kool cloud deploy --domain <unique-name>.kool.cloud` from the repo root, in the foreground, with the longest tool timeout (20 minutes or more). Never background it, and never end the session while it runs.
- **Architecture.** The shared cloud is linux/amd64, the CLI default: don't pass `--platform`. For a BYOC cluster, pass `--platform linux/arm64` only if the user says its nodes are arm64. On Apple Silicon, build JS/assets in a `FROM --platform=$BUILDPLATFORM` stage, because QEMU can segfault.
- **Never `--verbose`.** It prints the registry password and the environment's variables.
- **Hooks.**
  - `after` runs once per deploy, after the rollout.
  - `before` is an init container in every new pod: it runs on the first deploy, on restarts and in each replica.
  - Hook, `daemons` and `cron` commands are passed as arguments to the image ENTRYPOINT, so start the app with `CMD`.
- **Debug.**
  - `kool cloud logs --domain <d> <service>` (daemons: `-c daemon0`) and `kool cloud exec --domain <d> <service> <cmd>`.
  - On `deploy failed`, read `.error` from `GET https://kool.dev/api/deploy/<ID>/status`. The ID is in the printed `Logs available at:` URL; the skill shows a token-safe curl.
  - "Pods not ready" or `context deadline exceeded` means the new pods weren't ready within 5 minutes and the deploy rolled back. Common causes:
    - a wrong-arch image;
    - the app bound to 127.0.0.1 or to another port;
    - a missing env var or database;
    - an out-of-memory kill at the default `micro` size (set `size: small`);
    - a crashing daemon.
- **`persists`.** It gives one disk per service and makes each deploy stop-then-start. Never set `sync: true` on a database path.
- **Verify.** Curl the URL, then a real page or API call. `Deploy finished` only means the rollout finished.
<!-- kool:end -->
