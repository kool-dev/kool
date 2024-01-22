# Installation

## Requirements

**kool** is powered by **Docker**. To use **kool**, you need to **[install the Docker Engine and Docker Compose](https://docs.docker.com/get-docker/)** on your local machine, and make sure they're both running.

Important to notice: `kool` relies on the `docker` and `docker compose` CLI commands - that being said, you can pick any other Docker-compatibable container engine and use it seamlessly, like OrbStack for example.

> `kool` now requires [Docker Compose V2](https://docs.docker.com/compose/install/), make sure you have it enabled in your system.

## For Linux and macOS

Install **kool** by running the following script. It will download the latest **kool** binary from [https://github.com/kool-dev/kool/releases](https://github.com/kool-dev/kool/releases), and save it in your `/usr/local/bin` folder.

```bash
$ curl -fsSL https://kool.dev/install | bash
```

## For Windows

You must run `kool` on Windows via [WSL - Windows Subsystem for Linux](https://learn.microsoft.com/en-us/windows/wsl/install) - once you have a WSL environment properly set up, make sure you have [Docker available on it](https://docs.docker.com/desktop/wsl/), then you can install the CLI as you would in any Linux or MacOS (see above).

### Verification

Verify **kool** is installed correctly by running `kool` in your terminal to see a list of available commands.

# Stay Up-to-Date

Update **kool** to a newer version using the `self-update` command. This command will compare your local version with the latest release, and, if a newer version is available, automatically download and install it.

```bash
$ kool self-update
```

> Sign up for new release notifications and stay up-to-date on our latest features! [Go to our main GitHub repository](https://github.com/kool-dev/kool), click on "Watch" at the top, select the "Custom" option, check "Releases", and hit Apply.

# Autocompletion

If you want to use **kool** autocompletion in your Unix shell, follow the appropriate instructions below.

## Bash

Temporarily enable autocompletion for your **current session only**:

```bash
$ source <(kool completion bash)
```

Permanently enable autocompletion for **all sessions**:

Linux:

```bash
$ kool completion bash > /etc/bash_completion.d/kool
```

macOS:

```bash
$ kool completion bash > /usr/local/etc/bash_completion.d/kool
```

> After running one of the above commands, remember to start a new shell for autocompletion to take effect.

## Zsh

If Zsh tab completion is not already initialized on your machine, run the following command to turn it on.

```bash
$ echo "autoload -U compinit; compinit" >> ~/.zshrc
```

Permanently enable autocompletion for **all sessions**:

```bash
$ kool completion zsh > "${fpath[1]}/_kool"
```

> After running the above command, remember to start a new shell for autocompletion to take effect.

## Fish

Temporarily enable autocompletion for your **current session only**:

```bash
$ kool completion fish | source
```

Permanently enable autocompletion for **all sessions**:

```bash
$ kool completion fish > ~/.config/fish/completions/kool.fish
```

> After running one of the above commands, remember to start a new shell for autocompletion to take effect.
