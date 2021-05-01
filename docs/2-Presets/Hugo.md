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

```bash
$ kool create hugo my-project
```

Under the hood, this command will run `kool docker klakegg/hugo new site my-project` using the [klakegg/hugo](https://hub.docker.com/r/klakegg/hugo/) Docker image.

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

## 2. Run `kool run setup`

> Say hello to **kool.yml**, say goodbye to custom shell scripts!

As mentioned above, the [`kool preset` command](/docs/commands/kool-preset) added a **kool.yml** file to your project. Think of **kool.yml** as a super easy-to-use task _helper_. Instead of writing custom shell scripts, add your own scripts to **kool.yml** (under the `scripts` key), and run them with `kool run SCRIPT` (e.g. `kool run hugo`). You can add your own single line commands (see `hugo` below), or add a list of commands that will be executed in sequence (see `setup` below).

To help get you started, **kool.yml** comes prebuilt with an initial set of scripts (based on the choices you made earlier using the **preset** wizard), including a script called `setup`, which helps you spin up a project for the first time.

```yaml
scripts:
  hugo: kool docker -p 1313:1313 klakegg/hugo
  dev: kool run hugo server -D

  setup:
    - kool start
    - kool run dev
```

Go ahead and run `kool run setup` to start your Docker environment:

```bash
$ kool run setup
```

> As you can see in **kool.yml**, the `setup` script will do the following in sequence: run the `kool start` command to spin up your Docker environment; and then calls `kool run dev` to build your Hugo site.

Once `kool run setup` finishes, you should be able to access your new site at [http://localhost:1313](http://localhost:1313/). Hooray!.

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
app_1     | Start building sites …
app_1     | WARN 2021/05/01 18:12:00 found no layout file for "HTML" for kind "home": You should create a template file which matches Hugo Layouts Lookup Rules for this combination.
app_1     | WARN 2021/05/01 18:12:00 found no layout file for "HTML" for kind "taxonomy": You should create a template file which matches Hugo Layouts Lookup Rules for this combination.
app_1     | WARN 2021/05/01 18:12:00 found no layout file for "HTML" for kind "taxonomy": You should create a template file which matches Hugo Layouts Lookup Rules for this combination.
app_1     |
app_1     |                    | EN
app_1     | -------------------+-----
app_1     |   Pages            |  3
app_1     |   Paginator pages  |  0
app_1     |   Non-page files   |  0
app_1     |   Static files     |  0
app_1     |   Processed images |  0
app_1     |   Aliases          |  0
app_1     |   Sitemaps         |  1
app_1     |   Cleaned          |  0
app_1     |
app_1     | Built in 8 ms
app_1     | Watching for changes in /app/{archetypes,content,data,layouts,static}
app_1     | Watching for config changes in /app/config.toml
app_1     | Environment: "DEV"
app_1     | Serving pages from memory
app_1     | Running in Fast Render Mode. For full rebuilds on change: hugo server --disableFastRender
app_1     | Web Server is available at http://localhost:80/ (bind address 0.0.0.0)
app_1     | Press Ctrl+C to stop
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
- **[Laravel](/docs/2-Presets/Laravel.md)**
- **[NestJS](/docs/2-Presets/NestJS.md)**
- **[Next.js](/docs/2-Presets/NextJS.md)**
- **[Node.js](/docs/2-Presets/NodeJS.md)**
- **[Nuxt.js](/docs/2-Presets/NuxtJS.md)**
- **[PHP](/docs/2-Presets/PHP.md)**
- **[Symfony](/docs/2-Presets/Symfony.md)**

Missing a preset? **[Make a request](https://github.com/kool-dev/kool/issues/new)**, or contribute by opening a Pull Request. Go to [https://github.com/kool-dev/kool/tree/master/presets](https://github.com/kool-dev/kool/tree/master/presets) and browse the code to learn more about how presets work.




### Using Hugo preset

#### Creating a new Hugo website

To make things easier we will use **kool** to install it for you.

```console
$ kool create hugo my-website

$ cd my-website
```
- **kool create** already executes **kool preset** internally so you can skip the command in the next step.

#### Adding kool to an existing Hugo website

Go to the project folder and run:

```console
$ cd my-website/
$ kool preset hugo
```

**kool preset** will create a few configuration files in order to enable you to configure / extend it. You don't need to execute it if you created the project with `kool create`.

### Using kool for Hugo development

- To start the container to serve your Hugo website:

```console
$ kool start
```

Then check out your site at `http://localhost`. If you wanna stop the container just run `kool stop`.

- To create some new content:

```console
$ kool run hugo new posts/my-super-post.md
```

---

Check your **kool.yml** to see what scripts you can run and add more.
