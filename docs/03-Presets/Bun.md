# Start a Bun Project with Docker in 2 Easy Steps

1. Run `kool create bunjs my-project`
2. Run `kool run setup`

> Yes, using **kool** + Docker to create and work on new [Bun](https://bun.sh) projects is that easy!

## Requirements

If you haven't done so already, you first need to [install Docker and the kool CLI](/docs/getting-started/installation).

Also, make sure you're running the latest version of **kool**. Run the following command to compare your local version of **kool** with the latest release, and, if a newer version is available, automatically download and install it.

```bash
$ kool self-update
```

> Please note that it helps to have a basic understanding of how Docker and Docker Compose work to use Kool with Docker.

## 1. Run `kool create bunjs my-project`

Use the [`kool create PRESET FOLDER` command](/docs/commands/kool-create) to create your new Bun project:

```bash
$ kool create bunjs my-project
```

Under the hood, this command will create a "Hello World" **app.js** file (powered by `Bun.serve`) and a minimal **package.json** in the root of your new project directory, using the official <a href="https://hub.docker.com/r/oven/bun" target="_blank">oven/bun:1</a> Docker image.

After creating the project, `kool create` automatically runs the `kool preset bunjs` command, which sets up the initial tech stack for your project.

```bash
$ Preset bunjs is initializing!

...

Preset bunjs created successfully!
```

Now, move into your new Bun project:

```bash
$ cd my-project
```

The [`kool preset` command](/docs/commands/kool-preset) auto-generated the following configuration files and added them to your project, which you can modify and extend.

```bash
+docker-compose.yml
+kool.yml
+app.js
+package.json
```

> Now's a good time to review the services added to the **docker-compose.yml** file. The `app` service runs the <a href="https://hub.docker.com/r/oven/bun" target="_blank">oven/bun:1</a> image with a `command` of `bun app.js`.

## 2. Run `kool run setup`

> Say hello to **kool.yml**, say goodbye to custom shell scripts!

As mentioned above, the [`kool preset` command](/docs/commands/kool-preset) added a **kool.yml** file to your project. Think of **kool.yml** as a super easy-to-use task _helper_. Instead of writing custom shell scripts, add your own scripts to **kool.yml** (under the `scripts` key), and run them with `kool run SCRIPT` (e.g. `kool run bun`). You can add your own single line commands (see `bun` below), or add a list of commands that will be executed in sequence.

To help get you started, **kool.yml** comes prebuilt with an initial set of scripts (based on the **preset**).

```yaml
scripts:
  bun: kool exec app bun
  bunx: kool exec app bunx
```

Go ahead and run `kool start` to start running the container:

```bash
$ kool start
```

> The **docker-compose.yml** file includes a `command` to automatically run `bun app.js` when the `app` container starts.

Once the container is up, you should be able to access your new site at [http://localhost:3000](http://localhost:3000) and see the "Hello World" page. Hooray!

Verify your Docker container is running using the [`kool status` command](/docs/commands/kool-status):

```bash
$ kool status
+---------+---------+-------------------------------------------+--------------+
| SERVICE | RUNNING | PORTS                                     | STATE        |
+---------+---------+-------------------------------------------+--------------+
| app     | Running | 0.0.0.0:3000->3000/tcp, :::3000->3000/tcp | Up 4 seconds |
+---------+---------+-------------------------------------------+--------------+
```

Run `kool logs app` to see the logs from the running `app` container, and confirm the Bun server was started.

> Use `kool logs` to see the logs from all running containers. Add the `-f` option after `kool logs` to follow the logs (i.e. `kool logs -f app`).

```bash
$ kool logs app
Attaching to my-project_app_1
app_1  | Server running at http://localhost:3000/
```

---

### Install dependencies with Bun

Use the `bun` helper script to manage your dependencies. Because Bun is both the runtime and the package manager, the same `app` container is used:

```bash
$ kool run bun add hono
$ kool run bun install
```

### Run Commands in Docker Containers

Use [`kool exec`](/docs/commands/kool-exec) to execute a command inside a running service container:

```bash
# kool exec [OPTIONS] SERVICE COMMAND [--] [ARG...]

$ kool exec app bun --version
```

Try `kool run bun --help` to execute the `kool exec app bun --help` command in your running `app` container and print out information about Bun.

### Open Sessions in Docker Containers

Similar to SSH, if you want to open a Bash session in your `app` container, run `kool exec app bash`, where `app` is the name of the service container in **docker-compose.yml**. If you prefer, you can use `sh` instead of `bash` (`kool exec app sh`).

```bash
$ kool exec app bash
root@app:/app#

$ kool exec app sh
/app #
```

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

- **[AdonisJs](/docs/03-Presets/AdonisJs.md)**
- **[CodeIgniter](/docs/03-Presets/CodeIgniter.md)**
- **[Express.js](/docs/03-Presets/ExpressJS.md)**
- **[Hugo](/docs/03-Presets/Hugo.md)**
- **[Laravel](/docs/03-Presets/Laravel.md)**
- **[NestJS](/docs/03-Presets/NestJS.md)**
- **[Next.js](/docs/03-Presets/NextJS.md)**
- **[Node.js](/docs/03-Presets/NodeJS.md)**
- **[Nuxt.js](/docs/03-Presets/NuxtJS.md)**
- **[PHP](/docs/03-Presets/PHP.md)**
- **[Symfony](/docs/03-Presets/Symfony.md)**
- **[WordPress](/docs/03-Presets/WordPress.md)**

Missing a preset? **[Make a request](https://github.com/kool-dev/kool/issues/new)**, or contribute by opening a Pull Request. Go to [https://github.com/kool-dev/kool/tree/main/presets](https://github.com/kool-dev/kool/tree/main/presets) and browse the code to learn more about how presets work.
