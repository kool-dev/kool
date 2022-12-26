`kool.deploy.yml` will hold all the extra configuration needed to move your application from `docker-compose.yml` to run in the cloud.

## The basics

The `kool.deploy.yml` file is a an extension to your already familiar `docker-compose.yml`, having the same basic structure but introducing some configuration entries to enable you to fine tweak your deployment container needs.

Imagine you have the following `docker-compose.yml` file:

```yaml
services:
  app:
    image: some/image
      ports:
        - 80:80 # maps the container port 80 into your localhost
```

Now if you want to have this single-container app deployed to the Cloud using Kool you need the following `kool.deploy.yml` file:

```yaml
services:
  app:
    public:
      - port: 80
```

Provided you have already signed up and got you access token to Kool Cloud in your `.env`, simply by running `kool cloud deploy` you will get your container with `some/image` deployed to the cloud and a URL with HTTPS that will route the incoming traffic to the port 80 of such container.

## Reference

## Full example

Here's an example of a feature-rich `kool.deploy.yml` file showcasing all the features and configuration entries you have available:

```yaml
# Here we can tweak service configurations from docker-compose.yml
# to better suite our cloud deployments only.
services:
  app:
    # Appications usually will require a special image built for deployment.
    build: Dockerfile

    # Tells Kool Cloud that this service is accessible through the deployment URL.
    # Notice: only one service can be set to be public.
    public:
      # Tell the port which should be used for routing the HTTP traffic.
      - port: 80
      # Sometimes you may have a second process you want to access as well
      # i.e some websocket service that you spin via a 'daemon' and listens
      # on another port. You may do so by specifying a second port with a
      # `path` - so all requests starting with such path prefix will be routed
      # to that port instead of the default port above.
      - port: 3000
        path: /ws

    # What port the app will listen to receive incoming traffic upon deployment.
    port: 80

    # Optional: inform your app root folder so all other paths can be relative.
    root: /app

    # Containers are ephemeral! So is their local disks. If you want to persist stuff
    # into the disk across deployments, you can do so by defining persistent paths here.
    persists:
      # Total size of the volume you want to attach to the running container.
      size: 1Gi
      # Paths to persist - within that single volume, you can have one of more paths
      # that are going to be mounted everytime your containrs are running. Note that
      # such mounts include before/after hooks as well as daemon containers.
      paths:
        # Path must be either aboslute within the containers or relative to the 'root' config.
        - path: storage
          # Tells the Deploy API to sync the folder from your built image to the persisted
          # storage. This is very helpful to start of with a known folder structure.
          sync: true
          # Tell which user/group should ownthat folder (only used along the sync: true)
          chown: kool:kool

    # By default Kool Cloud will rollout new deployments in a blue-green fashion.
    # In case you want to disable it and make sure the current running container
    # is stopped before the new one is created, set the 'recreate: true'.
    recreate: false

    # Sometimes you may do changes to your app that wouldn't necessarily trigger
    # new containers to be created by the `kool cloud deploy` process.
    #
    # i.e: if you change only environment variables, or you use a fixed image
    # tag like ':latest' that doesn't change.
    #
    # By setting the `force: true` you tell the API to always update this service.
    force: false

    # Here we can define processes that behave like services and must be run in the cloud only.
    # Usually this is very helpful for queue workers and that kind of stuff. The processes will
    # run as other containers in your deployment using the very same image as the main service.
    daemons:
      - command: [ start-queues, arg1, arg2 ]
      # As mentioned above in the 'public' config, you may also want to have a daemon that serves
      # specific requests, a common use case is WebSocket services. You can define that your daemon
      # exposes a TCP port with the 'expose' entry.
      - command: [ run-websocker-server, --port=3000 ]
        expose: 3000

    # Hooks
    #
    # It's possible that you want to run some extra steps either before or after
    # everytime your application is deployed.
    #
    # The 'before' hook is a special section where we can define commands to be executed
    # right before a new deployment happens.
    before:
      - script_to_run.sh
    # The 'after' hook is a special section where we can define procedures to be executed
    # right after a new deployment finished.
    after:
      - [ run-database-migrations, arg1, arg2 ]

    # The 'cron' config allows you to set scheduled actions to run in a Cronjob fashion.
    cron:
      - { schedule: "* * * * *", command: [ run-some-task-every-minute ] }

    # Environment Variables to the cloud container
    #
    # You usually will want to provide a file containing the environment variables
    # that your deployment should have. Such file may contain special Kool variables
    # that will be translated to their actual value by the Deploy API.
    #
    # There are two ways to provide such file.
    #
    # 'environment' refers to a file with all of the environment variables available.
    # The Deploy API will get the contents of this file and make all of them available
    # as true environment variables for the container. You can check them in your
    # deployed container via `kool cloud exec env` (assuming you have `env` available).
    environment: kool.deploy.env
    #
    # 'env' is a different option that actually allows you to build a file inside your
    # running container. This is useful for apps that are used for example to have a
    # `.env` file, but you do not want to have that built into your app Docker image.
    env:
      source: kool.deploy.env
      target: .env
```
