# CI Integration

Beyond local development, **kool** can also run inside CI/CD pipelines. A pre-built Docker image is published to Docker Hub as [**kooldev/kool**](https://hub.docker.com/r/kooldev/kool) for exactly this purpose — it bundles the `kool` CLI together with `docker`, `docker compose`, `git`, and `bash`, pre-configured to talk to a Docker-in-Docker sidecar.

> This image is intended for **CI pipelines** that need `kool` alongside a Docker daemon. For local development, install the native binary as described in the [Installation guide](/docs/getting-started/installation).

## When to use the image

The `kooldev/kool` image is primarily used with **GitLab CI**, where the `services:` + `docker:dind` pattern pairs naturally with **kool**'s requirement for a working `docker` and `docker compose`. It also works with GitHub Actions, Drone, Bitbucket Pipelines, and any other CI runner that supports linking a Docker daemon service container.

If your CI runner already has Docker available on the host (e.g. GitHub Actions' default `ubuntu-latest` runner), installing **kool** via `curl | bash` at the top of the job is often simpler than running it through the image.

## What's in the image

- `kool` CLI — built with Go 1.25, statically linked
- `docker` CLI
- `docker compose` plugin (Compose V2)
- `git` and `bash`
- `ENV DOCKER_HOST=tcp://docker:2375` — defaults the Docker CLI to a sidecar service reachable at the hostname `docker`

Published tags track **kool** releases, e.g. `kooldev/kool:3.2.0` pins to a specific release and `kooldev/kool:3` tracks the current major line. Pin a specific version in CI for reproducibility.

## GitLab CI example

```yaml
test:
  image: kooldev/kool:3
  services:
    - name: docker:dind
      alias: docker
  variables:
    DOCKER_TLS_CERTDIR: ""
    DOCKER_HOST: tcp://docker:2375
  script:
    - kool run test
```

The `alias: docker` is important — the `kooldev/kool` image defaults `DOCKER_HOST` to `tcp://docker:2375`, so the DinD sidecar must be reachable at the hostname `docker`.

## GitHub Actions example

GitHub Actions runners already have Docker available, so you can usually skip the image and install **kool** natively:

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install kool
        run: curl -fsSL https://kool.dev/install | bash
      - run: kool run test
```

If you prefer a consistent, pinned `kool` version via `container:`, the pattern below mirrors the GitLab setup:

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    container:
      image: kooldev/kool:3
    services:
      docker:
        image: docker:dind
        options: --privileged
        env:
          DOCKER_TLS_CERTDIR: ""
    env:
      DOCKER_HOST: tcp://docker:2375
    steps:
      - uses: actions/checkout@v4
      - run: kool run test
```

## Known security caveat

The `kooldev/kool` image is built `FROM docker:29-cli`, which bundles `docker-compose` and `docker-buildx` CLI plugins as pre-built Go binaries that we inherit as-is. At any given time those plugins may have been compiled with a Go toolchain that has since received CVE advisories — for example a `go1.25.8` build of the bundled compose plugin triggers `CVE-2026-27143` until the upstream [`docker-library/docker`](https://github.com/docker-library/docker) image is rebuilt with a newer Go toolchain.

**The advisory is in the Go standard library of Docker's own plugin binaries, not in kool.** For security-sensitive pipelines, consider either:

- Installing **kool** as a native binary on the runner (`curl -fsSL https://kool.dev/install | bash`) instead of using the image, which avoids the upstream plugin surface entirely,
- Pinning a specific `kooldev/kool:VERSION` tag and re-scanning when you adopt a new version.

If you find a security issue in **kool** itself, please follow the [Security](https://github.com/kool-dev/kool#security) instructions in the main README.
