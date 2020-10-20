## kool completion

Generate completion script

### Synopsis

To load completions:

Bash:

$ source <(kool completion bash)

#### To load completions for each session, execute once:
Linux:
  $ kool completion bash > /etc/bash_completion.d/kool
MacOS:
  $ kool completion bash > /usr/local/etc/bash_completion.d/kool

Zsh:

#### If shell completion is not already enabled in your environment you will need to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

#### To load completions for each session, execute once:
$ kool completion zsh > "${fpath[1]}/_kool"

#### You will need to start a new shell for this setup to take effect.

Fish:

$ kool completion fish | source

#### To load completions for each session, execute once:
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

* [kool](kool.md)	 - kool - Kool stuff

