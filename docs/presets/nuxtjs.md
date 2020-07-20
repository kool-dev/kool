# kool - nuxtjs

Start working with NuxtJS.

```bash
$ kool init nuxtjs
```

By default adonis preset only comes with node service, to add more services take a look at `docker-compose.yml` file.

Also comes with some scripts to bring you up to speed at `kool.yaml`, take a look at the defaults.

By default we already add a script `setup` to help you setup nuxtjs with kool for the first time.

```bash
$ kool run setup
```

Now you can see your site at `http://localhost:3000`, you can add more commands to your `kool.yml` or run away:

```bash
$ kool run npm install [package]
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
