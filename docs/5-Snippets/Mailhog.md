# Mailhog

MailHog is an email testing tool for developers.

## .env

```yaml
MAILER_URL=smtp://mailhog:1025
```

`mailhog` is the name of the service container created below in the **docker-compose.yml** file.

By default, Mailhog starts the SMTP server on port `1025`.

## docker-compose.yml

```yaml
# do not use in production!
mailhog:
  # Official Docker Image: https://hub.docker.com/r/mailhog/mailhog/
  image: mailhog/mailhog:latest
  environment:
    - MH_STORAGE=maildir
  # volumes:
  #   - ./docker/mailhog/maildir:/maildir:rw,delegated
  ports:
    - "8025:8025"
  networks:
    - kool_local
```

By default, MailHog uses in-memory message storage and starts the HTTP server on port 8025.

### Full Example

```diff
version: "3.7"
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
  database:
    image: mysql:8.0
    command: --default-authentication-plugin=mysql_native_password
    ports:
      - ${KOOL_DATABASE_PORT:-3306}:3306
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_PASSWORD-rootpass}
      MYSQL_DATABASE: ${DB_DATABASE-database}
      MYSQL_USER: ${DB_USERNAME-user}
      MYSQL_PASSWORD: ${DB_PASSWORD-pass}
      MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
    volumes:
      - database:/var/lib/mysql:delegated
    networks:
      - kool_local
    healthcheck:
      test:
        - CMD
        - mysqladmin
        - ping
  cache:
    image: redis:6-alpine
    volumes:
      - cache:/data:delegated
    networks:
      - kool_local
    healthcheck:
      test:
        - CMD
        - redis-cli
        - ping
+  mailhog:
+   # Official Docker Image: https://hub.docker.com/r/mailhog/mailhog/
+   image: mailhog/mailhog:latest
+   environment:
+     - MH_STORAGE=maildir
+   # volumes:
+   #   - mailhog:/maildir:delegated
+   ports:
+     - "8025:8025"
+   networks:
+     - kool_local
volumes:
  database: null
  cache: null
+ mailhog: null
networks:
  kool_local: null
  kool_global:
    external: true
    name: ${KOOL_GLOBAL_NETWORK:-kool_global}
```
