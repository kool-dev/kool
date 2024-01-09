# Start a Laravel Octane Project with Docker in 3 Easy Steps

1. Run `kool create laravel+octane my-project`
2. Update **.env.example**
3. Run `kool run setup`

> Yes, using **kool** + Docker to create and work on new Laravel Octane projects is that easy!

## Requirements

If you haven't done so already, you first need to [install Docker and the kool CLI](/docs/getting-started/installation).

> Also, make sure you're running the latest version of **kool**. Run the following command: `kool self-update`

> Please note that it helps to have a basic understanding of how Docker and Docker Compose work to use Kool with Docker.

## 1. Run `kool create laravel+octane my-project`

Use the [`kool create PRESET FOLDER` command](/docs/commands/kool-create) to create your new Laravel Octane project:

```bash
$ kool create laravel+octane my-project
```

This command will guide you through setting up a new Laravel project, installing Laravel Octane with your preferred server engine (either Swoole or RoadRunner), and setting up all the Docker Compose configuration files to manage your dockerized development environment.

After that, you will have the option to include a database or cache service, all of which helps you easily set up the initial tech stack for your project using an interactive wizard.

---

Now, move into your new Laravel Octane project:

```bash
$ cd my-project
```

The [`kool preset` command](/docs/commands/kool-preset) auto-generated the following configuration files and added them to your project, which you can modify and extend.

```bash
+docker-compose.yml
+kool.yml
```

> Now's a good time to review the **docker-compose.yml** file and verify the services match the choices you made earlier using the wizard.

## 2. Update .env.example

You need to update some default values in Laravel's **.env.example** file to match the services in your **docker-compose.yml** file.

### Database Services

MySQL or MariaDB

```diff
-DB_HOST=127.0.0.1
+DB_HOST=database
```

PostgreSQL

```diff
-DB_CONNECTION=mysql
+DB_CONNECTION=pgsql

-DB_HOST=127.0.0.1
+DB_HOST=database

-DB_PORT=3306
+DB_PORT=5432
```

> In order to avoid permission issues with mysql and MariaDB, add a user other than root and a password to your **.env.example** file

```diff
-DB_USERNAME=root
+DB_USERNAME=<some_user>

-DB_PASSWORD=
+DB_PASSWORD=<somepass>
```

### Cache Services

Redis

```diff
-REDIS_HOST=127.0.0.1
+REDIS_HOST=cache
```

Memcached

```diff
-MEMCACHED_HOST=127.0.0.1
+MEMCACHED_HOST=cache
```

## 3. Run `kool run setup`

Go ahead and run `kool run setup` to start your Docker environment and finish setting up your project:

```bash
# CAUTION: this script will reset your `.env` file with `.env.example`
$ kool run setup
```

### About `kool.yml` and `kool run setup`

> Say hello to **kool.yml**, say goodbye to custom shell scripts!

As mentioned above, the [`kool preset` command](/docs/commands/kool-preset) added a **kool.yml** file to your project. Think of **kool.yml** as a super easy-to-use task _helper_. Instead of writing custom shell scripts, add your own scripts to **kool.yml** (under the `scripts` key), and run them with `kool run SCRIPT` (e.g. `kool run artisan`). You can add your own single line commands (see `composer` below), or add a list of commands that will be executed in sequence (see `setup` below).

To help get you started, **kool.yml** comes prebuilt with an initial set of scripts (based on the choices you made earlier using the **preset** wizard), including a script called `setup`, which helps you spin up a project for the first time.

Go ahead and run `kool run setup` to start your Docker environment and finish setting up your project:

```bash
# CAUTION: this script will reset your `.env` file with `.env.example`
$ kool run setup
```

> As you can see in **kool.yml**, the `setup` script will do the following in sequence: copy your updated **.env.example** file to **.env**; start your Docker environment; use Composer to install vendor dependencies; generate your `APP_KEY` (in `.env`); and then build your Node packages and assets.

Once `kool run setup` finishes, you should be able to access your new site at [http://localhost](http://localhost) and see the Laravel welcome page. Hooray!

Verify your Docker container is running using the [`kool status` command](/docs/commands/kool-status).

Run `kool logs app` to see the logs from the running `app` container.

> Use `kool logs` to see the logs from all running containers. Add the `-f` option after `kool logs` to follow the logs (i.e. `kool logs -f app`).

---

### Run Commands in Docker Containers

Use [`kool exec`](/docs/commands/kool-exec) to execute a command inside a running service container:

```bash
# kool exec [OPTIONS] SERVICE COMMAND [--] [ARG...]

$ kool exec app ls
```

Try `kool run artisan --help` to execute the `kool exec app php artisan --help` command in your running `app` container and print out information about Laravel's CLI commands.

### Open Sessions in Docker Containers

Similar to SSH, if you want to open a Bash session in your `app` container, run `kool exec app bash`, where `app` is the name of the service container in **docker-compose.yml**. If you prefer, you can use `sh` instead of `bash` (`kool exec app sh`).

```bash
$ kool exec app bash
bash-5.1#

$ kool exec app sh
/app #
```

### Connect to Docker Database Container

You can easily start a new SQL client session inside your running `database` container by executing `kool run mysql` (MySQL) or `kool run psql` (PostgreSQL) in your terminal. This runs the single-line `mysql` or `psql` script included in your **kool.yml**.

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

- **[AdonisJs](/docs/03-Presets/AdonisJs.md)**
- **[CodeIgniter](/docs/03-Presets/CodeIgniter.md)**
- **[Express.js](/docs/03-Presets/ExpressJS.md)**
- **[Hugo](/docs/03-Presets/Hugo.md)**
- **[NestJS](/docs/03-Presets/NestJS.md)**
- **[Next.js](/docs/03-Presets/NextJS.md)**
- **[Node.js](/docs/03-Presets/NodeJS.md)**
- **[Nuxt.js](/docs/03-Presets/NuxtJS.md)**
- **[PHP](/docs/03-Presets/PHP.md)**
- **[Symfony](/docs/03-Presets/Symfony.md)**
- **[WordPress](/docs/03-Presets/WordPress.md)**

Missing a preset? **[Make a request](https://github.com/kool-dev/kool/issues/new)**, or contribute by opening a Pull Request. Go to [https://github.com/kool-dev/kool/tree/main/presets](https://github.com/kool-dev/kool/tree/main/presets) and browse the code to learn more about how presets work.
