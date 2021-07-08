### `kool run mysql` or `kool run setup`

**Problem:**

> Access denied for root on mysql

**Answer:**

> Stop your containers with the purge option in order to delete all volumes

```shell
    kool stop --purge
```

> Make sure your DB username is other than root and define a password on the .env file

```diff
-DB_USERNAME=root
+DB_USERNAME=<some_user>

-DB_PASSWORD=
+DB_PASSWORD=<somepass>
```

> Start your container

```shell
	kool start
```

> Run your setup

```shell
	kool run setup
```
