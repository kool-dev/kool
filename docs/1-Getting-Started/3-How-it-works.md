**Kool** works with Docker and Docker Compose under the hood to power up small, fast and reproducible local development environments. **Kool** provides simple, no-brainer ease-of-use while keeping all the power and flexibility of Docker and Docker Compose. **Kool** makes using Docker super easy without losing any control, ensuring you remain fully in charge of every detail related to your Docker environments.

### Presets

Out of the box, **kool** ships with a collection of presets to help you quickly kickstart local development using some popular frameworks and stacks. Check out the [Laravel preset](https://kool.dev/docs/presets/laravel) as an example of the developer experience **kool** offers.

### docker-compose.yml

This is the Docker Compose configuration file, and it should be placed inside your project and committed to version control. This file defines all the service containers needed to run your application (the Docker images to use, ports, volume mounts, etc). It follows the [Docker Compose implementation of the Compose format](https://docs.docker.com/compose/compose-file/). Over time, you'll probably make tweaks and improvements to this file according to the specific needs of your project.

### Environment Variables

**Kool** loads environment variables from a **.env** file. If there's a **.env.local** file, it will take precedence and get loaded first, overriding variables in the **.env** file which use the exact same name. This helps define host-specific settings that are only applicable to your local machine.

> It's important to keep in mind that **real** environment variables win (take precedence) over variables defined in your **.env** files.

### kool.yml

This is the **kool** configuration file, and it should be placed inside your project and committed to version control. This file defines scripts (commands) that you execute in your local environment or CI/CDs. Think of **kool.yml** as a super easy-to-use task helper. Instead of writing custom shell scripts, add your own scripts to **kool.yml** (under the `scripts:` root key), and run them with `kool run <script-name>`. You can add your own single line commands, or add a list of commands that will be executed in sequence.

> Use **environment variables** within your scripts to **parameterize** them, and give them an extra bit of power and flexibility.

Every **preset** includes a **kool.yml** file with prebuilt scripts for that stack. Of course, you can add your own custom scripts to facilitate your development process and share knowledge across the team.

```yaml
# ./kool.yml for Laravel preset

scripts:
  artisan: kool exec app php artisan

  setup:
    - kool start
    - kool run artisan key:generate
```

Using the `artisan` command:

```bash
kool run artisan tinker
```

#### Types of Commands

The **kool.yml** file is not just for **kool** commands. You can add any type of command you usually run in your shell, such as `cat`, `cp`, `mv`, etc.

However, there is one important caveat - the commands you add to **kool.yml** are parsed and executed by the **kool** binary, and not in a general **bash** context. This means you **cannot** directly use **bash** control structures like `if []; then fi`, or redirection with pipes (`cmd | cmd2`). With that said, **we do support** input and output redirection (see below). W00t!

> If you need to add a more complex shell command, you can use something like `kool docker <some bash image> bash -c ""`, which will parse any bash script you need.

#### Adding Arguments

At the end of single line scripts like `kool run <script-name>`, you can also add arguments you want to pass down to the encapsulated command. Single line commands, such as `artisan`, are like aliases, whereby additional arguments are forwarded to the actual command. For example, `kool run artisan key:generate` basically becomes `kool exec app php artisan key:generate`.

At this time, adding arguments is **only supported** by **single line** commands. **Multi-line** commands (like `setup`) will return an error if an extra argument is added to the end (i.e. `kool run setup something`).

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

- The redirect key **must be a single argument** (not glued to other arguments)
	- Correct: `write: echo "something" > output.txt`
	- Wrong: `write: echo "something">output.txt`
- When performing an output redirect, the last argument after the redirect key **must be a single file destination**

#### Learn More

Learn more by taking a closer look at the **kool.yml** files in our [presets](https://kool.dev/docs/presets/introduction). They contain good examples of prebuilt commands that are ready to use in a handful of different stacks. If you need help creating custom scripts based on your own unique needs, don't hesitate to ask on GitHub.
