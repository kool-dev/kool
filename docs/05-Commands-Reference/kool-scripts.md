## kool scripts

List scripts defined in kool.yml

### Synopsis

List the scripts defined in kool.yml or kool.yaml in the current working directory
and in ~/kool. Use the optional FILTER to show only scripts that start with a given
prefix.

```
kool scripts [FILTER]
```

### Options

```
  -h, --help   help for scripts
      --json   Output scripts as JSON
```

When using `--json`, output is an array of objects with `name`, `comments`, and `commands`.

### Options inherited from parent commands

```
      --verbose              Increases output verbosity
  -w, --working_dir string   Changes the working directory for the command
```

### SEE ALSO

* [kool](kool)	 - Cloud native environments made easy
