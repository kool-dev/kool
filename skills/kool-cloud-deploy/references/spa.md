# Vite / React SPA on Kool Cloud

Build the bundle in a node stage, then serve `dist/` from nginx with an SPA fallback.

```yaml
# kool.cloud.yml
services:
  app:
    build:
      context: .
      args:                                    # VITE_* values are baked in at build time
        VITE_API_URL: "https://api.example.com"
    expose: 80
    public: true
```
With no `VITE_*` variables, `build: .` is enough.

```dockerfile
# Dockerfile
# The bundle is arch-independent: build it natively. On Apple Silicon an amd64 node
# stage would run under QEMU, where Rolldown or esbuild can segfault.
FROM --platform=$BUILDPLATFORM node:24-slim AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

FROM nginx:1.27-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
```
For pnpm or yarn, replace the lockfile and the install command (`pnpm install --frozen-lockfile`, `yarn install --immutable`).

```nginx
# nginx.conf
server {
    listen 80;
    root /usr/share/nginx/html;
    index index.html;
    location /assets/ {
        expires 1y;
        add_header Cache-Control "public, immutable";
        try_files $uri =404;
    }
    location / {
        try_files $uri $uri/ /index.html;   # SPA fallback: BrowserRouter deep links
    }
}
```

```gitignore
# .dockerignore
.git
.env*
node_modules
dist
```

## Notes
- **The fallback is required.** Without `try_files ... /index.html`, deep links such as `/pricing` return 404 from nginx. HashRouter doesn't need it.
- **`VITE_*` values are compile-time.** Setting `environment:` in kool.cloud.yml has no effect on a built SPA. Pass build args as above: `"$VAR"` is read from your shell, and `"{{VAR}}"` from the environment's variables in the Kool panel. Never put secrets in `VITE_*`, because they ship to every browser.
- **`.env*` stays out of the image.** Vite would read a local `.env.production` during the build and bake its values into the bundle.
- **Other build output.** Some tools emit `build/` instead of `dist/` (for example Create React App), and some apps use a custom `base` path. Adjust the `COPY --from=build` path to match.
- **Verify:** `curl -sS https://<domain>/ | grep -i '<title>'`, then `curl -sS -o /dev/null -w '%{http_code}\n' https://<domain>/some/deep/route`, which must print 200.
