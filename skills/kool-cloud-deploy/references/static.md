# Static site (plain HTML/CSS/JS) on Kool Cloud

nginx serves the files as they are, with no build step. Files in a `dist/` or `public/` folder work the same way: copy that folder instead.

```yaml
# kool.cloud.yml
services:
  app:
    build: .
    expose: 80
    public: true
```

```dockerfile
# Dockerfile
FROM nginx:1.27-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY . /usr/share/nginx/html
RUN rm /usr/share/nginx/html/nginx.conf
EXPOSE 80
```
`nginx:alpine` is multi-arch, so the image builds for amd64 without emulating any build step. Keep the image's own `CMD`: it runs nginx in the foreground. The first `COPY` needs `nginx.conf` in the build context, so don't list it in `.dockerignore`; the `RUN rm` keeps it off the site.

```nginx
# nginx.conf
server {
    listen 80;
    root /usr/share/nginx/html;
    index index.html;
    location / {
        try_files $uri $uri/ $uri.html =404;
    }
}
```
`try_files` makes `/about` serve `about.html`. Change `=404` to `/404.html` if the site has a custom error page.

```gitignore
# .dockerignore
.git
.env*
Dockerfile
.dockerignore
kool.cloud.yml
docker-compose.yml
node_modules
README.md
```
`COPY .` would otherwise publish these files, including `.git` and any `.env`, at public URLs.

## Notes
- Don't write a docker-compose.yml: Kool doesn't need one. If one already exists, its services are merged into the deploy, so check that every service in it is meant to be deployed.
- The site's content is baked into the image, so redeploy after every change: `kool cloud deploy --domain <name>.kool.cloud`.
- Verify: `curl -sS https://<domain>/ | grep -i '<title>'`, and request one asset, such as `/app.js`, to check it returns 200.
- Forms or a waitlist need a backend that the browser can reach. Use an external form service, or deploy an API (see node.md) as its own app on its own domain.
