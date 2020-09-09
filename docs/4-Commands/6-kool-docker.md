## kool docker

Creates a new container and runs the command in it.

### Synopsis

This command acts as a helper for docker run.
You can start with options that go before the image name
for docker run itself, i.e --env='VAR=VALUE'. Then you must pass
the image name and the command you want to execute on that image.

```
kool docker [options] [image] [command] [flags]
```

### Options

```
  -T, --disable-tty           Disables TTY
  -e, --env stringArray       Environment variables
  -h, --help                  help for docker
  -p, --publish stringArray   Publish a container’s port(s) to the host
  -v, --volume stringArray    Bind mount a volume
```

### SEE ALSO

* [kool](kool.md)	 - kool - Kool stuff

