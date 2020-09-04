Kool works with Docker / Docker Compose under the hood, and comes with some cool presets to help you get started, everything is configurable / extendable.

Let's use **Laravel** preset as example and explain how it works.

When you run **kool init laravel** all it does is create a few files for you:

```bash
$ kool init laravel
Preset laravel is initializing!
  Preset file Dockerfile.build created.
  Preset file docker-compose.yml created.
  Preset file kool.yml created.
Preset  laravel  initialized!
```

### Dockerfile.build

This is a file you can use in case you want to build your application to use in production, in Docker world every release is usually a new image built with your application on it.

Soon we will give more examples on how to use Docker in production or use it with **Kool Cloud**.

### docker-compose.yml

This file defines all services that runs your application, docker images to use, ports, volume mounts, etc.

You can add/change/remove services as you will.

### kool.yml

This is where most of the magic happens, a way to make your life easy, orchestrating scripts for you to use on your local environment or CI/CDs. It is created in your working directory when you run **kool init**, but you can also create it inside a folder named **kool** in your user's home directory.

The **scripts** defined will be used by **kool run** command, for example:

```yaml
scripts:
  artisan: kool exec app php artisan

  setup:
    - kool start
    - cp .env.example .env
    - kool run artisan key:generate
```

Single commands like **artisan** are kind of aliases, so anything you input will be forwarded to the actual command, so if you run: **kool run artisan key:generate** it will basically translate into: **kool exec app php artisan key:generate**.

Multiple commands like **setup** will not forward your input, so **kool run setup something** will run every script and **something** will be ignored.

PS: This is not only limited to kool commands, so you can add any type command the **cp** in example.
