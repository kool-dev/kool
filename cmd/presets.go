package cmd

// auto generated file

var presets = make(map[string]map[string]string)
func init() {
	presets["adonis"] = map[string]string{
		"Dockerfile.build": `FROM fireworkweb/node:14-adonis

RUN npm install

EXPOSE 3333

CMD [ "npm", "start" ]`,
		"docker-compose.yml": `version: "3.8"
services:
  app:
    image: fireworkweb/node:14-adonis
    command: ["adonis", "serve", "--dev"]
    ports:
     - "${PORT:-3333}:3333"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:cached
     - ${HOME:-/dev/null}:/home/fwd/.ssh:cached
##########################################################################
# optionally you can enable mysql, redis or other services to your stack #
##########################################################################
#   database:
#     image: mysql:8.0 # possibly change to: mysql:5.7
#     ports:
#      - "3306:3306"
#     environment:
#       MYSQL_ROOT_PASSWORD: "${DB_PASSWORD:-rootpass}"
#       MYSQL_DATABASE: "${DB_DATABASE:-database}"
#       MYSQL_USER: "${DB_USERNAME:-user}"
#       MYSQL_PASSWORD: "${DB_PASSWORD:-}"
#       MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
#     volumes:
#      - db:/var/lib/mysql:cached
#   cache:
#     image: redis:alpine
#     volumes:
#      - cache:/data:cached
# volumes:
#   db:
#   cache:`,
		"kool.yml": `scripts:
  # node
  node: kool exec app node
  npm: kool exec app npm # can change to: yarn,pnpm
  adonis: kool exec app adonis

  # examples
  install:
    - cp .env.example .env
    - kool run --docker fireworkweb/node:14 npm install # can change to: yarn,pnpm
    - kool start`,
	}
	presets["laravel"] = map[string]string{
		"Dockerfile.build": `FROM fireworkweb/php:7.4 as composer

COPY . /app
RUN composer install --no-interaction --prefer-dist --optimize-autoloader --quiet

FROM fireworkweb/node:14 as node

COPY --from=composer /app /app
RUN yarn install && yarn prod

FROM fireworkweb/php:7.4-nginx

COPY --chown=fwd:fwd --from=node /app /app`,
		"docker-compose.yml": `version: "3.8"
services:
  app:
    image: fireworkweb/php:7.4-nginx
    ports:
     - "80:80"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:cached
     - ${HOME:-/dev/null}:/home/fwd/.ssh:cached
  database:
    image: mysql:8.0 # can change to: mysql:5.7
    command: --default-authentication-plugin=mysql_native_password # remove this line if you change to: mysql:5.7
    ports:
     - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: "${DB_PASSWORD:-rootpass}"
      MYSQL_DATABASE: "${DB_DATABASE:-database}"
      MYSQL_USER: "${DB_USERNAME:-user}"
      MYSQL_PASSWORD: "${DB_PASSWORD:-}"
      MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
    volumes:
     - db:/var/lib/mysql:cached
  cache:
    image: redis:alpine
    volumes:
     - cache:/data:cached
volumes:
  db:
  cache:`,
		"kool.yml": `scripts:
  # php
  php: kool exec app php
  composer: kool exec app composer

  # node
  node: kool run --docker fireworkweb/node:14 node
  npm: kool run --docker fireworkweb/node:14 npm

  # examples
  install:
    - kool start
    - cp .env.example .env
    - kool run composer install
    - kool run php artisan key:generate
    - kool run npm install
    - kool run npm dev

  reset:
    - kool run composer install
    - kool run php artisan migrate:fresh --seed
    - kool run npm install # can change to: yarn,pnpm
    - kool run npm dev # can change to: yarn,pnpm`,
	}
	presets["nextjs"] = map[string]string{
		"Dockerfile.build": `FROM fireworkweb/node:14

RUN npm install && npm run build

EXPOSE 3000

CMD [ "npm", "start" ]`,
		"docker-compose.yml": `version: "3.8"
services:
  app:
    image: fireworkweb/node:14
    command: ["npm", "run", "dev"]
    ports:
     - "3000:3000"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:cached
     - ${HOME:-/dev/null}:/home/fwd/.ssh:cached`,
		"kool.yml": `start:
  services: app

scripts:
  # node
  node: kool exec app node
  npm: kool exec app npm # can change to: yarn,pnpm

  # examples
  install:
    - kool run --docker fireworkweb/node:14 npm install # can change to: yarn,pnpm
    - kool start`,
	}
}
