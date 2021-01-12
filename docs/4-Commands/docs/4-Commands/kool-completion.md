## kool completion

Generate completion script.

### Synopsis

To load completions:

**Bash:**

```bash
$ source <(kool completion bash)
```

To load completions for each session, execute once:

**Linux:**

```bash
  $ kool completion bash > /etc/bash_completion.d/kool
```

**MacOS:**

```bash
  $ kool completion bash > /usr/local/etc/bash_completion.d/kool
```

**Zsh:**

If shell completion is not already enabled in your environment, you will need to enable it. You can execute the following once:

```bash
$ echo "autoload -U compinit; compinit" >> ~/.zshrc
```

To load completions for each session, execute once:

```bash
$ kool completion zsh > "${fpath[1]}/_kool"
```

You will need to start a new shell for this setup to take effect.

**Fish:**

```bash
$ kool completion fish | source
```

To load completions for each session, execute once:

```bash
$ kool completion fish > ~/.config/fish/completions/kool.fish
```


```bash
kool completion [bash|zsh|fish|powershell]
```

### Options

```bash
  -h, --help   help for completion
```

### Options Inherited from Parent Commands

```bash
      --verbose   increases output verbosity
```

### SEE ALSO

* [kool](kool)	 - kool - Kool stuff

