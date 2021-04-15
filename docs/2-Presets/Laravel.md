## Start a New Laravel Project with Docker in 3 Easy Steps

1. Run `kool create laravel my-project`
2. Update **.env.example**
3. Run `kool run setup`

> Yes, using **kool** + Docker to create and work on new Laravel projects is that easy!

### Requirements

If you haven't done so already, you first need to [install Docker and the kool CLI](/docs/getting-started/installation).

> Please note that you'll need a basic understanding of how Docker and Docker Compose work in order to build a new project from scratch using Kool with Docker.

Also, make sure you're running the latest version of **kool**. Run the following command to compare your local version of **kool** with the latest release, and, if a newer version is available, automatically download and install it.

```bash
$ kool self-update
```

### 1. Run `kool create laravel my-project`

Use the [`kool create <preset> <project-name>` command](docs/commands/kool-create) to create your new Laravel project:

```bash
$ kool create laravel my-project
```

Under the hood, this command will run `composer create-project --no-install --no-scripts --prefer-dist laravel/laravel <project-name>` (using a customized **kool** Docker image: <a href="https://github.com/kool-dev/docker-php" target="_blank">kooldev/php:7.4</a>).

After installing Laravel, `kool create` automatically runs the `kool preset laravel` command, which helps you easily set up the initial tech stack for your project using an interactive wizard.

```bash
$ Preset laravel is initializing!

? What app service do you want to use [Use arrows to move, type to filter]
  PHP 7.4
> PHP 8.0

? What database service do you want to use [Use arrows to move, type to filter]
> MySQL 8.0
  MySQL 5.7
  PostgreSQL 13.0
  none

? What cache service do you want to use [Use arrows to move, type to filter]
> Redis 6.0
  Memcached 1.6
  none

? What javascript package manager do you want to use [Use arrows to move, type to filter]
  npm
> yarn

? What composer version do you want to use [Use arrows to move, type to filter]
  1.x
> 2.x

$ Preset laravel initialized!
```

Now, move into your new Laravel project:

```bash
$ cd my-project
```

The [`kool preset` command](docs/commands/kool-preset) auto-generated the following configuration files and added them to your project, which you can modify and extend.

```bash
+.dockerignore
+docker-compose.yml
+kool.yml
```

> Now's a good time to open the **docker-compose.yml** file to review the services and verify they match the choices you made earlier using the wizard.

### 2. Update .env.example

You need to update Laravel's **.env.example** file with some default values that match the services in your **docker-compose.yml** file.

#### Database Services

MySQL 5.7 and 8.0

```diff
-DB_HOST=127.0.0.1
+DB_HOST=database
```

PostgreSQL 13.0

```diff
-DB_CONNECTION=mysql
+DB_CONNECTION=pgsql

-DB_HOST=127.0.0.1
+DB_HOST=database

-DB_PORT=3306
+DB_PORT=5432
```

#### Cache Services

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

### 3. Run `kool run setup`

> Say hello to **kool.yml**, say goodbye to custom shell scripts!

As mentioned above, the [`kool preset` command](docs/commands/kool-preset) added a **kool.yml** file to your project. Think of **kool.yml** as a super easy-to-use task _helper_. Instead of writing custom shell scripts, add your own scripts to **kool.yml** (under the `scripts` key), and run them with `kool run <script>` (e.g. `kool run artisan`). You can add your own single line commands (see `composer` below), or add a list of commands that will be executed in sequence (see `setup` below).

To help get you started, **kool.yml** comes prebuilt with an initial set of scripts (based on the choices you made earlier using the **preset** wizard).

```yaml
scripts:
  artisan: kool exec app php artisan
  composer: kool exec app composer2
  mysql: kool exec -e MYSQL_PWD=$DB_PASSWORD database mysql -uroot
  node: kool docker kooldev/node:14 node

  node-setup:
    - kool run yarn install
    - kool run yarn dev

  reset:
    - kool run composer install
    - kool run artisan migrate:fresh --seed
    - kool run node-setup

  setup:
    - cp .env.example .env
    - kool start
    - kool run composer install
    - kool run artisan key:generate
    - kool run node-setup

  yarn: kool docker kooldev/node:14 yarn
```

> Try `kool run artisan --help` to execute the `kool exec app php artisan --help` command in your running `app` container and print out information about Laravel's CLI commands.

**kool.yml** always includes a script called `setup`, which helps you spin up a project for the first time.

Go ahead and run `kool run setup` to start your Docker environment and finish setting up your project:

```bash
# CAUTION: this script will reset your `.env` file with `.env.example`
$ kool run setup
```

> As you can see in **kool.yml**, the `setup` script will do the following in sequence: copy your updated **.env.example** file to **.env**; start your Docker environment; use Composer to install vendor dependencies; generate your `APP_KEY` (in `.env`); and then build your Node packages and assets.

Once `kool run setup` finishes, you should be able to access your site at [http://localhost](http://localhost).

---

#### Run a Container Command

Use [`kool exec`](/docs/commands/kool-exec) to execute a command inside a running service container:

```bash
# kool exec [OPTIONS] SERVICE COMMAND [--] [ARG...]

$ kool exec app ls
```

#### Open a Container Session

Similar to SSH, if you want to open a Bash session in your `app` container, run `kool exec app bash`, where `app` is the name of the service container in **docker-compose.yml**. If you prefer, you can use `sh` instead of `bash` (`kool exec app sh`).

```bash
$ kool exec app bash
bash-5.1#

$ kool exec app sh
/app #
```

#### Connect to your Database Container

You can easily start a new SQL client session inside your Docker `database` container by executing `kool run mysql` (MySQL) or `kool run psql` (PostgreSQL) in your terminal. This runs the single-line `mysql` or `psql` script included in your **kool.yml**.

### Staying kool

When it's time to stop working on the project:

```bash
$ kool stop
```

And when you're ready to start work again:

```bash
$ kool start
```
