# kool - laravel

Start working with Laravel.

```bash
$ kool init laravel
```

By default laravel preset comes with `mysql` and `redis` configured, you can review how is configured at `docker-compose.yml`.

Also comes with some scripts to bring you up to speed at `kool.yaml`, take a look at the defaults.

```bash
$ kool run install
```

Now you can see your site at `http://localhost`, you can add more commands to your `kool.yml` or run away:

```bash
$ kool run php artisan tinker
```

```bash
$ kool run composer require something
```

```bash
$ kool run npm add something
```

After that you can simply start / stop.

```bash
$ kool start
```

```bash
$ kool stop
```
