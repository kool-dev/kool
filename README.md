<p align="center"><a href="https://kool.dev" target="_blank"><img src="https://kool.dev/img/logo.png" width="400" alt="kool - cloud native dev tool"></a></p>


<p align="center">
<a href="https://goreportcard.com/report/github.com/kool-dev/kool"><img src="https://goreportcard.com/badge/github.com/kool-dev/kool" alt="Go Report Card"></a>
<a href="https://codecov.io/gh/kool-dev/kool"><img src="https://codecov.io/gh/kool-dev/kool/branch/main/graph/badge.svg" alt="codecov"></a>
<a href="https://github.com/kool-dev/kool/workflows/docker"><img src="https://github.com/kool-dev/kool/workflows/docker/badge.svg" alt="Docker Hub"></a>
<a href="https://github.com/kool-dev/kool/workflows/golangci-lint"><img src="https://github.com/kool-dev/kool/workflows/golangci-lint/badge.svg" alt="Golang CI Lint"></a>
<a href="https://codeclimate.com/github/kool-dev/kool/maintainability"><img src="https://api.codeclimate.com/v1/badges/1511f826de92d2ab39cc/maintainability" alt="Maintainability"></a>
<a href="https://kool.dev/slack"><img src="https://img.shields.io/badge/Join%20Slack-kool--dev-orange?logo=slack" alt="Join Slack Kool community"></a>
<a href="https://github.com/sindresorhus/awesome"><img src="https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg" alt="Awesome"></a>
</p>

## About `kool`

**Kool** is a CLI tool that brings the complexities of modern software development environments down to earth - making these environments lightweight, fast and reproducible. It reduces the complexity and learning curve of _Docker_ containers for local environments, and offers a simplified interface for using _Kubernetes_ to deploy staging and production environments to the cloud.

**Kool** gets your local development environment up and running easily and quickly, so you have more time to build a great application. When the time is right, you can then use Kool Cloud to deploy and share your work with the world!

**Kool** is suitable for solo developers and teams of all sizes. It provides a hassle-free way to handle the Docker basics and immediately start using containers for development, while simultaneously guaranteeing no loss of control over more specialized Docker environments.

[Learn more at kool.dev](https://kool.dev).

## Installation

Requirements: Kool is powered by [Docker](https://docs.docker.com/get-docker/). If you haven't done so already, you first need to [install Docker and the kool CLI](https://kool.dev/docs/getting-started/installation).

**Important**: make sure you are running the latest version of Docker and that you do have Compose V2 available (`docker compose`). You can read more about [Compose V2 release via its documentation](https://docs.docker.com/compose/reference/). Checkout out instructions for [installing Docker Compose V2 in the official documentation](https://docs.docker.com/compose/install/#scenario-two-install-the-compose-plugin).

### For Linux and MacOS

Install **kool** by running the following script. It will download the latest **kool** binary from [https://github.com/kool-dev/kool/releases](https://github.com/kool-dev/kool/releases), and save it in your `/usr/local/bin` folder.

```bash
curl -fsSL https://kool.dev/install | bash
```

### For Windows

You must run `kool` on Windows via [WSL - Windows Subsystem for Linux](https://learn.microsoft.com/en-us/windows/wsl/install) - once you have a WSL environment properly set up, make sure you have [Docker available on it](https://docs.docker.com/desktop/wsl/), then you can install the CLI as you would in any Linux or MacOS (see above).

## Getting Started

It's really easy to get started with `kool`. Check out our [Getting Started documentation for a generic PHP web app](https://kool.dev/docs/getting-started/starting-new-project).

To help you start building real-world applications, we've created Kool Presets as a starting point for some popular frameworks and stacks.

### Available Presets

#### Popular stacks and frameworks

- **Node**: [NestJS](docs/2-Presets/NestJS.md), [AdonisJs](docs/2-Presets/AdonisJs.md), [Express.js](/docs/2-Presets/ExpressJS.md)
- **PHP**: [Laravel](docs/2-Presets/Laravel.md), [Laravel Octane](docs/2-Presets/Laravel+Octane.md), [Symfony](docs/2-Presets/Symfony.md), [CodeIgniter](docs/2-Presets/CodeIgniter.md)
- **Javascript**: [Next.js](docs/2-Presets/NextJS.md), [NuxtJS](docs/2-Presets/NuxtJS.md)
- **Others**: [Hugo](docs/2-Presets/Hugo.md), [WordPress](docs/2-Presets/WordPress.md)

#### Monorepo structures

It's a common challange mixing up different frameworks for the frontned and a backend API. Working with containers and having them both working an communicating properly can be a huge differential for good development experience and productivity.

Check out our pre-shaped [mono-repo structures](https://monorepo.tools/#what-is-a-monorepo) in a single preset:

- [Monorepo NestJS + Next.js](docs/2-Presets/2-Monorepo-NestJS-with-NextJS.md) with Typescript on both the frontend and the backend.

> If you don't see your favorite framework in the list above, please let us know by creating a GitHub issue, or, better yet, feel free to submit a PR!

## Documentation

Read the documentation at [**https://kool.dev/docs**](https://kool.dev/docs) or [docs/](docs/).

## Community, Contributing and Support

We invite you to contribute and help in our mission of making software development *kool* for everyone.

- [Issues](/issues) are the primary channel for bringing issues and proposals to our attention.
- [Kool on Slack](https://kool.dev/slack) is a great place to get help and interact with Kool developers.
- Learn how to get involved by reading our [contributing guide](CONTRIBUTING.md).

## Roadmap

Our work is organized according to a loosely defined but clear roadmap. Check out [the Roadmap page](https://blog.kool.dev/page/roadmap) on [our blog](https://blog.kool.dev/).

## Security

If you find a security issue, please let us know right away, before making it public, by creating a GitHub issue. We'll take action as soon as possible. You can email questions and concerns to `contact@kool.dev`.

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
