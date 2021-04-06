## kool docker

Create a new container (a powered up 'docker run')

### Synopsis

This command acts as a helper for 'docker run'.
You can provide one or more [OPTIONS] before the IMAGE name that will be used
by 'docker run' itself (i.e --env='VAR=VALUE'). Then you must pass
the IMAGE name and the [COMMAND] you want to execute on that IMAGE. After that you can use -- and follow with any extra arguments that command may require.

```
kool docker [OPTIONS] IMAGE [COMMAND] -- [ARG...] [flags]
```

### Options

```
  -T, --disable-tty           Deprecated - no effect.
  -e, --env stringArray       Environment variables.
  -h, --help                  help for docker
  -n, --network stringArray   Connect a container to a network.
  -p, --publish stringArray   Publish a container's port(s) to the host.
  -v, --volume stringArray    Bind mount a volume.
```

### Options inherited from parent commands

```
      --verbose   increases output verbosity
```

### SEE ALSO

* [kool](kool)	 - Development environments made easy

