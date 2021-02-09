<p align="center"><a href="https://kool.dev" target="_blank"><img src="https://kool.dev/img/logo.png" width="400" alt="kool - cloud native"></a></p>


<p align="center">
<a href="https://goreportcard.com/report/github.com/kool-dev/kool"><img src="https://goreportcard.com/badge/github.com/kool-dev/kool" alt="Go Report Card"></a>
<a href="https://codecov.io/gh/kool-dev/kool"><img src="https://codecov.io/gh/kool-dev/kool/branch/master/graph/badge.svg" alt="codecov"></a>
<a href="https://github.com/kool-dev/kool/workflows/docker"><img src="https://github.com/kool-dev/kool/workflows/docker/badge.svg" alt="Docker Hub"></a>
<a href="https://github.com/kool-dev/kool/workflows/golangci-lint"><img src="https://github.com/kool-dev/kool/workflows/golangci-lint/badge.svg" alt="Golang CI Lint"></a>
<a href="https://app.fossa.com/projects/git%2Bgithub.com%2Fkool-dev%2Fkool?ref=badge_shield"><img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkool-dev%2Fkool.svg?type=shield" alt="FOSSA Status"></a>
<a href="https://codeclimate.com/github/kool-dev/kool/maintainability"><img src="https://api.codeclimate.com/v1/badges/1511f826de92d2ab39cc/maintainability" alt="Maintainability"></a>
<a href="https://join.slack.com/t/kool-dev/shared_invite/zt-jeh36s5g-kVFHUsyLjFENUUH4ucrxPg"><img src="https://img.shields.io/badge/Join%20Slack-kool--dev-orange?logo=slack" alt="Join Slack Kool community"></a>
</p>

## About `kool`

**kool** is a CLI tool that helps bringing down to earth the complexities of modern software development environments - making them lightweight, fast and reproduceable. It takes off the complexity and learning curve of _Docker_ and _Docker Compose_ for local environments, as well as offers a highly simplified interface for leveraging _Kubernetes_ cloud deployment for staging and production deployments.

Get your local development environment up and running easy and quickly, put time and effort on making a great application, and then leverage the Kool cloud for deploying and sharing your work with the world! This tool is suitable for single developers or large teams, powering them with a simple start and still provides all flexibility the DevOps team needs to tailor up everything.

### Why adopt kool for your development environment?

- Provides out-of-the-box simple and fast development environments.
- No problems with running several projects with different versions and dependency needs.
- Do not install other project specific dependency ever again in your machine.
- Removes the learning curve of managing Docker and Docker Compose yourself (yet does not remove one inch of flexibity if you know your way around them).
- Standard tool for different stacks - helps micro-services teams to communicate and navigate amongst projects.

## Installation

Requirements: Kool is powered by [Docker](https://docs.docker.com/get-docker/) so you need to have it already installed on your machine. If you haven't already, please [get Docker first](https://docs.docker.com/get-docker/).

#### For Linux or MacOS

In order to obtain `kool` under **Linux** and **MacOS** run the following script:

```bash
curl -fsSL https://kool.dev/install | bash
```

#### For Windows

Download and run the latest installer from our releases artifacts [here](https://github.com/kool-dev/kool/releases).

## Getting started

It is easy to get started leveraging `kool`.

To create a new Laravel project you only need to:

```console
$ kool create laravel my-laravel-project
$ cd my-laravel-project/
$ # make sure your `.env` points to the proper database and Redis hosts (`database` and `cache`)
$ kool start
```

To get started in an existing Laravel project you only need to:

```console
$ cd my-laravel-project/
$ kool preset laravel
$ # make sure your `.env` points to the proper database and Redis hosts (`database` and `cache`)
$ kool start
```


- There you go! Now you have a PHP 7.4, Mysql and Redis environment. You are encouraged to take a look and make changes you see fit at `docker-compose.yml` or `kool.yml` to better adjust your project specifications.

- The steps above will create some configuration files in your project folder - all of which you should commit to your version control system.

## Documentation

You can check the documentation at **https://kool.dev/docs** or at [docs/](docs/).


## Frameworks Presets

To help getting you started we've built presets as a starting point for some popular stacks and frameworks. In case you miss one let us know in an issue or feel free to open up a PR for it!

Some or our current presets to get you started in no time:

- [Laravel](docs/2-Presets/Laravel.md)
- [Symfony](docs/2-Presets/Symfony.md)
- [Wordpress](docs/2-Presets/Wordpress.md)
- [Adonis](docs/2-Presets/Adonis.md)
- [NestJS](docs/2-Presets/NestJS.md)
- [NextJS](docs/2-Presets/NextJS.md)
- [NuxtJS](docs/2-Presets/NuxtJS.md)
- [Hugo](docs/2-Presets/Hugo.md)

## See it in action (DEMO)

You can check out a couple of sample commands in action at [asciinema.org/~kooldev](https://asciinema.org/~kooldev). We will be continuously uploading more samples.

## Community, Contributing and Support

You are most welcome to contribute and help in our missiong of making software development *kool* for everyone.

- [Issues](/issues) are the primary channel for bringing up and tracking issues or proposals.
- [Kool community on Slack](https://kool-dev.slack.com) is the a great place to get help and reach Kool developers.
- Check out our draft on [contributing guide](CONTRIBUTING.md) for getting involved.

### Roadmap

We have been working in a loosely defined but clear roadmap. You can check it out [in our blog Roadmap page](https://blog.kool.dev/page/roadmap).

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

If you find security issue please let us know ahead of making it public like in an issue so we can take action as soon as possible. Please email the concern to `contact@kool.dev`.

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.


[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkool-dev%2Fkool.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkool-dev%2Fkool?ref=badge_large)
