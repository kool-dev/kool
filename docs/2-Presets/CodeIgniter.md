# Start a CodeIgniter Project with Docker in 3 Easy Steps

1. Run `kool create codeigniter my-project`
2. Update **env**
3. Run `kool run setup`

> Yes, using **kool** + Docker to create and work on new CodeIgniter projects is that easy!

## Requirements

If you haven't done so already, you first need to [install Docker and the kool CLI](/docs/getting-started/installation).

Also, make sure you're running the latest version of **kool**. Run the following command to compare your local version of **kool** with the latest release, and, if a newer version is available, automatically download and install it.

```bash
$ kool self-update
```

> Please note that it helps to have a basic understanding of how Docker and Docker Compose work to use Kool with Docker.

## 1. Run `kool create codeigniter my-project`

Use the [`kool create PRESET FOLDER` command](/docs/commands/kool-create) to create your new CodeIgniter project:

```bash
$ kool create codeigniter my-project
```

Under the hood, this command will run `composer create-project --no-install --no-scripts --prefer-dist codeigniter4/appstarter my-project` using a customized **kool** Docker image: <a href="https://github.com/kool-dev/docker-php" target="_blank">kooldev/php:7.4</a>.

After installing CodeIgniter, `kool create` automatically runs the `kool preset codeigniter` command, which helps you easily set up the initial tech stack for your project using an interactive wizard.

```bash
$ Preset CodeIgniter is initializing!

? Which app service do you want to use [Use arrows to move, type to filter]
> PHP 7.4
  PHP 8.0

? Which database service do you want to use [Use arrows to move, type to filter]
> MySQL 8.0
  MySQL 5.7
  MariaDB 10.5
  PostgreSQL 13.0
  none

? Which cache service do you want to use [Use arrows to move, type to filter]
> Redis 6.0
  Memcached 1.6
  none

$ Preset codeigniter initialized!
```

Now, move into your new CodeIgniter project:

```bash
$ cd my-project
```

The [`kool preset` command](/docs/commands/kool-preset) auto-generated the following configuration files and added them to your project, which you can modify and extend.

```bash
+docker-compose.yml
+kool.yml
```

> Now's a good time to review the **docker-compose.yml** file and verify the services match the choices you made earlier using the wizard.

## 2. Update env

You need to update some default values in CodeIgniter's **env** file to match the services in your **docker-compose.yml** file.

### Database Services

MySQL 5.7 and 8.0 or MariaDB 10.5

```diff
-database.default.hostname = localhost
+database.default.hostname = database

-database.default.database = ci4
+DB_DATABASE = ci4
+database.default.database = "${DB_DATABASE}"

+database.default.port = 3306
```

PostgreSQL 13.0

```diff
-database.default.DBDriver = MySQLi
+database.default.DBDriver = Postgre

-database.default.hostname = localhost
+database.default.hostname = database

-database.default.database = ci4
+DB_DATABASE = ci4
+database.default.database = "${DB_DATABASE}"

+database.default.port = 5432
```

> In order to avoid permission issues with mysql and mariaDB, add a user other than root and a password to your **env** file

```diff
-database.default.username = root
+DB_USERNAME = <some_user>
+database.default.username = "${DB_USERNAME}"

-database.default.password = root
+DB_PASSWORD = <somepass>
+database.default.password = "${DB_PASSWORD}"
```

## 3. Update app/Config/Cache.php

You need to update some default values in CodeIgniter's **app/Config/Cache.php** file to match the services in your **docker-compose.yml** file.

### Cache Services

Redis

```diff
public $redis = [
-    'host' => '127.0.0.1',
+    'host' => cache
```

Memcached

```diff
public $memcached = [
-    'host' => '127.0.0.1',
+    'host' => cache
```

## 4. Run `kool run setup`

> Say hello to **kool.yml**, say goodbye to custom shell scripts!

As mentioned above, the [`kool preset` command](/docs/commands/kool-preset) added a **kool.yml** file to your project. Think of **kool.yml** as a super easy-to-use task _helper_. Instead of writing custom shell scripts, add your own scripts to **kool.yml** (under the `scripts` key), and run them with `kool run SCRIPT` (e.g. `kool run spark`). You can add your own single line commands (see `composer` below), or add a list of commands that will be executed in sequence (see `setup` below).

To help get you started, **kool.yml** comes prebuilt with an initial set of scripts (based on the choices you made earlier using the **preset** wizard), including a script called `setup`, which helps you spin up a project for the first time.

```yaml
scripts:
  composer: kool exec app composer
  spark: kool exec app php spark
  phpunit: kool exec app php ./bin/phpunit

  setup:
    - cp env .env
    - kool start
    - kool run composer install
    - kool run spark key:generate

  reset:
    - kool run composer install
    - kool run spark migrate:refresh -f
```

Go ahead and run `kool run setup` to start your Docker environment and finish setting up your project:

```bash
# CAUTION: this script will reset your `.env` file with `env`
$ kool run setup
```

> As you can see in **kool.yml**, the `setup` script will do the following in sequence: copy your updated **env** file to **.env**; start your Docker environment; use Composer to install vendor dependencies; generate your `APP_KEY` (in `.env`); and then build your Node packages and assets.

Once `kool run setup` finishes, you should be able to access your new site at [http://localhost](http://localhost) and see the CodeIgniter welcome page. Hooray!

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

Try `kool run spark --help` to execute the `kool exec app php spark --help` command in your running `app` container and print out information about CodeIgniter's CLI commands.

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

- **[AdonisJs](/docs/2-Presets/AdonisJs.md)**
- **[Hugo](/docs/2-Presets/Hugo.md)**
- **[Laravel](/docs/2-Presets/Laravel.md)**
- **[NestJS](/docs/2-Presets/NestJS.md)**
- **[Next.js](/docs/2-Presets/NextJS.md)**
- **[Node.js](/docs/2-Presets/NodeJS.md)**
- **[Nuxt.js](/docs/2-Presets/NuxtJS.md)**
- **[PHP](/docs/2-Presets/PHP.md)**
- **[Symfony](/docs/2-Presets/Symfony.md)**
- **[WordPress](/docs/2-Presets/WordPress.md)**

Missing a preset? **[Make a request](https://github.com/kool-dev/kool/issues/new)**, or contribute by opening a Pull Request. Go to [https://github.com/kool-dev/kool/tree/master/presets](https://github.com/kool-dev/kool/tree/master/presets) and browse the code to learn more about how presets work.
