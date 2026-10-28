## kool preset

Install configuration files customized for Kool in the current directory

### Synopsis

Initialize a project using the specified [PRESET] by installing configuration
files customized for Kool in the current working directory. If no [PRESET] is provided,
an interactive wizard will present the available options.

Questions asked by the preset can be answered beforehand with --answer REF=value,
where REF is the question ref (or the question text when it has no ref). Options
can be given by value or label, and multiple choices are comma separated. Without
a terminal, or with --no-interaction, a question without an answer fails.

```
kool preset [PRESET]
```

### Options

```
  -a, --answer stringArray   Answer to a preset question as REF=value (repeatable; comma separate multiple choices)
  -h, --help                 help for preset
  -n, --no-interaction       Fail on any preset question without an answer instead of prompting
```

### Options inherited from parent commands

```
      --verbose              Increases output verbosity
  -w, --working_dir string   Changes the working directory for the command
```

### SEE ALSO

* [kool](kool)	 - Cloud native environments made easy

