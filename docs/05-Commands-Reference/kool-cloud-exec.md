## kool cloud exec

Execute a command inside a running service container deployed to Kool.dev Cloud

### Synopsis

After deploying an application to Kool.dev Cloud using 'kool deploy',
execute a COMMAND inside the specified SERVICE container (similar to an SSH session).
Must use a KOOL_API_TOKEN environment variable for authentication.

```
kool cloud exec SERVICE [COMMAND] [--] [ARG...]
```

### Options

```
  -c, --container string   Container target. (default "default")
  -h, --help               help for exec
```

### Options inherited from parent commands

```
      --domain string        Environment domain name to deploy to
      --token string         Token to authenticate with Kool.dev Cloud API
      --verbose              Increases output verbosity
  -w, --working_dir string   Changes the working directory for the command
```

### SEE ALSO

* [kool cloud](kool_cloud)	 - Interact with Kool.dev Cloud and manage your deployments.

