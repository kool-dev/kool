# Start a PHP Project with Docker in 2 Easy Steps

1. Run `kool create php my-project`
2. Run `kool run setup`

> Yes, using **kool** + Docker to create and work on new PHP projects is that easy!

## Requirements

If you haven't done so already, you first need to [install Docker and the kool CLI](/docs/getting-started/installation).

Also, make sure you're running the latest version of **kool**. Run the following command to compare your local version of **kool** with the latest release, and, if a newer version is available, automatically download and install it.

```bash
$ kool self-update
```

> Please note that it helps to have a basic understanding of how Docker and Docker Compose work to use Kool with Docker.

## 1. Run `kool create php my-project`

Use the [`kool create PRESET FOLDER` command](/docs/commands/kool-create) to create your new PHP project:

```bash
$ kool create php my-project
```

Under the hood, this command will create a "Hello world!" **index.php** file inside a `public` directory, which is the default root used by NGINX.

After setting up the project, `kool create` automatically runs the `kool preset php` command, which helps you easily set up the initial tech stack for your project using an interactive wizard.

```bash
$ Preset php is initializing!

? What app service do you want to use [Use arrows to move, type to filter]
> PHP 7.4
  PHP 8.0

? What composer version do you want to use [Use arrows to move, type to filter]
> 1.x
  2.x

$ Preset php initialized!
```

Now, move into your new PHP project:

```bash
$ cd my-project
```

The [`kool preset` command](/docs/commands/kool-preset) auto-generated the following configuration files and added them to your project, which you can modify and extend.

```bash
+docker-compose.yml
+kool.yml
```

> Now's a good time to review the **docker-compose.yml** file and verify the services match the choices you made earlier using the wizard.

## 2. Run `kool run setup`

> Say hello to **kool.yml**, say goodbye to custom shell scripts!

As mentioned above, the [`kool preset` command](/docs/commands/kool-preset) added a **kool.yml** file to your project. Think of **kool.yml** as a super easy-to-use task _helper_. Instead of writing custom shell scripts, add your own scripts to **kool.yml** (under the `scripts` key), and run them with `kool run SCRIPT` (e.g. `kool run composer`). You can add your own single line commands (see `composer` below), or add a list of commands that will be executed in sequence (see `setup` below).

To help get you started, **kool.yml** comes prebuilt with an initial set of scripts (based on the choices you made earlier using the **preset** wizard), including a script called `setup`, which helps you spin up a project for the first time.

```yaml
scripts:
  composer: kool exec app composer
  php: kool exec app php

  setup:
    - kool start
    # - add more setup commands
```

Go ahead and run `kool run setup` to start your Docker environment:

```bash
# CAUTION: this script will reset your `.env` file with `.env.example`
$ kool run setup
```

> As you can see in **kool.yml**, the `setup` script will do the following in sequence: run the `kool start` command to spin up your Docker environment; and then run any additional commands you add to the list.

Once `kool run setup` finishes, you should be able to access your new site at [http://localhost](http://localhost) and see "Hello world!". Hooray!

Verify your Docker container is running using the [`kool status` command](/docs/commands/kool-status).

```bash
$ kool status
+---------+---------+---------------------------------------------+--------------+
| SERVICE | RUNNING | PORTS                                       | STATE        |
+---------+---------+---------------------------------------------+--------------+
| app     | Running | 0.0.0.0:80->80/tcp, :::80->80/tcp, 9000/tcp | Up 3 seconds |
+---------+---------+---------------------------------------------+--------------+
```

Run `kool logs app` to see the logs from the running `app` container.

> Use `kool logs` to see the logs from all running containers. Add the `-f` option after `kool logs` to follow the logs (i.e. `kool logs -f app`).

```bash
$ kool logs app
Attaching to my-project_app_1
app_1  | time="2021-05-01T16:06:07Z" level=info msg="create process:php-fpm"
app_1  | time="2021-05-01T16:06:07Z" level=info msg="create process:nginx"
app_1  | time="2021-05-01T16:06:07Z" level=info msg="stop listening"
app_1  | time="2021-05-01T16:06:07Z" level=info msg="try to start program" program=php-fpm
app_1  | time="2021-05-01T16:06:07Z" level=info msg="try to start program" program=nginx
app_1  | time="2021-05-01T16:06:07Z" level=debug msg="wait program exit" program=php-fpm
app_1  | time="2021-05-01T16:06:07Z" level=debug msg="wait program exit" program=nginx
app_1  | [01-May-2021 16:06:07] NOTICE: fpm is running, pid 24
app_1  | [01-May-2021 16:06:07] NOTICE: ready to handle connections
app_1  | time="2021-05-01T16:06:08Z" level=info msg="success to start program" program=php-fpm
app_1  | time="2021-05-01T16:06:08Z" level=info msg="success to start program" program=nginx
```

---

### Run Commands in Docker Containers

Use [`kool exec`](/docs/commands/kool-exec) to execute a command inside a running service container:

```bash
# kool exec [OPTIONS] SERVICE COMMAND [--] [ARG...]

$ kool exec app ls
```

Try `kool run composer list` to execute the `kool exec app composer list` command in your running `app` container and print out a list of Composer commands.

### Open Sessions in Docker Containers

Similar to SSH, if you want to open a Bash session in your `app` container, run `kool exec app bash`, where `app` is the name of the service container in **docker-compose.yml**. If you prefer, you can use `sh` instead of `bash` (`kool exec app sh`).

```bash
$ kool exec app bash
bash-5.1#

$ kool exec app sh
/app #
```

### Access Private Repos and Packages in Docker Containers

If you need your `app` container to use your local SSH keys to pull private repositories and/or install private packages (which have been added as dependencies in your `composer.json` or `package.json` file), you can simply add `$HOME/.ssh:/home/kool/.ssh:delegated` under the `volumes` key of the `app` service in your **docker-compose.yml** file. This maps a `.ssh` folder in the container to the `.ssh` folder on your host machine.

```diff
volumes:
  - .:/app:delegated
+ - $HOME/.ssh:/home/kool/.ssh:delegated
```

## Staying kool

When it's time to stop working on the project:

```bash
$ kool stop
```

And when you're ready to start work again:

```bash
$ kool start
```

## Additional Presets

We have more presets to help you start projects with **kool** in a standardized way across different frameworks.

- **[AdonisJs](/docs/2-Presets/AdonisJs.md)**
- **[Hugo](/docs/2-Presets/Hugo.md)**
- **[Laravel](/docs/2-Presets/Laravel.md)**
- **[NestJS](/docs/2-Presets/NestJS.md)**
- **[Next.js](/docs/2-Presets/NextJS.md)**
- **[Nuxt.js](/docs/2-Presets/NuxtJS.md)**
- **[Symfony](/docs/2-Presets/Symfony.md)**
- **[WordPress](/docs/2-Presets/WordPress.md)**

Missing a preset? **[Make a request](https://github.com/kool-dev/kool/issues/new)**, or contribute by opening a Pull Request. Go to [https://github.com/kool-dev/kool/tree/master/presets](https://github.com/kool-dev/kool/tree/master/presets) and browse the code to learn more about how presets work.
