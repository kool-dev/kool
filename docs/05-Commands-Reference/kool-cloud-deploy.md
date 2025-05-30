## kool cloud deploy

Deploy a local application to a Kool.dev Cloud environment

```
kool cloud deploy
```

### Options

```
      --domain-extra stringArray   List of extra domain aliases
  -h, --help                       help for deploy
      --platform string            Platform for docker build (default: linux/amd64) (default "linux/amd64")
      --timeout uint               Timeout in minutes for waiting the deployment to finish
      --www-redirect               Redirect www to non-www domain
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

