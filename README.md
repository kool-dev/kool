# kool

Dev environment made easy, a standarized way for running applications no matter the stack on your local machine and deploying it to a development environment.

Have the same feeling working on multiple projects with different stacks.

## Installation

One Liner

## Usage

To help learning how to use kool we've built presets with good starting point for some popular stacks, feel free to open a PR in case you miss one.

### Presets

- [Adonis](docs/presets/adonis.md)
- [Adonis with NextJS](docs/presets/adonis-nextjs.md)
- [Laravel](docs/presets/laravel.md)
- [NextJS](docs/presets/nextjs.md)
- [NextJS Static](docs/presets/nextjs-static.md)
- [NuxtJS](docs/presets/nuxtjs.md)
- [NuxtJS Static](docs/presets/nuxtjs-static.md)

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

### kool install

```bash
$ kool install [preset] [flags]
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
$ kool run [script/image] [command]
```

Execute script or a docker image.

| Name | Type | Description |
| ---- | ---- | ----------- |
| `script/image` | `string` | Script to run within your `kool.yaml` file.|
| `command` | `string` | The command to run, i.e: `composer install`, `yarn install` |
| `--docker` | `none` | If enabled, `script` param will become `image` and you will run a docker image |

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
