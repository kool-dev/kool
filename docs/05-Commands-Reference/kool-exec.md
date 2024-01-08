## kool exec

Execute a command inside a running service container

### Synopsis

Execute a COMMAND inside the specified SERVICE container (similar to an SSH session).

```
kool exec [OPTIONS] SERVICE COMMAND [--] [ARG...]
```

### Options

```
  -d, --detach            Detached mode: Run command in the background.
  -e, --env stringArray   Environment variables.
  -h, --help              help for exec
```

### Options inherited from parent commands

```
      --verbose              Increases output verbosity
  -w, --working_dir string   Changes the working directory for the command
```

### SEE ALSO

* [kool](kool)	 - Cloud native environments made easy

