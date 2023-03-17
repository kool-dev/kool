## kool docker

Create a new container (a powered up 'docker run')

### Synopsis

A helper for 'docker run'. Any [OPTIONS] added before the
IMAGE name will be used by 'docker run' itself (i.e. --env='VAR=VALUE').
Add an optional [COMMAND] to execute on the IMAGE, and use [--] after
the [COMMAND] to provide optional arguments required by the COMMAND.

```
kool docker [OPTIONS] IMAGE [COMMAND] [--] [ARG...]
```

### Options

```
  -e, --env stringArray       Environment variables.
  -h, --help                  help for docker
  -n, --network stringArray   Connect a container to a network.
  -p, --publish stringArray   Publish a container's port(s) to the host.
  -v, --volume stringArray    Bind mount a volume.
```

### Options inherited from parent commands

```
      --verbose              Increases output verbosity
  -w, --working_dir string   Changes the working directory for the command
```

### SEE ALSO

* [kool](kool)	 - Cloud native environments made easy

