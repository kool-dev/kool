You have the ability to define hooks to run before or after every deploy you make.

These hooks will run using the very same image that is being deployed. This is usually needed for common routines, as for example running database migrations, sending alerts of updates to some other service, etc.

To illustrate the options you have on `kool.cloud.yml` file:

```yaml
services:
  app:
    # ...

    # The 'before' hook is a special section where we can define commands to be executed
    # right before a new deployment happens.
    before:
      - [ sh, run-database-migrations.sh, arg1, arg2 ]

    # The 'after' hook is a special section where we can define procedures to be executed
    # right after a new deployment finishes.
    after:
      - [ sh, run-cache-version-update.sh, arg1, arg2 ]
```

### Failures on lifecycle hooks

Please notice that these lifecycle hooks are required in order for the new deploy to be successfull - this mean that **if any of them fail** - either `before` or `after` new deployed contianer versions are running - **the whole deploy is going to be rolledback**. As you can imagine this poses a challange specially on database migrations since they can be problematic and not backwards compatible with previously running container version.
