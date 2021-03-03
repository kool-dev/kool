## Installation

### Requirements

**kool** is powered by **Docker**. To use **kool**, you need to **[install the Docker Engine](https://docs.docker.com/get-docker/)** and **[Docker Compose](https://docs.docker.com/compose/install/)** on your machine, and make sure they're both running.

> Docker Compose is included with **Docker Desktop for Mac** and **Docker Desktop for Windows**.

### For Linux and macOS

Install **kool** by running the following script. It will download the latest **kool** binary from [https://github.com/kool-dev/kool/releases](https://github.com/kool-dev/kool/releases), and save it in your `/usr/local/bin` folder.

```bash
$ curl -fsSL https://kool.dev/install | bash
```

### For Windows

Install **kool** by downloading the appropriate installer from [https://github.com/kool-dev/kool/releases](https://github.com/kool-dev/kool/releases). At the bottom of the release notes for the latest version, expand the list of "Assets", and download the installer that matches your machine.

### Verification

Verify **kool** is installed correctly by running `kool` in your terminal to see a list of available commands.

## Update to a Newer Version

Update **kool** using the `self-update` command. This command will compare your local version with the latest release, and, if a newer version is available, automatically download and install it.

```bash
$ kool self-update
```

> Sign up for new release notifications and stay up-to-date on our latest features! [Go to our main GitHub repository](https://github.com/kool-dev/kool), click on "Watch" at the top, select the "Custom" option, check "Releases", and hit Apply.

## Autocompletion

If you want to use **kool** autocompletion in your Unix shell, follow the appropriate instructions below.

### Bash

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

### Zsh

If Zsh tab completion is not already initialized on your machine, run the following command to turn it on.

```bash
$ echo "autoload -U compinit; compinit" >> ~/.zshrc
```

Permanently enable autocompletion for **all sessions**:

```bash
$ kool completion zsh > "${fpath[1]}/_kool"
```

> After running the above command, remember to start a new shell for autocompletion to take effect.

### Fish

Temporarily enable autocompletion for your **current session only**:

```bash
$ kool completion fish | source
```

Permanently enable autocompletion for **all sessions**:

```bash
$ kool completion fish > ~/.config/fish/completions/kool.fish
```

> After running the above command, remember to start a new shell for autocompletion to take effect.
