## kool docker

Create a new container using the specified [image] and run a [command] inside it.

### Synopsis

This command acts as a helper for 'docker run'.
You can provide one or more [option...] before the [image] name that will be used
by 'docker run' itself (i.e --env='VAR=VALUE'). Then you must pass
the [image] name and the [command] you want to execute on that [image].

```
kool docker [option...] [image] [command] [flags]
```

### Options

```
  -T, --disable-tty           Deprecated - no effect.
  -e, --env stringArray       Environment variables.
  -h, --help                  help for docker
  -p, --publish stringArray   Publish a container's port(s) to the host.
  -v, --volume stringArray    Bind mount a volume.
```

### Options inherited from parent commands

```
      --verbose   increases output verbosity
```

### SEE ALSO

* [kool](kool)	 - kool - Kool stuff

