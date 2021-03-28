Kool works with Docker and Docker Compose under the hood to power up small, fast and reproduceable local development environms. We also ship out of the box a collection of what we call *presets* to help you get started local development on some popular frameworks and stacks.

We aim at bringing easiness and no-brainer solutions while still keeping open and accessible to you absolute all the power and flexibity of Docker and Docker Compose. We make them simple to use, but we do not toll this facilitation on lack of control - you continue in charge of everything that is going on.

Check out our [**Laravel** preset as example](https://kool.dev/docs/presets/laravel) to get the feeling on the experience we offer.

### docker-compose.yml

This file defines all services that runs your application, docker images to use, ports, volume mounts, etc. This file is the configuration file for Docker Compose and follow its [reference for content and formatting](https://docs.docker.com/compose/compose-file/) to define our service containers.

This file belongs to your project and should be versioned with it. As such you are can and probably eventually will make tweaks and improvements to it, to better suit your project specific needs.

#### Environment variables

`kool` will load environment variables from a `.env` file (if there is a `.env.local` file, it gets loaded first). This is helpful for setting up host specific settings you don't want to tie up to your setup.

Previously defined variables never get overriden.

### kool.yml

This is Kool configuration file. Your best usage for it is to encapsule scripts for you to use on your local environment or CI/CDs. It is a YML file with `scripts:` root key, where you should define *scripts* (single line or an array with several lines). You are then able to call this scripts with `kool run my-script`.

Each of our *presets* comes with a helpful `kool.yml` suggestion for that stack. Of course you are encouraged to expand it with your own custom scripts to help along the development process and knowledge sharing across the team. This file should also be versioned along the project.

> Using **environment variables** within your scripts gives them an extra bit of power and flexibility. You are encouraged to use them in order to **parameterize your scripts**.

kool.yml:
```yaml
scripts:
  artisan: kool exec app php artisan

  setup:
    - kool start
    - kool run artisan key:generate
```

Usage:
```bash
kool run artisan tinker
```

You can pass in after `kool run my-script` any other arguments you wish to pass down the encapsulated command.

> The arguments passing at this moment is limited to single-line scripts only.

#### Arguments to kool run <script>

Single commands like **artisan** are kind of aliases, so anything you input will be forwarded to the actual command, so if you run: **kool run artisan key:generate** it will basically translate into: **kool exec app php artisan key:generate**.

Multiple commands like **setup** will not forward your input, so **kool run setup something** will run every script and **something** will be ignored.

#### What kind of commands can be encasulated on `kool.yml`

This is not meant only for `kool` commands, you can add any type commands as you usually run them in your shell like `cat`, `cp`, `mv`, etc.

There is just one caveat we need to be aware of - the commands within a script on `kool.yml` are parsed and executed by `kool` and not in a general `bash` context, so you **cannot** directly use bash script structures like `if []; then fi` or redirection with pipes (`cmd | cmd2`). Although, we support output redirection (see below).

> In case you need some more complex shell command, you can use something like `kool docker <some bash image> bash -c ""` which then parses any bash script you need.

#### Input and output redirects on `kool.yml`

Although the previous warning about commands within a script at `kool.yml` not running under an actual *shell*, we do support some *shell* helping syntax like input and output redirects.

So you are totally able to do things like:

kool.yml
```yaml
scripts:
  # performing input injection from files
  import-db: kool docker mysql:8 -hhost -uuser -p db < path/to/some/file.sql

  # redirecting standard output to a file
  write-output: echo "writing something" > output.txt

  # redirecting standard output to a file in append mode
  append-output: echo "something else in a new line" >> output.txt

  # is supports multi redirecting within a single command
  input-and-output: cat < some-file > some-new-file
```

Again, of course the syntax is not as flexible as you would have in straight out `bash`, please notice:

- The redirect key must be a single argument (not glued to the other arguments).
    - Correct: `write: echo "something" > output`
    - Wrong: `write: echo "something">output`
- When performing a output redirect, the last argument after the redirect key must be a single file destination.

Hope you enjoy this feature! Take a look at the presets which already contain good examples of `kool.yml` files ready to be used in a handlful of different stacks. In case you need help yo create your own based on your needs make sure to ask for help on Github.
