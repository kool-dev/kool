# Start a Hugo Project with Docker in 2 Easy Steps

1. Run `kool create hugo my-project`
2. Run `kool run setup`

> Yes, using **kool** + Docker to create and work on new PHP projects is that easy!

## Requirements

If you haven't done so already, you first need to [install Docker and the kool CLI](/docs/getting-started/installation).

Also, make sure you're running the latest version of **kool**. Run the following command to compare your local version of **kool** with the latest release, and, if a newer version is available, automatically download and install it.

```bash
$ kool self-update
```

> Please note that it helps to have a basic understanding of how Docker and Docker Compose work to use Kool with Docker.

## 1. Run `kool create hugo my-project`

Use the [`kool create PRESET FOLDER` command](/docs/commands/kool-create) to create your new Hugo project:

> IMPORTANT: if you're on **Windows WSL** or **Linux**, you should run `sudo kool create hugo my-project` as the superuser (via `sudo`) to avoid permissions issues when creating the project directory and files.

```bash
$ kool create hugo my-project
```

Under the hood, this command will run `kool docker klakegg/hugo:ext-alpine new site my-project` using the [klakegg/hugo](https://hub.docker.com/r/klakegg/hugo/) Docker image.

Now, move into your new Hugo project:

```bash
$ cd my-project
```

The [`kool preset` command](/docs/commands/kool-preset) auto-generated the following configuration files and added them to your project, which you can modify and extend.

```bash
+docker-compose.yml
+kool.yml
```

> Now's a good time to review the **docker-compose.yml** file and verify the services suit the needs of your project.

## 2. Run `kool run quickstart`

> Say hello to **kool.yml**, say goodbye to custom shell scripts!

As mentioned above, the [`kool preset` command](/docs/commands/kool-preset) added a **kool.yml** file to your project. Think of **kool.yml** as a super easy-to-use task _helper_. Instead of writing custom shell scripts, add your own scripts to **kool.yml** (under the `scripts` key), and run them with `kool run SCRIPT` (e.g. `kool run hugo`). You can add your own single line commands (see `hugo` below), or add a list of commands that will be executed in sequence (see `setup` below).

To help get you started, **kool.yml** comes prebuilt with an initial set of scripts (based on the choices you made earlier using the **preset** wizard), including a script called `setup`, which helps you spin up a project for the first time. However, since Hugo requires a few extra steps to create a Hello World site, we've added a special `quickstart` script to make it super easy.

```yaml
scripts:
	hugo: kool docker -p 1313:1313 klakegg/hugo:ext-alpine
	dev: kool run hugo server -D

	# remove or modify to suit the needs of your project
	quickstart:
		- kool start
		- git init
		- git submodule add https://github.com/theNewDynamic/gohugo-theme-ananke.git themes/ananke
		- echo theme = \"ananke\" >> config.toml
		- kool run hugo new posts/my-first-post.md
		- kool run dev

	setup:
		- kool start
		- kool run dev
```

Go ahead and run `kool run quickstart` to start your Docker environment and initialize your Hugo site:

```bash
$ kool run quickstart
```

> As you can see in **kool.yml**, the `quickstart` script will do the following in sequence: run the `kool start` command to spin up your Docker environment; call `git init` to create a Git repository; download the Ananke theme; use an `echo` command to add the theme to your Hugo config file; add your first post; and then call `kool run dev` to build your Hugo site.

Once `kool run quickstart` finishes, you should be able to access your new site at [http://localhost](http://localhost/) and see the "My New Hugo Site" page. Hooray!

Verify your Docker container is running using the [`kool status` command](/docs/commands/kool-status).

```bash
$ kool status
+---------+---------+------------------------------+--------------+
| SERVICE | RUNNING | PORTS                        | STATE        |
+---------+---------+------------------------------+--------------+
| app     | Running | 0.0.0.0:80->80/tcp, 1313/tcp | Up 2 minutes |
| static  | Running | 80/tcp                       | Up 2 minutes |
+---------+---------+------------------------------+--------------+
```

Run `kool logs app` to see the logs from the running `app` container.

> Use `kool logs` to see the logs from all running containers. Add the `-f` option after `kool logs` to follow the logs (i.e. `kool logs -f app`).

```bash
$ kool logs app
Attaching to my-project_app_1
app_1     |   Non-page files   |  0
app_1     |   Static files     |  0
app_1     |   Processed images |  0
app_1     |   Aliases          |  0
app_1     |   Sitemaps         |  1
app_1     |   Cleaned          |  0
app_1     |
app_1     | Built in 1 ms
app_1     | Watching for changes in /app/{archetypes,content,data,layouts,static}
app_1     | Watching for config changes in /app/config.toml
app_1     | Environment: "DEV"
app_1     | Serving pages from memory
app_1     | Running in Fast Render Mode. For full rebuilds on change: hugo server --disableFastRender
app_1     | Web Server is available at http://localhost:80/ (bind address 0.0.0.0)
app_1     | Press Ctrl+C to stop
app_1     |
app_1     | Change of config file detected, rebuilding site.
app_1     | 2021-05-01 20:34:06.306 +0000
app_1     | Rebuilt in 136 ms
app_1     | adding created directory to watchlist /app/content/posts
app_1     |
app_1     | Change detected, rebuilding site.
app_1     | 2021-05-01 20:34:07.305 +0000
app_1     | Source changed "/app/content/posts/my-first-post.md": CREATE
app_1     | Total in 26 ms
```

---

### Run Commands in Docker Containers

Use [`kool exec`](/docs/commands/kool-exec) to execute a command inside a running service container:

```bash
# kool exec [OPTIONS] SERVICE COMMAND [--] [ARG...]

$ kool exec app ls
```

Try `kool run hugo version` to execute the `kool exec app hugo version` command in your running `app` container and verify your installation.

### Open Sessions in Docker Containers

Similar to SSH, if you want to open a Bash session in your `app` container, run `kool exec app sh`, where `app` is the name of the service container in **docker-compose.yml**.

```bash
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

- **[AdonisJs](/docs/2-Presets/AdonisJs.md)**
- **[CodeIgniter](/docs/2-Presets/CodeIgniter.md)**
- **[Express.js](/docs/2-Presets/ExpressJS.md)**
- **[Laravel](/docs/2-Presets/Laravel.md)**
- **[NestJS](/docs/2-Presets/NestJS.md)**
- **[Next.js](/docs/2-Presets/NextJS.md)**
- **[Node.js](/docs/2-Presets/NodeJS.md)**
- **[Nuxt.js](/docs/2-Presets/NuxtJS.md)**
- **[PHP](/docs/2-Presets/PHP.md)**
- **[Symfony](/docs/2-Presets/Symfony.md)**
- **[WordPress](/docs/2-Presets/WordPress.md)**

Missing a preset? **[Make a request](https://github.com/kool-dev/kool/issues/new)**, or contribute by opening a Pull Request. Go to [https://github.com/kool-dev/kool/tree/main/presets](https://github.com/kool-dev/kool/tree/main/presets) and browse the code to learn more about how presets work.
