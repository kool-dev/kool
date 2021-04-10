> We use PHP for our Hello World example, but **kool** is a **multi-stack** tool that can be used with any language and framework. For instance, [check out our presets](/docs/presets/introduction) for easily creating projects using JavaScript frameworks like Next.js, Nuxt.js, NestJS, and AdonisJs.

## Start a New Project with Docker in 4 Easy Steps

### Requirements

If you haven't done so already, you first need to [install Docker and the kool CLI](/docs/getting-started/installation).

Also, make sure you're running the latest version of **kool**. Run the following command to compare your local version of **kool** with the latest release, and, if a newer version is available, automatically download and install it.

```bash
$ kool self-update
```

### 1. Create a New Project

Create a working directory for your new project, and move into it.

```bash
$ mkdir my-app
$ cd my-app
```

### 2. Add docker-compose.yml

Create a **docker-compose.yml** file in the root of your working directory, and then copy/paste one of the YAML templates provided below into **docker-compose.yml**.

```bash
$ touch docker-compose.yml
```

Here's a simple, generic setup for **docker-compose.yml**:

```yaml
version: "3.8"  # optional since v1.27.0
services:
  app:
    image: kooldev/php:8.0-nginx
    ports:
      - "80:80"
    volumes:
      - .:/app:delegated
```

Here's a more extensible, **Kool-optimized** setup for **docker-compose.yml** into which you can easily add additional services (i.e. database, cache, etc):

```yaml
version: "3.8"  # optional since v1.27.0
services:
  app:
    image: kooldev/php:8.0-nginx
    ports:
      - ${KOOL_APP_PORT:-80}:80
    environment:
      ASUSER: ${KOOL_ASUSER:-0}
      UID: ${UID:-0}
    volumes:
      - .:/app:delegated
    networks:
      - kool_local
      - kool_global
networks:
  kool_local: null
  kool_global:
    external: true
    name: ${KOOL_GLOBAL_NETWORK:-kool_global}
```

> As you can see, we're using a [Kool-optimized Docker image](https://github.com/kool-dev/docker-php) for **PHP 8** (using `php:8.0-fpm-alpine`as its base), which also includes an NGINX web server. We're mapping **localhost** to container port `80`.

### 3. Hello World!

Create a `/public` sub-directory, and an **index.php** file inside it, and then copy/paste the code provided below into **index.php**.

```bash
$ mkdir public # in order to match NGINX's default root /app/public
$ touch public/index.php
```

```php
<?php

echo "Hello World! <br>";
echo "Yes, using Kool to create a new project running on Docker is super easy!";
```

> TIP: If you don't want to create a `/public` sub-directory for your **index.php** file (which matches NGINX's default setup), you can change the NGINX root by setting an `NGINX_ROOT` environment variable to `/app` in your **docker-compose.yml** file.

```yaml
services:
  app:
    environment:
      NGINX_ROOT: "/app"
```

### 4. Run Your App

Use the `kool start` command to start up the service containers defined in your **docker-compose.yml** file.

```bash
$ kool start
```

Verify your Docker containers are running using the `kool status` command:

```bash
$ kool status

+---------+---------+------------------------------+--------------+
| SERVICE | RUNNING | PORTS                        | STATE        |
+---------+---------+------------------------------+--------------+
| app     | Running | 0.0.0.0:80->80/tcp, 9000/tcp | Up 4 seconds |
+---------+---------+------------------------------+--------------+
```

You should now be able to access your new Hello World! app at [http://localhost](http://localhost). Hooray!

---

#### How to Open Container Sessions

Similar to SSH, if you want to open a Bash session in your `app` container, run `kool exec app bash`, where `app` is the name of the container service in **docker-compose.yml**. If you prefer, you can use `sh` instead of `bash` (`kool exec app sh`).

```bash
$ kool exec app bash
bash-5.1#

$ kool exec app sh
/app #
```

### Staying kool

When it's time to stop working on the project:

```bash
$ kool stop
```

When you're ready to start work again:

```bash
$ kool start
```

### BONUS: Add kool.yml

> Say hello to **kool.yml**, say goodbye to custom shell scripts!

Think of **kool.yml** as a super easy-to-use task _helper_. Instead of writing custom shell scripts, add your own scripts to **kool.yml** (under the `scripts` key), and run them with `kool run <script>` (e.g. `kool run setup`). You can add your own single line commands (see `composer` below), or add a list of commands that will be executed in sequence (see `setup` below).

Create a **kool.yml** file in the root of your working directory, and then copy/paste the code provided below into **kool.yml**.

```bash
$ touch kool.yml
```

Here's a **kool.yml** example to show you the types of commands you can add and use in your project to facilitate development:

```yaml
scripts:
  composer: kool exec app composer2
  node: kool docker kooldev/node:14 node
  yarn: kool docker kooldev/node:14 yarn
  node-setup:
    - kool run yarn install
    - kool run yarn dev
  setup:
    - cp .env.example .env
    - kool start
    - kool run composer install
    - kool run node-setup
```

> As you can see in this **kool.yml** example, the `setup` script will do the following in sequence: copy your updated **.env.example** file to **.env**; start your Docker environment; use Composer to install vendor dependencies; generate your `APP_KEY` (in `.env`); and then build your Node packages and assets.
