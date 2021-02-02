### Laravel Example

Check presets for how to create other types of projects.

To make things easier, we will use **kool** to install it for you.

```bash
kool create laravel my-project

cd my-project
```
- **kool create** already executes **kool preset** internally, so you can skip the command in the next step.

Or,

```bash
kool docker kooldev/php:7.4 composer create-project --prefer-dist laravel/laravel my-project

cd my-project
```

### Start Using kool

Go to the project folder and run:

```bash
$ kool preset laravel
```

Basically, **kool preset** creates a few configuration files that you can configure and extend. You don't need to execute it if you ran the **kool create** command to start the new project.

By default, the Laravel preset comes with **mysql** and **redis** already configured for you. You can review how it's configured in **docker-compose.yml**. The preset also comes with some scripts in **kool.yml** to help bring you up-to-speed. Take a look at the defaults.

We always add a script called **setup** to help you set up a project for the first time.

```bash
# CAUTION, this script will reset your `.env` file with `.env.example`
$ kool run setup
```

Now you should see your site at **http://localhost**.

To get back to work on a project:

```bash
$ kool start
```

Then, when you're done for the day:

```bash
$ kool stop
```

Check your **kool.yml** to see what scripts you can run, and then add more.
