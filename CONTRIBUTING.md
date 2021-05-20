# Contributing Guidelines

The Kool project accepts contributions via GitHub pull requests. This document outlines the process to help get your contribution accepted.

At this point, the workflow is focused on supporting the core team with shipping new features and writing documentation.

There are issues with a [`good first issue`](https://github.com/kool-dev/kool/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) label. Feel free to browse open issues, enter discussions, or get straight to work!

## Reporting a Security Issue

As stated on [SECURITY.md](SECURITY.md), in case of a security issue or concern, please consider reporting it privately at first to [contact@kool.dev](mailto:contact@kool.dev).

### Rich Content Issues and PRs

If applicable, please consider adding screenshots to help explain your issues or solutions. Recordings can also be very helpful to facilitate communication.

- Check out [ASCII Cinema](https://asciinema.org/) for recording your terminal only.
- For recording your whole screen, check out [ShareX Opensource Screen Capture](https://getsharex.com), or similar tools, and upload your recordings somewhere to share.

## Local Development Directions

### Lint, Formatting and Tests

Before submitting a PR, make sure to run `fmt` and linters.

```bash
kool run fmt
kool run lint
kool run test
```

### Updating Command Signature Documentation

The Cobra library offers a simple solution for generating markdown documentation for the usage of its commands. In order to update the generated markdown files, you must do the following:

```bash
kool run make-docs
git add .
git commit -m "Updated commands docs (auto generated)"
```
