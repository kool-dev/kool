### Create Laravel Project from Scratch

To make things easier we will use **kool** to install it for you.

```bash
kool docker kooldev/php:7.4 composer create-project --prefer-dist laravel/laravel my-project

cd my-project
```

### Start using kool

Go to the project folder and run:

```bash
$ kool init laravel
```

**kool init** basically creates a few configuration files in order to enable you to configure / extend it.

By default laravel preset comes with **mysql** and **redis** configured, you can review how is configured at **docker-compose.yml**.

Also comes with some scripts to bring you up to speed at **kool.yaml**, take a look at the defaults.

By default we always add a script called **setup** to help you setup a project for first time, but symfony also requires some environment changes in order to work with **docker**, bellow are the following changes:

```bash
DB_USERNAME=user
DB_PASSWORD=pass
DB_HOST=database
DB_PORT=3306
DB_DATABASE=database
DB_VERSION=8.0
DATABASE_URL=mysql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?serverVersion=${DB_VERSION}
```

The host you will use in your application config for any service using docker like mysql, redis or mongo will be the service name in the **docker-compose.yml**, so **mysql** will be referenced as **database**, feel free to change the name.

We recommend making these changes to you **.env.example** file to avoid steps on future installations.

```bash
$ kool run setup
```

Now you can see your site at **http://localhost**.

Check your **kool.yml** to see what scripts you can run and add more.
