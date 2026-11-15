# Shared Docker image CI

The active image repositories use the same release boundary, implemented with
Bash and `jq` here in `kool-dev/kool`. Keep image-specific builds, smoke tests,
scanners, supported tags and dispatch recipients in each image repository.

## Pipeline

1. Run generation checks once. Export the complete expected stable-tag list.
2. Build/test on native `ubuntu-24.04` and/or `ubuntu-24.04-arm` jobs through
   `image-build`. Its Docker driver keeps fresh same-job parents available.
3. Pull upstream parents only for base builds. Do not force-pull derived parents.
   Use `image-parent` to assert the tested base ID and full layer prefix.
4. After **all** tests/scans in a native job pass, call `image-stage` with a JSON
   array of `{tag,id}` pairs. It checks the local IDs/platforms, tags those IDs,
   pushes them once and captures Docker's returned digests without remote
   inspections. It uploads one run/attempt/group/architecture artifact.
5. A single `image-publish.yml` job needs the entire native matrix and checks.
   It validates every expected tag/platform record before login or stable writes,
   then assembles manifests from those exact immutable digests. Single-arch
   releases keep a single image manifest. Node can retain its historical
   `-x86_64`/`-aarch64` aliases through `architecture-aliases`.
6. Repository-specific downstream dispatch jobs need successful publication.

Pin **all actions, the reusable workflow and its `tooling-ref` input to the same
40-character kool commit**. This prevents a branch update from silently changing
other repositories' release code. Update pins deliberately after testing shared
changes. The publisher checks out this pinned tooling, not mutable caller files.

## Publication guards

Call staging/publication only for the canonical repository's default branch and
non-PR events. The shared Bash guard repeats this check before login/staging.
Do not pass credentials into build/test steps or PR checks. The reusable workflow
only accepts push, schedule, manual and repository-dispatch releases.

Artifacts are scoped to the current workflow run **and attempt**. The publisher
rejects missing/duplicate/foreign/stale/malformed records and requires the exact
expected tag × architecture matrix. It trusts platform/ID checks from the tested
local images; it does not re-read every staged image's remote manifest/config.
Docker manifest assembly still performs registry reads and can hit pull quotas.

**Recovery remains Re-run all jobs.** There is no resume workflow or cross-run
artifact reuse. Stable-tag updates are sequential, not atomic across a release.
On partial publication, downstream notification stays blocked. Run-scoped staging
tags accumulate; these actions do not implement registry cleanup.

Keep existing architecture support: PHP, Swoole, WordPress and Node use amd64 +
arm64; NGINX, SQL Server and OCI8 remain amd64-only until their dependencies are
separately validated for ARM. This standardizes CI, not runtime support policy.

## Tests

```bash
bash .github/actions/image-stage/tests/release.sh
bash .github/actions/image-stage/tests/registry.sh
```

The first test mocks Docker and checks guards, immutable IDs, negative metadata,
absence of extra registry inspections and fresh-parent assertions. The second
uses only an ephemeral loopback registry and synthetic scratch images to verify
actual push digest capture, two-platform membership and single-platform output.
Scratch fixtures do not execute ARM application code; consumers' native CI does.
