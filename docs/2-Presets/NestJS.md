# Start a NestJS Project with Docker in 2 Easy Steps

1. Run `kool create nestjs my-project`
2. Run `cd my-project && kool run setup`

> Yes, using **kool** + Docker to create and work on new NestJS projects is that easy!

## Requirements

If you haven't done so already, you first need to [install Docker and the kool CLI](/docs/getting-started/installation).

Also, make sure you're running the latest version of **kool**. Run the following command to compare your local version of **kool** with the latest release, and, if a newer version is available, automatically download and install it.

```bash
$ kool self-update
```

> Please note that it helps to have a basic understanding of how Docker and Docker Compose work to use Kool with Docker.

## 1. Run `kool create nestjs my-project`

Use the [`kool create PRESET FOLDER` command](/docs/commands/kool-create) to create your new NestJS project:

```bash
$ kool create nestjs my-project
```

Under the hood, this command will run `nest new my-project` to install NestJS with Typescript and NPM as the package manager.

After installing NestJS, `kool create` automatically runs the `kool preset nestjs` command, which helps you easily set up the initial tech stack for your project using an interactive wizard.

```bash
$ Preset nestjs is initializing!

? Which database service do you want to use [Use arrows to move, type to filter]
> MySQL 8.0
  MySQL 5.7
  PostgreSQL 13.0
  none

? Which cache service do you want to use [Use arrows to move, type to filter]
> Redis 6.0
  Memcached 1.6
  none

$ Preset nestjs initialized!
```

Now, move into your new NestJS project:

```bash
$ cd my-project
```

The [`kool preset` command](/docs/commands/kool-preset) auto-generated the following configuration files and added them to your project, which you can modify and extend.

```bash
+docker-compose.yml
+kool.yml
+.env.dist
```

> Now's a good time to review the **docker-compose.yml** file and verify the services match the choices you made earlier using the wizard.

## 2. Run `kool run setup`

> Say hello to **kool.yml**, say goodbye to custom shell scripts!

As mentioned above, the [`kool preset` command](/docs/commands/kool-preset) added a **kool.yml** file to your project. Think of **kool.yml** as a super easy-to-use task _helper_. Instead of writing custom shell scripts, add your own scripts to **kool.yml** (under the `scripts` key), and run them with `kool run SCRIPT` (e.g. `kool run nest`). You can add your own single line commands (see `nest` below), or add a list of commands that will be executed in sequence (see `setup` below).

To help get you started, **kool.yml** comes prebuilt with an initial set of scripts (based on the **preset**), including a script called `setup`, which helps you spin up a project for the first time.

```yaml
scripts:
  setup:
    # copy .env file
    - cp .env.dist .env
    # install backend deps
    - kool docker kooldev/node:16 npm install

  # helpers
  npm: kool exec app npm
  npx: kool exec app npx
  nest: kool run npx @nestjs/cli
```

Go ahead and run `kool run setup` to start your Docker environment and finish setting up your project:

```bash
$ kool run setup
$ kool start
```

> As you can see in **kool.yml**, the `setup` script will do the following in sequence: run `npm install` to build your Node packages and dependencies (by spinning up and down a temporary Node container). After that youo can start your Docker environment using **docker-compose.yml**  with the command `kool start` (which includes a container to running `npm run start:dev`).

Once `kool start` finishes, you should be able to access your new site at [http://localhost:3000](http://localhost:3000) and see the NestJS "Hello World!" welcome page.

Verify your Docker container is running using the [`kool status` command](/docs/commands/kool-status):

```bash
$ kool status
+---------+---------+------------------------+--------------+
| SERVICE | RUNNING | PORTS                  | STATE        |
+---------+---------+------------------------+--------------+
| app     | Running | 0.0.0.0:3000->3000/tcp | Up 5 seconds |
+---------+---------+------------------------+--------------+
```

Run `kool logs app` to see the logs from the running `app` container, and confirm the NestJS server was started.

> Use `kool logs` to see the logs from all running containers. Add the `-f` option after `kool logs` to follow the logs (i.e. `kool logs -f app`).

```bash
$ kool logs app
[2:55:57 AM] Starting compilation in watch mode...
my-project-app-1  |
my-project-app-1  | [2:56:11 AM] Found 0 errors. Watching for file changes.
my-project-app-1  |
my-project-app-1  | [Nest] 32  - 08/11/2022, 2:56:14 AM     LOG [NestFactory] Starting Nest application...
my-project-app-1  | [Nest] 32  - 08/11/2022, 2:56:14 AM     LOG [InstanceLoader] AppModule dependencies initialized +85ms
my-project-app-1  | [Nest] 32  - 08/11/2022, 2:56:14 AM     LOG [RoutesResolver] AppController {/}: +25ms
my-project-app-1  | [Nest] 32  - 08/11/2022, 2:56:14 AM     LOG [RouterExplorer] Mapped {/, GET} route +8ms
my-project-app-1  | [Nest] 32  - 08/11/2022, 2:56:14 AM     LOG [NestApplication] Nest application successfully started +6ms
```

---

### NestJS Configuration

If you added a database and/or cache service when answering the **preset** wizard questions, you'll need to add some local environment variables to match the services in your **docker-compose.yml** file (see below). To set these variables, it's common to use a **.env** file in your project root directory. Learn more about [how to configure NestJS](https://docs.nestjs.com/techniques/configuration).

#### Database Services

MySQL 5.7 and 8.0

```diff
+DB_CONNECTION=mysql
+DB_HOST=database
```

PostgreSQL 13.0

```diff
+DB_CONNECTION=pgsql
+DB_HOST=database
+DB_PORT=5432
```

#### Cache Services

Redis

```diff
+REDIS_HOST=cache
+REDIS_PORT=6379
```

Memcached

```diff
+MEMCACHED_HOST=cache
+MEMCACHED_PORT=11211
```

### Run Commands in Docker Containers

Use [`kool exec`](/docs/commands/kool-exec) to execute a command inside a running service container:

```bash
# kool exec [OPTIONS] SERVICE COMMAND [--] [ARG...]

$ kool exec app node -v
```

Try `kool run nest --help` to execute the `kool exec app nest --help` command in your running `app` container and print out information about NestJS' commands.

### Open Sessions in Docker Containers

Similar to SSH, if you want to open a Bash session in your `app` container, run `kool exec app bash`, where `app` is the name of the service container in **docker-compose.yml**. If you prefer, you can use `sh` instead of `bash` (`kool exec app sh`).

```bash
$ kool exec app bash
bash-5.1#

$ kool exec app sh
/app #
```

### Connect to Docker Database Container

If you added a database service, you can easily start a new SQL client session inside your running `database` container by executing `kool run mysql` (MySQL) or `kool run psql` (PostgreSQL) in your terminal. This runs the single-line `mysql` or `psql` script included in your **kool.yml**.

### Access Private Repos and Packages in Docker Containers

If you need your `app` container to use your local SSH keys to pull private repositories and/or install private packages (which have been added as dependencies in your `package.json` file), you can simply add `$HOME/.ssh:/home/kool/.ssh:delegated` under the `volumes` key of the `app` service in your **docker-compose.yml** file. This maps a `.ssh` folder in the container to the `.ssh` folder on your host machine.

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
- **[CodeIgniter](/docs/2-Presets/CodeIgniter.md)**
- **[Express.js](/docs/2-Presets/ExpressJS.md)**
- **[Hugo](/docs/2-Presets/Hugo.md)**
- **[Laravel](/docs/2-Presets/Laravel.md)**
- **[Next.js](/docs/2-Presets/NextJS.md)**
- **[Nuxt.js](/docs/2-Presets/NuxtJS.md)**
- **[PHP](/docs/2-Presets/PHP.md)**
- **[Symfony](/docs/2-Presets/Symfony.md)**
- **[WordPress](/docs/2-Presets/WordPress.md)**

Missing a preset? **[Make a request](https://github.com/kool-dev/kool/issues/new)**, or contribute by opening a Pull Request. Go to [https://github.com/kool-dev/kool/tree/main/presets](https://github.com/kool-dev/kool/tree/main/presets) and browse the code to learn more about how presets work.
