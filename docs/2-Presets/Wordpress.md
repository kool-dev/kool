### Start from Scratch

For wordpress preset start from scratch or a existing wordpress is the same, go to the next section.

### Start with existing project

Go to the project folder and run:

```bash
$ kool init wordpress
```

**kool init** basically creates a few configuration files in order to enable you to configure / extend it.

By default wordpress preset comes with **mysql** and **redis** configured, you can review how is configured at **docker-compose.yml**.

At **docker-compose.yml** you can also see the default environment variables for database configuration that will be used for installing wordpress later on, the default are:

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
