### Requirements

Kool is powered by **[Docker](https://docs.docker.com/get-docker/)** and **[Docker Compose](https://docs.docker.com/compose/install/)**, you need to have it installed on your machine.

### Installation

To install **kool** simply run the following script.

```bash
curl -fsSL https://raw.githubusercontent.com/kool-dev/kool/master/install.sh | bash
```

In case it fails due to permission then run it using sudo:

```bash
curl -fsSL https://raw.githubusercontent.com/kool-dev/kool/master/install.sh | sudo bash
```

All this script will do is download latest kool bin from https://github.com/kool-dev/kool/releases for your OS and put in your `/usr/local/bin` folder.

## For Windows

Download the installer [here](https://github.com/kool-dev/kool/releases)

To check if everything looks good simply run **kool** and you will see something like this:

```bash
$ kool

An easy and robust software development environment
tool helping you from project creation until deployment.
Complete documentation is available at https://kool.dev

Usage:
  kool [command]

Available Commands:
  deploy      Deploys your application usin Kool Dev
  docker      Creates a new container and runs the command in it.
  exec        Execute a command within a running service container
  help        Help about any command
  info        Prints out information about kool setup (like environment variables)
  init        Initialize kool preset in the current working directory
  run         Runs a custom command defined at kool.yaml
  start       Start Kool environment containers
  status      Shows the status for containers
  stop        Stop kool environment containers

Flags:
  -T, --disable-tty   Disables TTY
  -h, --help          help for kool
  -v, --version       version for kool

Use "kool [command] --help" for more information about a command.
```
