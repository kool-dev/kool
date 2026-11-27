## kool cloud

Interact with Kool.dev Cloud and manage your deployments.

### Synopsis

The cloud subcommand encapsulates a set of APIs to interact with Kool.dev Cloud and deploy, access and tail logs from your deployments.

Authentication: the API token is read from the KOOL_API_TOKEN environment variable.
Prefer it over the --token flag, whose value shows up in the process list (ps) and in your shell history.

### Examples

```
kool cloud deploy
```

### Options

```
      --domain string   Environment domain name to deploy to
  -h, --help            help for cloud
      --token string    Token to authenticate with Kool.dev Cloud API. Prefer the KOOL_API_TOKEN environment variable, read when --token is not given: a --token value shows up in the process list
```

### Options inherited from parent commands

```
      --output string        Output format: json (machine-readable). For 'kool run', place before the script name.
      --verbose              Increases output verbosity
  -w, --working_dir string   Changes the working directory for the command
```

### SEE ALSO

* [kool](kool)	 - Cloud native environments made easy
* [kool cloud deploy](kool_cloud_deploy)	 - Deploy a local application to a Kool.dev Cloud environment
* [kool cloud destroy](kool_cloud_destroy)	 - Destroy an environment deployed to Kool.dev Cloud
* [kool cloud exec](kool_cloud_exec)	 - Execute a command inside a running service container deployed to Kool.dev Cloud
* [kool cloud logs](kool_cloud_logs)	 - See the logs of running service container deployed to Kool.dev Cloud
* [kool cloud setup](kool_cloud_setup)	 - Set up local configuration files for deployment

