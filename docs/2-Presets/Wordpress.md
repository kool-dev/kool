### Start from Scratch

If you want to create a project from scratch, just go to the directory you wanna have the new Wordpress install and run:

```bash
$ kool preset wordpress
$ kool start
```

Upon first start, if the current working directory is not an existing Wordpress source code, it will create it for you.

### Start with existing project

Go to the project folder and run:

```bash
$ kool preset wordpress
```

**kool preset** creates configuration files in order to enable you to configure and extend its behaviour.

The **wordpress** preset uses **mysql** and **redis** out of the box. You can review and change that at **docker-compose.yml**.

Run the following to start:

```bash
$ kool start
```

Now you can see your site at **http://localhost**.

To install wordpress, it will ask database credentials, the default ones are:

| Field         | Value    |
|---------------|----------|
| Database Name | database |
| Username      | user     |
| Password      | pass     |
| Database Host | database |
| Table Prefix  | wp_      |

Check your **kool.yml** to see what scripts you can run and add more.

### Publishing

For publishing we recommend using [Updraft Plus](https://wordpress.org/plugins/updraftplus):

* Install Plugin on your local
* Create Backup
* Download Backup
* Install Plugin on the destination with a fresh Wordpress Instalation
* Upload Backup
* Restore it
