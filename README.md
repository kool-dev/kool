<p align="center"><a href="https://kool.dev" target="_blank"><img src="https://kool.dev/img/logo.png" width="400" alt="kool - cloud native dev tool"></a></p>


<p align="center">
<a href="https://goreportcard.com/report/github.com/kool-dev/kool"><img src="https://goreportcard.com/badge/github.com/kool-dev/kool" alt="Go Report Card"></a>
<a href="https://codecov.io/gh/kool-dev/kool"><img src="https://codecov.io/gh/kool-dev/kool/branch/master/graph/badge.svg" alt="codecov"></a>
<a href="https://github.com/kool-dev/kool/workflows/docker"><img src="https://github.com/kool-dev/kool/workflows/docker/badge.svg" alt="Docker Hub"></a>
<a href="https://github.com/kool-dev/kool/workflows/golangci-lint"><img src="https://github.com/kool-dev/kool/workflows/golangci-lint/badge.svg" alt="Golang CI Lint"></a>
<a href="https://codeclimate.com/github/kool-dev/kool/maintainability"><img src="https://api.codeclimate.com/v1/badges/1511f826de92d2ab39cc/maintainability" alt="Maintainability"></a>
<a href="https://kool.dev/slack"><img src="https://img.shields.io/badge/Join%20Slack-kool--dev-orange?logo=slack" alt="Join Slack Kool community"></a>
</p>

## About `kool`

**kool** is a CLI tool that helps bringing down to earth the complexities of modern software development environments - making them lightweight, fast and reproduceable. It takes off the complexity and learning curve of _Docker_ and _Docker Compose_ for local environments, as well as offers a highly simplified interface for leveraging _Kubernetes_ cloud deployment for staging and production deployments.

Get your local development environment up and running easy and quickly, put time and effort on making a great application, and then leverage the Kool cloud for deploying and sharing your work with the world! This tool is suitable for single developers or large teams, powering them with a simple start and still provides all flexibility the DevOps team needs to tailor up everything.

To learn more [check out our website kool.dev](https://kool.dev).

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

It is easy to get started leveraging `kool`. Check out our [getting started documentation for a generic PHP web app](https://kool.dev/docs/getting-started/starting-new-project).

To help getting you started on real life applications, we've built presets as a starting point for some popular stacks and frameworks.

#### Available Presets

- [Laravel](docs/2-Presets/Laravel.md)
- [Symfony](docs/2-Presets/Symfony.md)
- [Wordpress](docs/2-Presets/Wordpress.md)
- [Adonis](docs/2-Presets/Adonis.md)
- [NestJS](docs/2-Presets/NestJS.md)
- [NextJS](docs/2-Presets/NextJS.md)
- [NuxtJS](docs/2-Presets/NuxtJS.md)
- [Hugo](docs/2-Presets/Hugo.md)

> In case you miss your favorite framework of choice from the list above, please let us know in an issue or feel free to open up a PR for it!

## Documentation

You can check the documentation at [**https://kool.dev/docs**](https://kool.dev/docs) or at [docs/](docs/).

## Community, Contributing and Support

You are most welcome to contribute and help in our mission of making software development *kool* for everyone.

- [Issues](/issues) are the primary channel for bringing up and tracking issues or proposals.
- [Kool community on Slack](https://kool.dev/slack) is the a great place to get help and reach Kool developers.
- Check out our [contributing guide](CONTRIBUTING.md) for getting involved.

### Roadmap

We have been working in a loosely defined but clear roadmap. You can check it out [in our blog Roadmap page](https://blog.kool.dev/page/roadmap).

## Security

If you find security issue please let us know ahead of making it public like in an issue so we can take action as soon as possible. Please email the concern to `contact@kool.dev`.

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
