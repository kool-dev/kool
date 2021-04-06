## kool exec

Execute a new command inside a running service container

### Synopsis

This command allows to spawn a new process (specified by COMMAND) within a running service container (specified by SERVICE).

```
kool exec [OPTIONS] SERVICE COMMAND [--] [ARG...]
```

### Options

```
  -d, --detach            Detached mode: Run command in the background.
  -T, --disable-tty       Deprecated - no effect.
  -e, --env stringArray   Environment variables.
  -h, --help              help for exec
```

### Options inherited from parent commands

```
      --verbose   increases output verbosity
```

### SEE ALSO

* [kool](kool)	 - Development environments made easy

