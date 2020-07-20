# kool - adonis

Start working with Adonis.

```bash
$ kool init adonis
```

By default adonis preset comes with `mysql` and `redis` pre-configured, to enable you can uncomment it at `docker-compose.yml` file.

Also comes with some scripts to bring you up to speed at `kool.yaml`, take a look at the defaults.

By default we already add a script `setup` to help you setup adonis with kool for the first time, but adonis requires some environment changes in order to work, for that we added a file `kool.env` to show what changes you have to do yo our `.env` file.

Make these changes to you `.env.example` file or run the commands from `setup` script manually.

```bash
# CAUTION, this script will reset your `.env` file with `.env.example`
$ kool run setup
```

Now you can see your site at `http://localhost`, you can add more commands to your `kool.yml` or run away:

```bash
$ kool run npm install [package]
```

```bash
$ kool run adonis migration:run
```

```bash
$ kool run node path/to/file.js
```

After that you can simply start / stop.

```bash
$ kool start
```

```bash
$ kool stop
```



