# kool

[![Kool.dev](https://kool.dev/img/logo.png)](https://kool.dev)

---

[![Go Report Card](https://goreportcard.com/badge/github.com/kool-dev/kool)](https://goreportcard.com/report/github.com/kool-dev/kool)
[![codecov](https://codecov.io/gh/kool-dev/kool/branch/master/graph/badge.svg)](https://codecov.io/gh/kool-dev/kool)
![Docker Hub](https://github.com/kool-dev/kool/workflows/docker/badge.svg)
![Golang CI Lint](https://github.com/kool-dev/kool/workflows/golangci-lint/badge.svg)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkool-dev%2Fkool.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkool-dev%2Fkool?ref=badge_shield)
[![Maintainability](https://api.codeclimate.com/v1/badges/1511f826de92d2ab39cc/maintainability)](https://codeclimate.com/github/kool-dev/kool/maintainability)

**kool** is a CLI tool that helps bringing down to earth the complexities of modern software development environments - making them lightweight, fast and reproduceable. It takes off the complexity and learning curve of _Docker_ and _Docker Compose_ for local environments, as well as offers a highly simplified interface for leveraging Kubernetes cloud deployment for staging and production deployments.

Get your local development environment up and running easy and quickly, put time and effort on making a great application, and then leverage the Kool cloud for deploying and sharing your work with the world! This tool is suitable for single developers or large teams, powering them with a simple start and still provide all flexibility the DevOps team needs to tailor up everything.

### Why adopt kool for your development environment?

- Provides out-of-the-box simple and fast development environments.
- No problems with running several projects with different versions and dependency needs.
- Do not install other project specific dependency ever again in your machine.
- Removes the learning curve of managing Docker and Docker Compose yourself (yet does not remove one inch of flexibity if you know your way around them).
- Standard tool for different stacks - helps micro-services teams to communicate and navigate amongst projects.

## Installation

Kool is powered by [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/), you need to have them already installed on your machine.

#### For Linux or MacOS

In order to obtain `kool` under **Linux** and **MacOS** run the following script:

```bash
curl -fsSL https://raw.githubusercontent.com/kool-dev/kool/master/install.sh | sudo bash
```

#### For Windows

Download and run the latest installer from our releases artifacts [here](https://github.com/kool-dev/kool/releases).

## Getting started

It is easy to get started leveraging `kool`. Provided you have all requirements (Docker and Docker Compose), to get started in an existing Laravel project you only need to:

```console
$ cd my-laravel-project/
$ kool preset laravel
$ # make sure your `.env` points to the proper database and Redis hosts (`database` and `cache`)
$ kool start
$ kool run reset
```


- There you go! Now you have a PHP 7.4, Mysql and Redis environment. You are encouraged to take a look and make changes you see fit at `docker-compose.yml` or `kool.yml` to better adjust your project specifications.

- The steps above will create some configuration files in your project folder - all of which you should commit to your version control system.

## Documentation

You can check the documentation at **https://kool.dev/docs** or at [docs/](docs/).


## Frameworks Presets

To help getting you started we've built presets as a starting point for some popular stacks and frameworks. In case you miss one let us know in an issue or feel free to open up a PR for it!

Out current presets:

- [Laravel](docs/2-Presets/Laravel.md)
- [Symfony](docs/2-Presets/Symfony.md)
- [Wordpress](docs/2-Presets/Wordpress.md)
- [Adonis](docs/2-Presets/Adonis.md)
- [NestJS](docs/2-Presets/NestJS.md)
- [NextJS](docs/2-Presets/NextJS.md)
- [NuxtJS](docs/2-Presets/NuxtJS.md)

### Examples

You can see projects using it here: https://github.com/kool-dev/examples

## Demo

<a href="https://www.youtube.com/watch?v=c4LonyQkFEI" target="_blank" title="Click to see full demo">
    <img src="https://user-images.githubusercontent.com/347400/87970968-fad10c80-ca9a-11ea-9bef-a88400b01f2c.png" alt="kool - demo" style="max-width:100%;">
</a>

## Contributing

You are most welcome to contribute! There are issues with [`good first issue`](https://github.com/kool-dev/kool/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) label, feel free to browse through and enter discussions or get to work!

The workflow is not written in stone, so you may go ahead and fork, code and PR with clear and direct communication on the feature/improvement or fix you developed.

### Roadmap

We have been working in a loosely defined but clear roadmap:

- Recently: we focused on tests coverage for moving forward condifently; we got from 0 to 90% coverage in a couple of weeks! *check!*
- Currently focusing in: improving overall UX and stabilize features - error messages, progress display, output control, presets creation, getting started, etc...
- Next steps:
    - Continunously improve tests coverage;
    - Continunously improve UX based on usage feedbacks;
    - Expand features (Proposal issues);

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
git commit -m "Updated commands docs"
```

## Security

If you find security issue please let us know ahead of making it public like in an issue so we can take action as soon as possible. Please email the concern to `fabricio.souza@fireworkweb.com`.

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.


[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkool-dev%2Fkool.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkool-dev%2Fkool?ref=badge_large)
