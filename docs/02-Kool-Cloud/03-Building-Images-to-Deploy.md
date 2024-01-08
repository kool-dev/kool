This topic is usually the biggest source or problems and trial and error frustrations when deploying cloud native applications for the first time.

As much as the Kool.dev project and the whole community tries to help and facilitate container images building, it is a times ultimately an individual and singular process for your web application.

That being said there's no scape from having some knowledge on how to properly build your images to deploy your app to the cloud - or at least seek such knowledgable hands to assist your in this moment.

For the most basic cases - like if you are using one of our presets - you will have a great start point by using our utility along `kool cloud setup` - this command will inquiry you about basic options on building your container images.

### `kool cloud deploy` building images

`kool` CLI is going to handle the build of your images locally - in your own host system. That means it's required that the environment where you are going to run `kool cloud deploy` have a working Docker-like engine running that can process successfully a `docker build ...` command.

The syntax configuration for building your deploy image for a given service on `kool.cloud.yml` is the very same as you use it locally on `docker-compose.yml`:

Check out the [Docker Compose `build` Documentation](https://docs.docker.com/compose/compose-file/compose-file-v3/#build) for reference.

```yaml
services:
  app:
    # ...
    build: . # this uses the root folder as context, and expects a Dockerfile to exist on it
```

or

```yaml
services:
  app:
    # ...
    build:
      context: ./dir # changes the context folder
      dockerfile: Dockerfile-alternate # name a different file than default 'Dockerfile'
      args:
        buildno: 1 # define values for ARGS used in your Dockerfile
```

Your image will be built locally when running the `kool` CLI for a deploy, and then pushed securely to Kool.dev Cloud registry to a repository dedicated to your app environment.

### Using a Private Registry

You may already have or use your own private registry for handling images. You are welcome to hold the build process apart from the `kool cloud deploy` step, and just use the already built images in your `kool.cloud.yml` file:

```yaml
services:
  app:
    # ...
    image: myrepo-registry/my-built-image
```

If that registry is private you need to provide Kool.dev Cloud with credentials to read from that repo. As this is not yet fully automated you can [ping us via email to `contact@kool.dev`](contact@kool.dev) to set it up for you.
