## kool completion

Generate shell completion configuration script

### Synopsis

Autocompletion:

If you want to use kool autocompletion in your Unix shell, follow the appropriate instructions below.

After running one of the below commands, remember to start a new shell for autocompletion to take effect.

#### Bash

Temporarily enable autocompletion for your current session only:

$ source <(kool completion bash)

Permanently enable autocompletion for all sessions:

  Linux:

  $ kool completion bash > /etc/bash_completion.d/kool

  macOS:

  $ kool completion bash > /usr/local/etc/bash_completion.d/kool

#### Zsh

If Zsh tab completion is not already initialized on your machine, run the following command to turn it on.

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

Permanently enable autocompletion for all sessions:

$ kool completion zsh > "${fpath[1]}/_kool"

#### Fish

Temporarily enable autocompletion for your current session only:

$ kool completion fish | source

Permanently enable autocompletion for all sessions:

$ kool completion fish > ~/.config/fish/completions/kool.fish


```
kool completion [bash|zsh|fish|powershell]
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --verbose   increases output verbosity
```

### SEE ALSO

* [kool](kool)	 - Cloud native environments made easy

