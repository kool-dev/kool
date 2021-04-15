## Start a New Adonis Project with Docker in 3 Easy Steps

1. Run `kool create adonis my-project`
2. Update **.env**
3. Run `kool run setup`

> Yes, using **kool** + Docker to create and work on new Adonis projects is that easy!

### Requirements

If you haven't done so already, you first need to [install Docker and the kool CLI](/docs/getting-started/installation).

> Please note that you'll need a basic understanding of how Docker and Docker Compose work in order to build a new project from scratch using Kool with Docker.

Also, make sure you're running the latest version of **kool**. Run the following command to compare your local version of **kool** with the latest release, and, if a newer version is available, automatically download and install it.

```bash
$ kool self-update
```

### 1. Run `kool create adonis my-project`

Use the [`kool create <preset> <project-name>` command](docs/commands/kool-create) to create your new Adonis project:

```bash
$ kool create adonis my-project
```

Under the hood, this command will run `adonis new <project-name>` (using a customized **kool** Docker image: <a href="https://github.com/kool-dev/docker-node" target="_blank">kooldev/node:14-adonis</a>), which installs the Adonis [fullstack blueprint](https://github.com/adonisjs/adonis-fullstack-app).

After installing Adonis, `kool create` automatically runs the `kool preset adonis` command, which helps you easily set up the initial tech stack for your project using an interactive wizard.

```bash
$ Preset adonis is initializing!

? What database service do you want to use [Use arrows to move, type to filter]
  MySQL 8.0
  MySQL 5.7
  PostgreSQL 13.0
> none

? What cache service do you want to use [Use arrows to move, type to filter]
  Redis 6.0
  Memcached 1.6
> none

? What javascript package manager do you want to use [Use arrows to move, type to filter]
> npm
  yarn

$ Preset adonis initialized!
```

Now, move into your new Adonis project:

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

### 2. Update .env

You need to update your **.env** file with some values that match the services in your **docker-compose.yml** file.

> We recommend you make the same changes in your default **.env.example** file.

#### Host

```diff
-HOST=127.0.0.1
+HOST=0.0.0.0
```

#### Database Services

MySQL 5.7 and 8.0

```diff
-DB_CONNECTION=sqlite
+DB_CONNECTION=mysql

-DB_HOST=127.0.0.1
+DB_HOST=database

-DB_USER=root
+DB_USERNAME=root
```

PostgreSQL 13.0

```diff
-DB_CONNECTION=sqlite
+DB_CONNECTION=pgsql

-DB_HOST=127.0.0.1
+DB_HOST=database

-DB_PORT=3306
+DB_PORT=5432

-DB_USER=root
+DB_USERNAME=root
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
+MEMCACHED_PORT=?
```

### 3. Run `kool run setup`

> Say hello to **kool.yml**, say goodbye to custom shell scripts!

As mentioned above, the [`kool preset` command](docs/commands/kool-preset) added a **kool.yml** file to your project. Think of **kool.yml** as a super easy-to-use task _helper_. Instead of writing custom shell scripts, add your own scripts to **kool.yml** (under the `scripts` key), and run them with `kool run <script>` (e.g. `kool run adonis`). You can add your own single line commands (see `adonis` below), or add a list of commands that will be executed in sequence (see `setup` below).

To help get you started, **kool.yml** comes prebuilt with an initial set of scripts (based on the **preset**).

```yaml
scripts:
  node: kool exec app node
  npm: kool exec app npm # can change to: yarn,pnpm
  adonis: kool exec app adonis

  setup:
    - kool docker kooldev/node:14 npm install # can change to: yarn,pnpm
    - kool start
```

> Try `kool run adonis --help` to execute the `kool exec app adonis --help` command in your running `app` container and print out information about Adonis' commands.

**kool.yml** always includes a script called `setup`, which helps you spin up a project for the first time.

Go ahead and run `kool run setup` to start your Docker environment and finish setting up your project:

```bash
$ kool run setup
```

> As you can see in **kool.yml**, the `setup` script will do the following in sequence: run `npm install` to build your Node packages and dependencies (by spinning up and down a temporary Node container); and then start your Docker environment using **docker-compose.yml** (which includes a `command` to automatically run `adonis serve --dev` for you).

Verify your Docker container is running using the [`kool status` command](docs/commands/kool-status):

```bash
$ kool status

+---------+---------+------------------------+--------------+
| SERVICE | RUNNING | PORTS                  | STATE        |
+---------+---------+------------------------+--------------+
| app     | Running | 0.0.0.0:3333->3333/tcp | Up 5 seconds |
+---------+---------+------------------------+--------------+
```

Run `kool logs` to see the logs from the running Docker service containers, and confirm the Adonis server was started (and add the `-f` option to follow the logs):

```
kool logs
Attaching to my-project_app_1
app_1  |
app_1  |  SERVER STARTED
app_1  | > Watching files for changes...
app_1  |
app_1  | info: serving app on http://0.0.0.0:3333
```

 Once `kool run setup` finishes, you should now be able to access your new app at [http://localhost:3333](http://localhost:3333) and see the Adonis "It works!" welcome page. Hooray!

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

If you added a database service, you can easily start a new SQL client session inside your Docker `database` container by executing `kool run mysql` (MySQL) or `kool run psql` (PostgreSQL) in your terminal. This runs the single-line `mysql` or `psql` script included in your **kool.yml**.

### Staying kool

When it's time to stop working on the project:

```bash
$ kool stop
```

And when you're ready to start work again:

```bash
$ kool start
```
