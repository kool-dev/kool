# fwd2 - powered by Golang

To build:

`go build -o /usr/local/bin/kool`

Cross-compile:

`env GOOS=linux GOARCH=amd64 go build -o main_linux_amd64 main.go`

### Env vars

- `KOOL_TTY_DISABLE`: sets to `1` or `true` will make `kool exec` disable TTY for container interaction commands.

---

References:

- CLI tooling - https://github.com/spf13/cobra
- Environment variables - https://github.com/joho/godotenv
-
