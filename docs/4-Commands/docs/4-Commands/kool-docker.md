## kool docker

Create a new container and run the command in it.

### Synopsis

This command acts as a helper for docker run.
You can start with options that go before the image name
for docker run itself, i.e --env='VAR=VALUE'. Then you must pass
the image name and the command you want to execute on that image.

```bash
kool docker [options] [image] [command] [flags]
```

### Options

```bash
  -T, --disable-tty           Deprecated - no effect
  -e, --env stringArray       Environment variables
  -h, --help                  help for docker
  -p, --publish stringArray   Publish a container’s port(s) to the host
  -v, --volume stringArray    Bind mount a volume
```

### Options Inherited from Parent Commands

```bash
      --verbose   increases output verbosity
```

### SEE ALSO

* [kool](kool)	 - kool - Kool stuff

