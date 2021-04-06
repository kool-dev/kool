# Contributing Guidelines

The Kool project accepts contributions via GitHub pull requests. This document outlines the process to help get your contribution accepted.

The workflow is on the making with the core team focused on shiping features and writting documentation at this point.

There are issues with [`good first issue`](https://github.com/kool-dev/kool/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) label, feel free to browse through and enter discussions or get to work!

## Reporting a Security Issue

As stated on [SECURITY.md](SECURITY.md), in case of a security issue or concern, please consider reporting it privately at first to [contact@kool.dev](mailto:contact@kool.dev).

### Rich content Issues and PRs

If applicable, please consider adding screenshots to help explain your issues or solutions. Recording can also be very helpful to communicate.

- Check out [ASCII Cinema](https://asciinema.org/) for recording your terminal only.
- For recording your whole screen check out [ShareX Opensource Screen Capture](https://getsharex.com) or similars (then upload it somewhere to share).

## Local Development directions
### Lint, formatting and tests

Before submitting a PR make sure to run `fmt` and linters.

```bash
kool run fmt
kool run lint
kool run test
```

### Updating commands signature documentation

The Cobra library offers a simple solution for getting markdown documentation for usage of its commands. In order to generate update the generated markdown files do as follow:

```bash
kool run make-docs
git add .
git commit -m "Updated commands docs (auto generated)"
```
