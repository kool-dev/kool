Most of time the web applications you are going to deploy will usually have its own companying database.

Kool.dev Cloud offers two ways for you do deploy your databases.

1. Deploy databases as regular containers.
1. Use a managed database on a shared structure.
1. Use a dedicated Database Custom Resource (i.e RDS dedicated instance).

All of this options come with basic daily backup routines for your safety.

### Deploy Databases as Regular Containers

Deploying a container that is going to run your database is pretty straight forward - just like your have it in your local `docker-compose.yml` for your local environment, you can deploy that very same container. The benefit of this is you have full control at your container configuration and can use any type of database.

#### Caveats of deploying databases on containers are

- **Persistent disk storage**: by default deployed containers are ephemeral and DO NOT have any disk persistent storage. This may cause problems if you deploy a database and upon restart, all your data is lost. **You must make sure to incoude in your container deploy configuration a `persistent` disk storage**, so upon restarts your data is kept safe and is no longer ephemeral.
- **Environment variables**: your database image may require specific environment variables to determine credentials and other settings. You need to make sure you set them correctly, different than your local ones.

### Managed Database in shared structure

This option is the easiest to get started - but currently only supports MySQL 8 database deployments.

If you have a MySQL database in your `docker-compose.yml`, you can just assign that service the `cloud: true` on your `kool.cloud.yml` and Kool.dev Cloud is going to setup a new database on a running shared RDS instance.

This managed options will provide you with variables placeholders for you to get a hold of the credentials automatically generated as well as the database name/host.

Here is the list of Kool.dev variables placeholders available and how you would use them in your environment variables definition to use the managed database:

```
DB_HOST="{{KOOL_DATABASE_HOST}}"
DB_PORT={{KOOL_DATABASE_PORT}}
DB_DATABASE="{{KOOL_DATABASE_DATABASE}}"
DB_USERNAME="{{KOOL_DATABASE_USERNAME}}"
DB_PASSWORD="{{KOOL_DATABASE_PASSWORD}}"
```

The placeholders always have the `{{PLACEHOLDER}}` syntax. When used anywhere in your `kool.cloud.yml` configuration they are going to be replaced by their managed values when deploying.

#### Caveats of using managed shared database

- Currently **only supports MySQL 8** deployments.
- Being a shared resource, top performance is not guarenteed (unless you have it running in your own Cloud vendor account in the Enterprise offer).
- Best suited for development and staging workloads.

### Dedicated Database Resource

You can have any sort of custom resource for your application, including dedicated databases (i.e RDS or ElastiCache).

As this is not yet fully automated you need to [contact our support to set it up](mailto:contact@kool.dev) for you in your account.

One of the benefits is having total control of your set up not only on disk/computing performance, but as well as tailored backup and replication options.
