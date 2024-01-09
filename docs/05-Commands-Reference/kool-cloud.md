## kool cloud

Interact with Kool.dev Cloud and manage your deployments.

### Synopsis

The cloud subcommand encapsulates a set of APIs to interact with Kool.dev Cloud and deploy, access and tail logs from your deployments.

### Examples

```
kool cloud deploy
```

### Options

```
      --domain string   Environment domain name to deploy to
  -h, --help            help for cloud
      --token string    Token to authenticate with Kool.dev Cloud API
```

### Options inherited from parent commands

```
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

