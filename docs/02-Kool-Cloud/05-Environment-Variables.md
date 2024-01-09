Most applications and frameworks nowadays rely on environment variables to configure important aspects of their functions, mainly providing credentials and other secrets your app needs to work and access other resources.

Kool.dev Cloud supports a few different ways you can define your environment variables for a deploying container, so pick the one that best suits you.

### Using `kool.deploy.env` file for deploy

`kool.deploy.env` should be a `.env` formatted file. You can point to it like this:

```yaml
services:
  app:
    # ...
    environment: kool.deploy.env
```

Upon deployment, all of the variables within that file will be parsed, placeholders replaced—if you have any—and then **each variable will become a real environment variable in the running container**.

This option is usually best suited for automated CI routines since you work your way to have a different `kool.deploy.env` file for each of your deploying environments (i.e., staging and production).

### Using a plain YAML object for environment variables

```yaml
services:
  app:
    # ...
    environment:
      FOO: bar
```

You can simply use a YAML map of values that will become your environment variables in the running deployed container. This is handy sometimes when you have simple and not sensitive variables you want to add to a container for deploy.

### Build a `.env` file inside the running container

If your application does rely on and requires a `.env` file existing in the running container, you may achieve so by using the `env:` entry:

```yaml
services:
  app:
    # ...

    # 'env' is a different option that allows you to build a file inside your running container.
    env:
      source: kool.deploy.env
      target: .env
```

This is useful for apps that require the .env file, but you do not wish to have that built into your Docker image itself.
