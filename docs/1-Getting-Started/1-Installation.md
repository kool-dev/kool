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

To check if everything looks good simply run **kool** to see the list of available commands.

## Autocompletion

To load completions:

### Bash

```
$ source <(kool completion bash)
```

To load completions for each session, execute once:
Linux:
```
  $ kool completion bash > /etc/bash_completion.d/kool
```
MacOS:
```
  $ kool completion bash > /usr/local/etc/bash_completion.d/kool
```

Attention: MacOS bash doesn't come with completion enabled by default, you need to look into enabling it.

### Zsh

```
# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ kool completion zsh > "${fpath[1]}/_kool"

# You will need to start a new shell for this setup to take effect.
```


### Fish

```
$ kool completion fish | source

# To load completions for each session, execute once:
$ kool completion fish > ~/.config/fish/completions/kool.fish
```
