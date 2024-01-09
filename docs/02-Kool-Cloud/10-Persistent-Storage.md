Containers were built to be ephemeral—and that is how we like them and how Kubernetes and all other container orchestrators usually work the best with them as well.

But at times, we know that traditional web applications may not be ready to switch to network-based object storage instead of local disk storage.

Kool.dev Cloud does offer you the ability to create persisted paths within your deployed containers.

```yaml
services:
  app:
    # ...

    # Tells your app's root folder so all other paths can be relative (optional).
    root: /app

    # Containers are ephemeral, that means their filesystem do not persist across restarts.
    # If you want to persist stuff into the disk across deployments, you can do so by defining persistent paths here.
    persists:
      # Total size of the volume you want to attach to the running container.
      # This can be increased later, but it may take a while to apply the change.
      size: 10Gi
      # Paths to persist - within that single volume, you can have one or more paths
      # that are going to be mounted every time your containers are running. Note that
      # such mounts will be there for before/after hooks as well as daemon containers.
      paths:
        # The path within the container. Must be either aboslute or relative to the 'root' config.
        - path: /app/some/path/persisted
          # Tells the Deploy API to sync the folder from your built image to the persisted storage.
          # This is very helpful to start off with a known folder structure.
          sync: true
          # Tells what user and group should own the persisted folder (only used when sync: true)
          chown: user:group
```