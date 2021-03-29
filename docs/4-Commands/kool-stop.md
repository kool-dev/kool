## kool stop

Stop and destroy the service containers

### Synopsis

Stop and destroy running [SERVICE] containers started with the 'kool start' command. If no [SERVICE] is provided, all containers will be stopped.

```
kool stop [SERVICE...] [flags]
```

### Options

```
  -h, --help    help for stop
      --purge   Remove all persistent data from volume mounts on containers
```

### Options inherited from parent commands

```
      --verbose   increases output verbosity
```

### SEE ALSO

* [kool](kool)	 - Development environments made easy

