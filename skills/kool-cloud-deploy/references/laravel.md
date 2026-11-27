# Laravel on Kool Cloud
This is the shape of a production Laravel app on Kool: the `kooldev/php` prod image, migrations in a hook, the queue worker as a daemon, the scheduler as a cron job, and env from a git-ignored file.
```yaml
# kool.cloud.yml
services:
  app:
    build: .
    expose: 80
    public: true
    size: small                      # micro (256Mi) is too small for php-fpm + nginx
    root: /app
    environment: kool.deploy.env     # git-ignored dotenv file, uploaded at each deploy
    readinessProbe:                  # the /up route ships with Laravel 11+
      httpGet: { path: /up, port: 80 }
    after:                           # once per deploy, after the new pods are up
      - [php, artisan, migrate, --force]
    daemons:                         # Horizon: [php, artisan, horizon] at size medium
      - command: [php, artisan, queue:work, --tries=3, --max-time=3600]
        size: small
    cron:
      - schedule: "* * * * *"
        command: [php, artisan, schedule:run]
    persists:                        # only if files are stored locally (drop it for S3)
      size: 1Gi
      paths:
        - { path: storage, sync: true, chown: "kool:kool" } # /app/storage, seeded from the image
  database:
    image: mysql:8.0
    cloud: true                      # managed shared MySQL 8; fills {{KOOL_DATABASE_*}}
```
```dockerfile
# Dockerfile. Prod image: opcache on (validate_timestamps=0), nginx and php-fpm workers run
# as `kool`, and the entrypoint runs hooks/daemons/cron as `kool`. Don't add USER.
FROM kooldev/php:8.4-debian-nginx-prod AS app
COPY --chown=kool:kool composer.json composer.lock /app/
RUN composer install --no-dev --prefer-dist --no-progress --no-scripts --no-autoloader
COPY --chown=kool:kool . /app
RUN composer install --no-dev --prefer-dist --no-progress --optimize-autoloader

# Vite assets are arch-independent: build them natively (no package.json? delete the stage)
FROM --platform=$BUILDPLATFORM node:24-slim AS assets
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
COPY --from=app /app/vendor ./vendor
RUN npm run build

FROM app
COPY --chown=kool:kool --from=assets /app/public/build /app/public/build
RUN ln -sfn /app/storage/app/public /app/public/storage
```
Match the image tag to composer.json's PHP version: 8.2 to 8.5 exist.

`.dockerignore`, one pattern per line: `.git`, `.env*`, `kool.deploy.env`, `node_modules`, `vendor`, `public/build`, `public/hot`, `public/storage`, `storage/logs/*.log`, `storage/framework/views/*.php`, `database/*.sqlite`.

## kool.deploy.env (add it to .gitignore; never print it)
```dotenv
APP_ENV=production
APP_DEBUG=false
APP_URL=https://myapp.kool.cloud
LOG_CHANNEL=stderr
DB_CONNECTION=mysql
DB_HOST="{{KOOL_DATABASE_HOST}}"
DB_PORT="{{KOOL_DATABASE_PORT}}"
DB_DATABASE="{{KOOL_DATABASE_DATABASE}}"
DB_USERNAME="{{KOOL_DATABASE_USERNAME}}"
DB_PASSWORD="{{KOOL_DATABASE_PASSWORD}}"
```
- **APP_KEY:** generate it once without displaying it, with `printf 'APP_KEY=base64:%s\n' "$(openssl rand -base64 32)" >> kool.deploy.env`, and never regenerate it: a new key breaks sessions and encrypted data. Or write `APP_KEY="{{APP_KEY}}"` and have the human set APP_KEY in the panel.
- **Variables:** only this file's variables reach the containers; the panel's variables only fill `{{ }}`. With your own database, delete the `database` service and write literal `DB_*` values.
- **Logs:** `LOG_CHANNEL=stderr` sends the logs to `kool cloud logs app`.

## Notes
- **Migrations:** `after` runs once per deploy, so the new code serves for a few seconds before they finish. A `before` hook works too, the first deploy included, but it runs at every pod start and in each replica. With `replicas` above 1, a `before` hook needs `--isolated=1` and a Redis cache (the database cache's lock table doesn't exist on the first deploy).
- **Scheduler and queue:** each scheduler run is a fresh pod, so `runInBackground()` tasks die when the pod exits; send long work to the queue. Queue workers are daemons in the app pod, so they restart with every deploy and scale with `replicas`.
- **`storage` persistence:** ownership on the disk comes from the image (`COPY --chown`); `chown` applies only when the path is missing from the image. `persists` makes each deploy stop-then-start, and on a first deploy `before` hooks see an empty `storage`.
- **HTTPS:** TLS ends at Kool's ingress, so trust the proxy, or `asset()` URLs and redirects come out as `http://`. On Laravel 11+, put `$middleware->trustProxies(at: '*');` in `bootstrap/app.php`.
- **Caching:** don't run `config:cache` at build time (the env only arrives at runtime) or in a hook (hooks have their own filesystem).
- **docker-compose.yml (kool preset, Sail):** its services are merged in. Keep `environment:` on `app` so the compose `${...}` values don't carry over. Set `deploy: false` on compose services you don't deploy. To run Redis in the cluster, add `cache: {image: redis:7-alpine, expose: 6379}` and set `REDIS_HOST=cache`.
