# Start a monorepo with NestJS and NextJS with Docker

> Note: this preset bootstraps both NestJS and NextJS with **Typescript** and `npm` as the package manager.

1. Run `kool create mono-ts-nest-next my-project`
2. Run `kool start`
3. Give a few moments for both Nest and Next to build/start (you can check out the outputs with `kool logs -f`)
4. After built and running you can access:
	- `http://localhost:81` - NestJS backend
	- `http://localhost` - NextJS frontend

> Yes, using **kool** + Docker to create and work on new monorepo for NestJS and NextJS projects is that easy!

## Requirements

If you haven't done so already, you first need to [install Docker and the kool CLI](/docs/getting-started/installation).

Also, make sure you're running the latest version of **kool**. Run the following command to compare your local version of **kool** with the latest release, and, if a newer version is available, automatically download and install it.

```bash
$ kool self-update
```

> Please note that it helps to have a basic understanding of how Docker and Docker Compose work to use Kool with Docker.

## Extra steps

If you are using Database or Cache services for your NestJS backend you need to setup some environment variables. [Check out our NestJS preset documentation]().

## Resources

### Run Commands in Docker Containers

Use [`kool exec`](/docs/commands/kool-exec) to execute a command inside a running service container:

```bash
# kool exec [OPTIONS] SERVICE COMMAND [--] [ARG...]

$ kool exec backend node -v # backend is the NestJS container service
$ kool exec frontend node -v # backend is the NextJS container service
```

Try `kool run nest --help` to execute the `kool exec backend nest --help` command in your running `app` container and print out information about NestJS' commands.

### Open Sessions in Docker Containers

Similar to SSH, if you want to open a Bash session in your `backend` or `frontend` containers, run `kool exec <container> bash`, where `<container>` is the name of the service container in **docker-compose.yml**. If you prefer, you can use `sh` instead of `bash` (`kool exec <service> sh`).

```bash
$ kool exec backend bash # opens a shell on the container running NestJS
$ kool exec frontend bash # opens a shell on the container running NextJS
```

## Related Presets

We have more presets to help you start projects with **kool** in a standardized way across different frameworks.

- **[Next.js](/docs/2-Presets/NextJS.md)**
- **[NestJS](/docs/2-Presets/NestJS.md)**
- **[AdonisJs](/docs/2-Presets/AdonisJs.md)**
- **[Nuxt.js](/docs/2-Presets/NuxtJS.md)**
