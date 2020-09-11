# kool

[![Go Report Card](https://goreportcard.com/badge/github.com/kool-dev/kool)](https://goreportcard.com/report/github.com/kool-dev/kool)
![Docker Hub](https://github.com/kool-dev/kool/workflows/docker/badge.svg)
![Golang CI Lint](https://github.com/kool-dev/kool/workflows/golangci-lint/badge.svg)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkool-dev%2Fkool.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkool-dev%2Fkool?ref=badge_shield)

### Development workspaces made easy

Dev environment made easy, a standardized way for running applications no matter the stack on your local machine and deploying it to a development environment.

Run any stack / tool with any version, powered by Docker and Docker Compose in a simple way avoiding you to install lots of stuff on your machine.

Have the same feeling working on multiple projects with different stacks.

## Documentation

Full documentation at **https://kool.dev/docs** or at [docs/](docs/).

## Demo

<a href="https://www.youtube.com/watch?v=c4LonyQkFEI" target="_blank" title="Click to see full demo">
    <img src="https://user-images.githubusercontent.com/347400/87970968-fad10c80-ca9a-11ea-9bef-a88400b01f2c.png" alt="kool - demo" style="max-width:100%;">
</a>

## Installation

Kool is powered by [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/), you need to have it installed on your machine.

The run the follow script to install `kool` bin in your machine.

```bash
curl -fsSL https://raw.githubusercontent.com/kool-dev/kool/master/install.sh | bash
```
In case you need sudo:

```bash
curl -fsSL https://raw.githubusercontent.com/kool-dev/kool/master/install.sh | sudo bash
```

## For Windows

Download the installer [here](https://github.com/kool-dev/kool/releases)

## Usage

To help learning how to use kool we've built presets with good starting point for some popular stacks, feel free to open a PR in case you miss one.

### Presets

- [Adonis](docs/2-resets/Adonis.md)
- [Laravel](docs/2-Presets/Laravel.md)
- [NextJS](docs/2-Presets/NestJS.md)
- [NextJS](docs/2-Presets/NextJS.md)
- [NuxtJS](docs/2-Presets/NuxtJS.md)
- [NuxtJS](docs/2-Presets/Symfony.md)

### Examples

You can see projects using it here: https://github.com/kool-dev/examples

## Contributing

Like what you see? You are most welcome to contribute! We are working in a backlog of issues, feel free to browse through and enter discussions or get to work!

The flow is not written in stone, so you may go ahead and fork, code and PR with clear and direct communication on the feature/improvement or fix you developed.

PS: our main pain point at this moment is the lack of testing. Might be a great starting point.

### Lint

Before submitting a PR make sure to run `fmt` and linters.

```bash
kool run lint
```

### Updating commands signature documentation

The Cobra library offers a simple solution for getting markdown documentation for usage of its commands. In order to generate update the generated markdown files do as follow:

```bash
kool run make-docs
git add .
git commit -m "Updated commands docs"
```

## Security

If you find security issue please let us know ahead of making it public like in an issue so we can take action as soon as possible. Please email the concern to `fabricio.souza@fireworkweb.com`.

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
