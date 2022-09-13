### `kool run npm run hot` or `kool run npm run dev`

**Problem:**

> Browser cannot request "HMR" (hot module replacement) server successfully in Laravel Mix or Laravel Vite.

> Console errors: **net::ERR_CONNECTION_REFUSED** or **net::ERR_EMPTY_RESPONSE**.

**Answer:**

> Publish the HMR's port to the host and change HMR settings to listen on all IPv4 addresses (`0.0.0.0`).

For the sake of clarity, let's elect port `8080` to publish.

In your `kool.yml`, apply the following changes:

```diff
-npm: kool docker kooldev/node:16 npm
+npm: kool docker -p 8080:8080 kooldev/node:16 npm
```

- Alternatively, if you don't want to publish the port for your general `kool run npm` commands, you may add a new entry.

**> Laravel Mix**

In your `webpack.mix.js`, include the following changes to the `mix.options` call.

```diff
mix.options({
+    hmrOptions: {
+        host: '0.0.0.0',
+        port: 8080,
+    },
});
```

**> Laravel Vite**

In your `vite.config.js` file, include the following changes to the `defineConfig` call.

```diff
export default defineConfig({
+    server: {
+        host: '0.0.0.0',
+        port: 8080,
+    },
    plugins: [
        laravel({
            input: ['resources/css/app.css', 'resources/js/app.js'],
            refresh: true,
        }),
    ],
});
```
