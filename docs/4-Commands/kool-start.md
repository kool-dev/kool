## kool start

Start service containers defined in docker-compose.yml

### Synopsis

Start one or more specified [SERVICE] containers. If no [SERVICE] is provided,
all containers are started. If the containers are already running, they are recreated.

```
kool start [SERVICE...]
```

### Options

```
  -f, --foreground       Start containers in foreground mode
  -h, --help             help for start
      --profile string   Specify a profile to enable
  -b, --rebuild          Updates and builds service's images
```

### Options inherited from parent commands

```
      --verbose   increases output verbosity
```

### SEE ALSO

* [kool](kool)	 - Cloud native environments made easy

