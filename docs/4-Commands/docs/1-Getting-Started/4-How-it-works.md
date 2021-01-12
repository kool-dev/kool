Kool works with Docker / Docker Compose under the hood, and comes with some cool presets to help you get started, everything is configurable / extendable.

Let's use **Laravel** preset as example and explain how it works.

When you run **kool preset laravel** all it does is create a few files for you:

```bash
$ kool preset laravel
Preset laravel is initializing!
  Preset file Dockerfile.build created.
  Preset file docker-compose.yml created.
  Preset file kool.yml created.
Preset laravel initialized!
```

### Dockerfile.build

This is a file you can use in case you want to build your application to use in production, in Docker world every release is usually a new image built with your application on it.

Soon we will give more examples on how to use Docker in production or use it with **Kool Cloud**.

### docker-compose.yml

This file defines all services that runs your application, docker images to use, ports, volume mounts, etc.

You can add/change/remove services as you will.

### kool.yml

This is where most of the magic happens, a way to make your life easy, encapsulating scripts for you to use on your local environment or CI/CDs. It is created in your working directory when you run **kool preset**, but you can also create it inside a folder named **kool** in your user's home directory.

The **scripts** defined can be used with **kool run <script>** command.

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

You can pass in after kool run any options or arguments you wish to pass down the encapsulated command.

#### Arguments to kool run <script>

Single commands like **artisan** are kind of aliases, so anything you input will be forwarded to the actual command, so if you run: **kool run artisan key:generate** it will basically translate into: **kool exec app php artisan key:generate**.

Multiple commands like **setup** will not forward your input, so **kool run setup something** will run every script and **something** will be ignored.

#### What kind of commands can be encasulated on `kool.yml`

This is not meant only for `kool` commands, you can add any type commands as you usually run them in your shell like `cat`, `cp`, `mv`, etc.

There is just one caveat we need to be aware of - the commands within a script on `kool.yml` are parsed and executed by `kool` and not in a general `bash` context, so you **cannot** directly use bash script structures like `if []; then fi`. In case you need something for that effect, you should use a `kool docker <some bash image> bash -c ""` which then parses any bash script you need.

#### Input and output redirects on `kool.yml`

Although the previous notice about commands within a script at `kool.yml` is not straight out a bash script, we do support some bash helping syntax like input and output redirects.

So you are totally able to do things like:

kool.yml
```yaml
scripts:
  # performing input injection from files
  import-db: kool docker mysql:8 -h$DB_HOST -u$DB_USER -p $DB_NAME < path/to/some/file.sql

  # redirecting standard output to a file
  write-output: echo "writing something" > output.txt

  # redirecting standard output to a file in append mode
  append-output: echo "something else in a new line" >> output.txt
```

Again, of course the syntax is not as flexible as you would have in straight out `bash`, please notice:

- The redirect key must be a single argument (not glued to the other arguments).
    - Correct: `write: echo "something" > output`
    - Wrong: `write: echo "something">output`
- When performing a redirect, the last argument after the redirect key must be a single file destination.

Hope you enjoy this feature! Take a look at the presets which already contain good examples of `kool.yml` files ready to be used in a handlful of different stacks. In case you need help yo create your own based on your needs make sure to ask for help on Github.
