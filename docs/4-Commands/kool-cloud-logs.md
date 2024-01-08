## kool cloud logs

See the logs of running service container deployed to Kool Cloud

### Synopsis

After deploying an application to Kool Cloud using 'kool deploy',
you can see the logs from the specified SERVICE container.
Must use a KOOL_API_TOKEN environment variable for authentication.

```
kool cloud logs [OPTIONS] SERVICE
```

### Options

```
  -c, --container string   Container target. (default "default")
  -f, --follow             Follow log output.
  -h, --help               help for logs
  -t, --tail int           Number of lines to show from the end of the logs for each container. A value equal to 0 will show all lines. (default 25)
```

### Options inherited from parent commands

```
      --domain string        Environment domain name to deploy to
      --token string         Token to authenticate with Kool Cloud API
      --verbose              Increases output verbosity
  -w, --working_dir string   Changes the working directory for the command
```

### SEE ALSO

* [kool cloud](kool_cloud)	 - Interact with Kool Cloud and manage your deployments.

