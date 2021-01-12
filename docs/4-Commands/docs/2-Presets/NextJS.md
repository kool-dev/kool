### Create NextJS Project from Scratch

To make things easier we will use **kool** to install it for you.

```bash
kool create nextjs my-project

cd my-project
```
- **kool create** already executes **kool preset** internally so you can skip the command in the next step

Or

```bash
kool docker kooldev/node:14 yarn create next-app my-project

cd my-project
```

### Start using kool

Go to the project folder and run:

```bash
$ kool preset nextjs
```

**kool preset** basically creates a few configuration files in order to enable you to configure / extend it. You don't need to execute it whether you chose kool create command to start the new project.

Also comes with some scripts to bring you up to speed at **kool.yaml**, take a look at the defaults.

By default we always add a script called **setup** to help you setup a project for first time.

```bash
$ kool run setup
```

Now you can see your site at **http://localhost:3000**.

Check your **kool.yml** to see what scripts you can run and add more.
