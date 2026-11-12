## kool create

Create a new project using a preset

### Synopsis

Create a new project using the specified PRESET in a directory named FOLDER.

Questions asked by the preset can be answered beforehand with --answer REF=value,
where REF is the question ref (or the question text when it has no ref). Options
can be given by value or label, and multiple choices are comma separated. Without
a terminal, or with --no-interaction, a question without an answer fails.

```
kool create PRESET FOLDER
```

### Options

```
  -a, --answer stringArray   Answer to a preset question as REF=value (repeatable; comma separate multiple choices)
  -h, --help                 help for create
  -n, --no-interaction       Fail on any preset question without an answer instead of prompting
```

### Options inherited from parent commands

```
      --output string        Output format: json (machine-readable). For 'kool run', place before the script name.
      --verbose              Increases output verbosity
  -w, --working_dir string   Changes the working directory for the command
```

### SEE ALSO

* [kool](kool)	 - Cloud native environments made easy

