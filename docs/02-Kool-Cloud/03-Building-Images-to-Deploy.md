This topic is usually the biggest source of problems and trial-and-error frustrations when deploying cloud-native applications for the first time.

As much as the Kool.dev project and the whole community try to help and facilitate container image building, it is at times ultimately an individual and singular process for your web application.

That being said, there's no escape from having some knowledge of how to properly build your images to deploy your app to the cloud—or at least seek such knowledgeable hands to assist you in this moment.

For the most basic cases — like if you are using one of our presets — you will have a great starting point by using our utility along with `kool cloud setup`. This command will inquire about basic options for building your container images.


### `kool cloud deploy` building images

The `kool` CLI is going to handle the build of your images locally—in your own host system. That means it's required that the environment where you are going to run `kool cloud deploy` has a working Docker-like engine running that can successfully process a `docker build ...` command.

The syntax configuration for building your deploy image for a given service in `kool.cloud.yml` is the very same as you use it locally in `docker-compose.yml`:

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

Your image will be built locally when running the `kool` CLI for a deploy and then pushed securely to the Kool.dev Cloud registry to a repository dedicated to your app environment.

### Build arguments (`args`) variables

As stated above the args provided to Docker when building the image will come from the `services.<service>.build.args` configuration entry.

It's a common need to have different values for different environments (i.e staging vs production). Kool Cloud supports two different ways for you to have a single `kool.cloud.yml` definition and still use different values per environment:

- **Environment Variables**: we will parse environment variables before passing the build args to Docker so you can use the common syntax `FOO: "$FOO"` and the value will be interpolated to the current value of `FOO` when running the `kool cloud deploy` command.
- **Kool Cloud Environment Variables**: you can define the variables under the Environment on Kool.dev Cloud panel, and then use the special `FOO: {{FOO}}` syntax to have the value interpolated with the web panel managed value for `FOO`.

#### Escaping variables

If the value you want to use contains `$` sign that could lead to trouble on having the sign mistakenly parsed as a variable marker. To scape it you need to double it: `FOO=$$bar` will have the desired effect of getting the actual `$bar` string as the value of `FOO` environment variable.

### Using a Private Registry

You may already have or use your own private registry for handling images. You are welcome to hold the build process apart from the `kool cloud deploy` step and just use the already built images in your `kool.cloud.yml` file:

```yaml
services:
  app:
    # ...
    image: myrepo-registry/my-built-image
```

If that registry is private, you need to provide Kool.dev Cloud with credentials to read from that repo. As this is not yet fully automated, you can [contact us via email at `contact@kool.dev`](contact@kool.dev) to set it up for you.

