# kool

Dev environment made easy, a standarized way for running applications no matter the stack on your local machine and deploying it to a development environment.

Run any stack / tool with any version, powered by Docker and Docker Compose in a simple way avoiding you to install lots of stuff on your machine.

Have the same feeling working on multiple projects with different stacks.

## Demo

<a href="https://www.youtube.com/watch?v=14KNzJI8Oh0" target="_blank" title="Click to see full demo">
    <img src="https://user-images.githubusercontent.com/347400/87970968-fad10c80-ca9a-11ea-9bef-a88400b01f2c.png" alt="kool - demo" style="max-width:100%;">
</a>

## Installation

Kool is powered by [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/), you need to have it installed on your machine.

```bash
curl -fsSL https://raw.githubusercontent.com/kool-dev/kool/master/install.sh | bash
```
In case you need sudo:

```bash
curl -fsSL https://raw.githubusercontent.com/kool-dev/kool/master/install.sh | sudo bash
```

## Usage

To help learning how to use kool we've built presets with good starting point for some popular stacks, feel free to open a PR in case you miss one.

### Presets

- [Adonis](docs/presets/adonis.md)
- [Laravel](docs/presets/laravel.md)
- [NextJS](docs/presets/nextjs.md)
- [NuxtJS](docs/presets/nuxtjs.md)

### Examples

You can see projects using it here: https://github.com/kool-dev/examples

## Commands

### kool start

```bash
$ kool start [flags]
```

Start services (containers) defined on docker-compose.yml

| Name | Type | Description |
| ---- | ---- | ----------- |
| `--services=` | `string` | Specific services to be started |

### kool status

```bash
$ kool status
```

Shows the status for containers

### kool info

```bash
$ kool info
```

Prints out information about kool setup (like environment variables)

### kool init

```bash
$ kool init [preset] [flags]
```

Enable Kool preset configuration in the current working directory

| Name | Type | Description |
| ---- | ---- | ----------- |
| `preset` | `string` | The preset [(Presets)](#presets) |
| `--override` | `none` | Force replace local existing files with the default preset files |

### kool exec

```bash
$ kool exec [service] [command]
```

Execute command in running container.

| Name | Type | Description |
| ---- | ---- | ----------- |
| `service` | `string` | The service from `docker-compose.yml`, i.e: `app`,`database`,`adonis` |
| `command` | `string` | The command to run, i.e: `php artisan migrate`, `adonis run:migration`, `npm build` |

### kool run

```bash
$ kool run [script] [command]
```

Execute script or a docker image.

| Name | Type | Description |
| ---- | ---- | ----------- |
| `script` | `string` | Script to run within your `kool.yaml` file. |
| `command` | `string` | The command to run, i.e: `composer install`, `yarn install` |

### kool docker

```bash
$ kool docker [image] [command]
```

Execute script or a docker image.

| Name | Type | Description |
| ---- | ---- | ----------- |
| `image` | `string` | Docker image to run, i.e: `kooldev/node:14` |
| `command` | `string` | The command to run, i.e: `composer install`, `yarn install` |
| `--disable-tty / -T` | `none` | Force disable tty, good for CI/CI/Automation |

### kool stop

```bash
$ kool stop [flags]
```

Stop containers.

| Name | Type | Description |
| ---- | ---- | ----------- |
| `--purge` | `none` | If enabled, docker volume will be deleted. |

### Understanding kool.yml

This is where most of the magic happens, a way to make your life easy, orchestrating scripts for you to use on your local environment or CI/CDs. Look at presets to see examples.

### Understanding docker-compose.yml

This is where you control your local environment.

### Understanding Dockerfile.build

This file gives you a way for building Docker images for production. More docs to come.

## Contributing

[Build](docs/build.md)

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
