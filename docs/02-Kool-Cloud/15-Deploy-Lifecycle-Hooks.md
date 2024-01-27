# Hooks

You have the ability to define hooks to run before or after every deploy you make.

These hooks will run using the very same image that is being deployed. This is usually needed for common routines, such as running database migrations, sending alerts of updates to some other service, etc.

# Before and After

To illustrate the options you have in the `kool.cloud.yml` file:

```yaml
services:
  app:
    # ...

    # The 'before' hook is a special section where we can define commands to be executed
    # right before a new deployment happens.
    # ATTENTION: current limitation - can only have a 'before' hook after a first deploy has created the environment.
    before:
      - [ sh, run-database-migrations.sh, arg1, arg2 ]

    # The 'after' hook is a special section where we can define procedures to be executed
    # right after a new deployment finishes.
    after:
      - [ sh, run-cache-version-update.sh, arg1, arg2 ]
```

## Failures on lifecycle hooks

Please notice that these lifecycle hooks are required for the new deploy to be successful—this means that **if any of them fail**—either `before` or `after` newly deployed container versions are running—**the whole deploy is going to be rolled back**. As you can imagine, this poses a challenge, especially on database migrations since they can be problematic and not backwards compatible with the previously running container version.

## Limitations

The `before` hook can only be used after a first deploy has succeded for your environment - that is a limitation that should be lifted in the future, but currently can halt your first deploy if it includes a `before` hook.
