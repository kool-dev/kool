## Installation

### Requirements

Kool is powered by **[Docker](https://docs.docker.com/get-docker/)** and **[Docker Compose](https://docs.docker.com/compose/install/)**, you need to have them installed on your machine.

### For Linux and MacOS

To install **kool** simply run the following script.

```bash
curl -fsSL https://raw.githubusercontent.com/kool-dev/kool/master/install.sh | sudo bash
```

All this script will do is download latest kool bin from https://github.com/kool-dev/kool/releases for your OS and put in your `/usr/local/bin` folder.

### For Windows

Download the installer [here](https://github.com/kool-dev/kool/releases)

To check if everything looks good simply run **kool** to see the list of available commands.

## Updating kool to a newer version

When a new version is released you can obtain it with the builtin command:

```bash
sudo kool self-update
```

This will check the latest realease and download/install it case there is a newer version.

**Alternative**: in order to get a new release you can always repeat the installation steps provided above as well, which should get you the latest version.

> We suggest that you *start* and sign up for new releases notifications on our Github main repository so you always will stay up to date with our latest features!

## Autocompletion

To load completions:

### Bash

`$ source <(kool completion bash)`

To load completions for each session, execute once:

Linux:
  `$ kool completion bash > /etc/bash_completion.d/kool`
MacOS:
  `$ kool completion bash > /usr/local/etc/bash_completion.d/kool`

### Zsh

**If shell completion is not already enabled in your environment you will need to enable it**.  You can execute the following once:

`$ echo "autoload -U compinit; compinit" >> ~/.zshrc`

To load completions for each session, execute once:

`$ kool completion zsh > "${fpath[1]}/_kool"`

**You will need to start a new shell for this setup to take effect**.

### Fish

`$ kool completion fish | source`

To load completions for each session, execute once:

`$ kool completion fish > ~/.config/fish/completions/kool.fish`
