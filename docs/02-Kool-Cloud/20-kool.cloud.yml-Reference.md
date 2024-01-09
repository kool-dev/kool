`kool.cloud.yml` will hold all the extra configuration needed to move your application from `docker-compose.yml` and run it in the cloud.

## The basics

The `kool.cloud.yml` file is an extension of your already familiar `docker-compose.yml`, having the same basic structure but introducing some configuration entries to enable you to fine-tune your deployment container needs.

Suppose you have the following `docker-compose.yml` file:

```yaml
services:
  app:
    image: some/image
      ports:
        - 80:80 # maps the container port 80 to your localhost
```

Now, if you want to deploy this single-container app to the cloud using Kool, you need the following `kool.cloud.yml` file:

```yaml
services:
  app:
    public:
      - port: 80
```

Provided you have already signed up and obtained your access token for Kool Cloud in your `.env`, simply by running `kool cloud deploy`, you will get your container with `some/image` deployed to the cloud and a URL with HTTPS that will route incoming traffic to port 80 of that container.

## Reference

### Full example

Here's an example of a `kool.cloud.yml` file showcasing all the features and configuration entries available:

```yaml
services:
  app:
    # Applications usually will require a special image built for deployment.
    # Reference: https://docs.docker.com/compose/compose-file/compose-file-v3/#build
    build: .

    # Tells Kool Cloud that this service is accessible through the deployment URL.
    # Note: only one service can be set to be public.
    public: true # simply defining true is enough to most cases where your `expose` port will be used for routing incoming HTTP requests.

    # Another option is advanced definition:
    public:
      # Tells the port that should be used for routing incoming HTTP traffic.
      - port: 80
      # Sometimes you may have a second process you want to access as well, i.e. some
      # websocket service that you spin up via a 'daemon' and listens on another port.
      # You may do so by specifying a second port with a path - so all requests starting
      # with such path prefix will be routed to that port instead of the default port above.
      - port: 3000
        path: /ws

    # Tells what port the app will listen to (optional).
    expose: 80

    # Tells your app's root folder so all other paths can be relative (optional).
    root: /app

    # Containers are ephemeral, that means their filesystem do not persist across deployments.
    # If you want to persist stuff into the disk across deployments, you can do so by defining persistent paths here.
    persists:
      # Total size of the volume you want to attach to the running container.
      size: 1Gi
      # Paths to persist - within that single volume, you can have one or more paths
      # that are going to be mounted every time your containers are running. Note that
      # such mounts include before/after hooks as well as daemon containers.
      paths:
        # The path within the container. Must be either aboslute or relative to the 'root' config.
        - path: /app/some/path/persisted
          # Tells the Deploy API to sync the folder from your built image to the persisted storage.
          # This is very helpful to start off with a known folder structure.
          sync: true
          # Tells what user and group should own the persisted folder (only used along the sync: true)
          chown: user:group

    # By default, Kool Cloud will rollout new deployments in a blue-green fashion.
    # If you want to disable it and make sure the current running container
    # is stopped before the new one is created, set 'recreate: true'.
    recreate: false

    # Sometimes you may make changes to your app that wouldn't necessarily trigger
    # new containers to be created by the kool cloud deploy process.
    #
    # For example, if you change only environment variables or use a fixed image tag
    # like ':latest' that doesn't change.
    #
    # By setting `force: true`, you tell the API to always update this service.
    force: false

    # Here we can define processes that behave like services and must be run in the cloud only.
    # Usually, this is very helpful for queue workers and that kind of stuff. The processes will
    # run as other containers in your deployment using the very same image as the main service.
    daemons:
      - command: [ start-queues, arg1, arg2 ]
      # As mentioned above in the 'public' config, you may also want to have a daemon
      # that serves specific requests. A common use case is WebSocket services. You can
      # define that your daemon exposes a TCP port with the 'expose' entry.
      - command: [ run-websocker-server, --port=3000 ]
        expose: 3000

    # Hooks
    #
    # It's possible that you want to run some extra steps either before or after every
    # time your application is deployed. Such hooks are executed in standalone containers
    # using the same Docker image about to be deployed.

    #
    # The 'before' hook is a special section where we can define commands to be executed
    # right before a new deployment happens.
    before:
      - script_to_run.sh
    # The 'after' hook is a special section where we can define procedures to be executed
    # right after a new deployment finishes.
    after:
      - [ run-database-migrations, arg1, arg2 ]

    # The 'cron' config allows you to set scheduled actions to run in a Cron job fashion.
    cron:
      - { schedule: "* * * * *", command: [ run-some-task-every-minute ] }

    # Environment Variables to the cloud container
    #
    # You usually will want to provide a file containing the environment variables
    # that your deployment should have. Such a file may contain special Kool variables
    # that will be translated to their actual value by the Deploy API.
    #
    # There are two ways to provide such a file.
    #
    # 'environment' refers to a file with all of the environment variables available.
    # The Deploy API will get the contents of this file and make all of them available
    # as true environment variables for the container. You can check them in your deployed
    # container via kool cloud exec env (assuming you have env available).
    environment: kool.deploy.env
    #
    # 'env' is a different option that allows you to build a file inside your running container.
    # This is useful for apps that have a .env file, but you do not want to have that built into your app Docker image.
    env:
      source: kool.deploy.env
      target: .env
```
