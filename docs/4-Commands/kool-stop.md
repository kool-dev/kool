## kool stop

Stop and destroy running service containers

### Synopsis

Stop and destroy the specified [SERVICE] containers, which were started
using 'kool start'. If no [SERVICE] is provided, all running containers are stopped.

```
kool stop [SERVICE...]
```

### Options

```
  -h, --help    help for stop
      --purge   Remove all persistent data from volume mounts on containers
```

### Options inherited from parent commands

```
      --verbose              Increases output verbosity
  -w, --working_dir string   Changes the working directory for the command
```

### SEE ALSO

* [kool](kool)	 - Cloud native environments made easy

