To build:

`go build -o /usr/local/bin/kool`

Cross-compile:

`env GOOS=linux GOARCH=amd64 go build -o main_linux_amd64 main.go`

### Env vars

- `KOOL_TTY_DISABLE`: sets to `1` or `true` will make `kool exec` disable TTY for container interaction commands.

### Presets parsing

The preset files are managed at the presets/ folder. After any changes on those files you are required to run parse_presets.sh.

---

References:

- CLI tooling - https://github.com/spf13/cobra
- Environment variables - https://github.com/joho/godotenv
-
