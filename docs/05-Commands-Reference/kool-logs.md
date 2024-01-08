## kool logs

Display log output from running service containers

### Synopsis

Display log output from all running service containers,
or one or more specified [SERVICE...] containers. Add a '-f' option to the
the command to follow the log output (i.e. 'kool logs -f [SERVICE...]').

```
kool logs [OPTIONS] [SERVICE...]
```

### Options

```
  -f, --follow     Follow log output.
  -h, --help       help for logs
  -t, --tail int   Number of lines to show from the end of the logs for each container. A value equal to 0 will show all lines. (default 25)
```

### Options inherited from parent commands

```
      --verbose              Increases output verbosity
  -w, --working_dir string   Changes the working directory for the command
```

### SEE ALSO

* [kool](kool)	 - Cloud native environments made easy

