**Kool** is a suite of open source dev tools that level up the way you build and deploy modern, cloud native web applications. By using Kool, you'll get a better development environment and a better development workflow, which adds up to **a better developer experience**.

<a name="better-development-environment"></a>
## A Better Development Environment

### Kool CLI

The `kool` CLI is the "command center" of Kool's suite of dev tools. If you haven't done so already, you first need to [install Docker and the Kool CLI](https://kool.dev/docs/getting-started/installation). And then, when you have time, get up-to-speed with the [Kool CLI commands](https://kool.dev/docs/commands/kool).

> If you already have `kool` installed, make sure you're running the latest version with `kool self-update`.

### Kool Presets

Out of the box, Kool ships with a collection of Presets that will help you quickly kickstart local development environments for new and existing projects using popular frameworks and tech stacks.

> As an example of the Preset developer experience, check out the [Laravel preset](https://kool.dev/docs/presets/laravel).

- [AdonisJs](/docs/2-Presets/AdonisJs.md)
- [CodeIgniter](/docs/2-Presets/CodeIgniter.md)
- [Hugo](/docs/2-Presets/Hugo.md)
- [Laravel](/docs/2-Presets/Laravel.md)
- [NestJS](/docs/2-Presets/NestJS.md)
- [Next.js](/docs/2-Presets/NextJS.md)
- [Node.js](/docs/2-Presets/NodeJS.md)
- [Nuxt.js](/docs/2-Presets/NuxtJS.md)
- [PHP](/docs/2-Presets/PHP.md)
- [Symfony](/docs/2-Presets/Symfony.md)
- [WordPress](/docs/2-Presets/WordPress.md)

### Docker and Docker Compose

Kool works with Docker and Docker Compose under the hood to spin up small, fast and reproducible local development environments. Kool provides a super easy, hassle-free way to handle the Docker basics and immediately start using containers for development, while simultaneously guaranteeing no loss of control over fully customizing and extending more specialized Docker environments. You always remain fully in charge of every detail of your Docker environments.

#### docker-compose.yml

This is your Docker Compose configuration file. It should be placed inside your project and committed to version control. This file defines all the service containers needed to run your application (Docker images, ports, volume mounts, networks, etc), and follows the [Docker Compose implementation of the Compose format](https://docs.docker.com/compose/compose-file/). Over time, you'll probably make changes to this file according to the evolving needs of your project.

#### Kool Docker Images

When you start developing in containers, you quickly realize the official Docker images are built for deployment, and not for the special considerations (and nuances) of local development. One of the most common and recurring problems we see are permission issues with mapped volumes, due to host users being different from container users. Kool fixes this problem, and many others, by creating custom Docker images optimized for local development environments. [Learn more about how Kool optimizes its Docker images](https://kool.dev/docs/docker-images/introduction).

### Environment Variables

Kool loads environment variables from a **.env** file. If there's a **.env.local** file, it will take precedence and get loaded first, overriding variables in the **.env** file which use the exact same name. This helps define host-specific settings that are only applicable to your local machine.

> It's important to keep in mind that **real** environment variables defined inline as `VAR=value kool ...` or via `export VAR=value` will take precedence over the same variables in your **.env** files.

### Kool Snippets

Add services to your development environment by copying and pasting code snippets into your project's **docker-compose.yml** file. We just have a couple at the moment, but we'll be adding more soon.

- [Generate PDFs](https://kool.dev/docs/snippets/generate-pdfs)
- [Mailhog](https://kool.dev/docs/snippets/mailhog)

### kool.yml

This is your Kool configuration file. It defines scripts (commands) that you execute in your local environment or CI/CD workflows. It should be placed inside your project and committed to version control. Think of **kool.yml** as a super easy-to-use task helper. Instead of writing custom shell scripts, you can add your own scripts to **kool.yml** (under the scripts key), and run them with `kool run SCRIPT`. You can add single line commands (`kool run artisan`), or add a list of commands that will be executed in sequence (`kool run setup`).

> Use **environment variables** within your scripts to **parameterize** them, and give them an extra bit of power and flexibility.

Every **Kool Preset** auto-generates a **kool.yml** file with prebuilt scripts customized for that framework and/or stack. You can then modify this file and add your own custom scripts to facilitate your development process and share knowledge across your team. For example, add scripts to run database migrations, reset local environments, run static analysis tools, and so on. Think about how much easier it will be to onboard a teammate or new developer :)

```yaml
# ./kool.yml for the Laravel preset

scripts:
  artisan: kool exec app php artisan
  composer: kool exec app composer
  mysql: kool exec -e MYSQL_PWD=$DB_PASSWORD database mysql -uroot
  node: kool docker kooldev/node:14 node
  npm: kool docker kooldev/node:14 npm # or yarn
  npx: kool exec app npx

  node-setup:
    - kool run npm install # or yarn install
    - kool run npm run dev # or yarn dev

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
```

For example, use the `artisan` script in Laravel's **kool.yml** as follows:

```bash
kool run artisan tinker
```

#### Types of Commands

The **kool.yml** file is not just for **kool** commands. You can add any type of command that you usually run in your shell, such as `cat`, `cp`, `mv`, etc.

However, there's one important caveat – the commands you add to **kool.yml** are parsed and executed by the **kool** binary, and not in a general **bash** context. This means you **cannot** directly use **bash** control structures like `if []; then fi`, or redirection with pipes (`cmd | cmd2`). With that said, **we do support** input and output redirection (see below).

> If you need to add a more complex shell command, you can use something like `kool docker <some bash image> bash -c ""`, which will parse any bash script you need.

#### Adding Arguments

At the end of single line scripts like `kool run SCRIPT`, you can also add arguments you want to pass down to the encapsulated command. Single line commands, such as `artisan`, are like aliases, whereby additional arguments are forwarded to the actual command. For example, `kool run artisan key:generate` basically becomes `kool exec app php artisan key:generate`.

At this time, adding arguments is **only supported** by **single line** scripts. **Multi-line** scripts (like `setup`) will return an error if an extra argument is added to the end (i.e. `kool run setup some-argument`).

#### Input and Output Redirects

While commands in **kool.yml** may not run under an actual shell, we do support some shell syntax like input and output redirects. This means you can do things like the following:

```yaml
# ./kool.yml

scripts:
  # perform input injection from files
  import-db: kool docker mysql:8 -hhost -uuser -p db < path/to/some/file.sql

  # redirect standard output to a file
  write-output: echo "write something" > output.txt

  # redirect standard output to a file in append mode
  append-output: echo "add something else as a new line" >> output.txt

  # multi-redirect within a single command
  input-and-output: cat < some-file > some-new-file
```

Of course, the syntax is not as flexible as you would get directly in **bash**.

The redirect key **must be a single argument** (not glued to other arguments):
  - Correct: `write: echo "something" > output.txt`
  - Wrong: `write: echo "something">output.txt`

When performing an output redirect, the last argument after the redirect key **must be a single file destination**.

#### Learn More

Learn more by taking a closer look at the **kool.yml** files in our [Presets](https://github.com/kool-dev/kool/tree/master/presets). They contain good examples of prebuilt commands that are ready to use in a handful of different stacks. If you need help creating custom scripts based on your own unique needs, don't hesitate to [ask us on Slack](https://kool.dev/slack).

<a name="better-development-workflow"></a>
## A Better Development Workflow

Start work on a project by moving into the project directory and using `kool start` to spin up your local Docker environment. Once you're up and running, use the `kool` CLI to level up your development workflow.

### Run Commands

Use the [`kool exec` command](https://kool.dev/docs/commands/kool-exec) to execute a command inside a running service container.

```bash
# kool exec [OPTIONS] SERVICE COMMAND [--] [ARG...]

$ kool exec app ls
```

Look at the **kool.yml** files in the [Presets](https://github.com/kool-dev/kool/tree/master/presets) to see other examples of `kool exec` commands used in scripts.

### Open Sessions

Similar to SSH, if you want to open a Bash session in your `app` container, run `kool exec app bash`, where `app` is the name of the service container in your **docker-compose.yml**. If you prefer, you can use `sh` instead of `bash` (`kool exec app sh`).

```bash
$ kool exec app bash
bash-5.1#

$ kool exec app sh
/app #
```

### Connect to the Database

If you added a database service, start a new SQL client session inside your running `database` container by executing `kool run mysql` (MySQL) or `kool run psql` (PostgreSQL). This runs the `mysql` or `psql` script in your **kool.yml**.

### Add Dependencies

If you're working on a Node-based project, and you add more dependencies to `package.json`, use the [`kool restart` command](https://kool.dev/docs/commands/kool-restart) to restart your `app` container and load the new packages.

```console
$ kool restart app

Stopping my-project_app_1 ... done
Going to remove my-project_app_1
Removing my-project_app_1 ... done
Creating my-project_app_1 ... done
```

### Access Private Repos and Packages

If you need your `app` container to use your local SSH keys to pull private repositories and/or install private packages (which have been added as dependencies in your `composer.json` or `package.json` file), you can simply add `$HOME/.ssh:/home/kool/.ssh:delegated` under the `volumes` key of the `app` service in your **docker-compose.yml** file. This maps a `.ssh` folder in the container to the `.ssh` folder on your host machine.

```diff
volumes:
  - .:/app:delegated
+ - $HOME/.ssh:/home/kool/.ssh:delegated
```

### View the Logs

View container logs using the [`kool logs` command](https://kool.dev/docs/commands/kool-logs). Run `kool logs` to see the logs for all running containers, or `kool logs app` to specify a service and only see the logs for the `app` container. Add the `-f` option after `kool logs` to follow the logs (i.e. `kool logs -f app`).

### Share Your Work

When you need to quickly share local changes with your team, use the [`kool share` command](https://kool.dev/docs/commands/kool-share) to share your local environment over the Internet via an HTTP tunnel. Specify your own subdomain using the `--subdomain` flag. Since the default port is `80`, you'll also need to use the `--port` option for Node apps (i.e. `kool share --port=3000`).

```console
$ kool share

Thank you for using expose.
Local-URL:     app
Dashboard-URL: http://127.0.0.1:4040
Expose-URL:    https://eeskzijcbe.kool.live

Remaining time: 00:59:59
Remaining time: 00:59:58
Remaining time: 00:59:57
```

### Switch Projects

Kool supports any language or framework, so you can standardize the way you work across all your tech stacks. When it's time to stop working on your current app and switch to a different project, you can easily change local Docker environments by running `kool stop`, moving into the other project directory, and running `kool start`.

```console
$ kool stop
$ cd ~/my-other-project
$ kool start
```

Pretty _kool_, right?

> If you like what we're doing, show your support for this new open source project by [**starring us on GitHub**](https://github.com/kool-dev/kool)!
