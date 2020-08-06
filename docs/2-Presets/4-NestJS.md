### Create NestJS Project from Scratch

To make things easier we will use **kool** to install it for you.

```bash
kool docker kooldev/node:14-nest nest new my-project

cd my-project
```

### Start using kool

Go to the project folder and run:

```bash
$ kool init nestjs
```

**kool init** basically creates a few configuration files in order to enable you to configure / extend it.

By default nestjs preset comes with **mysql** and **redis** pre-configured, to enable you can uncomment it at **docker-compose.yml** file.

The host you will use in your application config for any service using docker like mysql, redis or mongo will be the service name in the **docker-compose.yml**, so **mysql** will be referenced as **database**, feel free to change the name.

Also comes with some scripts to bring you up to speed at **kool.yaml**, take a look at the defaults.

By default we always add a script called **setup** to help you setup a project for first time.

```bash
$ kool run setup
```

Now you can see your site at **http://localhost:3000**.

Check your **kool.yml** to see what scripts you can run and add more.
