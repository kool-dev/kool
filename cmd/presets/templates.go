package presets

// auto generated file

// GetTemplates get all templates
func GetTemplates() map[string]map[string]string {
	var templates = make(map[string]map[string]string)
	templates["app"] = map[string]string{
		"hugo.yml": `version: "3.8"
services:
  app:
    image: klakegg/hugo
    command: ["server", "-p", "80", "-D"]
    working_dir: /app
    ports:
      - "${KOOL_APP_PORT:-80}:80"
    volumes:
      - .:/app:delegated
    networks:
      - kool_local
      - kool_global
  static:
    image: kooldev/nginx:static
    volumes:
      - .:/app:delegated
    networks:
      - kool_local
      - kool_global
networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"
`,
		"node14-adonis.yml": `services:
  app:
    image: kooldev/node:14-adonis
    command: ["adonis", "serve", "--dev"]
    ports:
      - "${KOOL_APP_PORT:-3333}:3333"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
      - .:/app:delegated
    networks:
      - kool_local
      - kool_global
`,
		"node14-nestjs.yml": `services:
  app:
    image: kooldev/node:14-nest
    command: ["npm", "run", "start:dev"]
    ports:
      - "${KOOL_APP_PORT:-3000}:3000"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
      - .:/app:delegated
    networks:
      - kool_local
      - kool_global
`,
		"php74.yml": `services:
  app:
    image: kooldev/php:7.4-nginx
    ports:
      - "${KOOL_APP_PORT:-80}:80"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
      - .:/app:delegated
    networks:
      - kool_local
      - kool_global
`,
		"php8.yml": `services:
  app:
    image: kooldev/php:8.0-nginx
    ports:
      - "${KOOL_APP_PORT:-80}:80"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
      - .:/app:delegated
    networks:
      - kool_local
      - kool_global
`,
		"wordpress74.yml": `services:
  app:
    image: kooldev/wordpress:7.4-nginx
    ports:
      - "${KOOL_APP_PORT:-80}:80"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
      - .:/app:delegated
    networks:
      - kool_local
      - kool_global
`,
	}
	templates["cache"] = map[string]string{
		"memcached16.yml": `services:
  cache:
    image: memcached:1.6-alpine
    volumes:
      - cache:/data:delegated
    networks:
      - kool_local

volumes:
  cache:
`,
		"redis6.yml": `services:
  cache:
    image: redis:6-alpine
    volumes:
      - cache:/data:delegated
    networks:
      - kool_local
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]

volumes:
  cache:
`,
	}
	templates["database"] = map[string]string{
		"mysql57-adonis.yml": `services:
  database:
    image: mysql:5.7
    ports:
      - "${KOOL_DATABASE_PORT:-3306}:3306"
    environment:
      MYSQL_ROOT_PASSWORD: "${DB_PASSWORD-rootpass}"
      MYSQL_DATABASE: "${DB_DATABASE-database}"
      MYSQL_USER: "${DB_USER-user}"
      MYSQL_PASSWORD: "${DB_PASSWORD-pass}"
      MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
    volumes:
      - database:/var/lib/mysql:delegated
    networks:
      - kool_local
    healthcheck:
      test: ["CMD", "mysqladmin", "ping"]

volumes:
  database:

scripts:
  mysql: kool exec -e MYSQL_PWD=$DB_PASSWORD database mysql -u $DB_USER $DB_DATABASE
`,
		"mysql57.yml": `services:
  database:
    image: mysql:5.7
    ports:
      - "${KOOL_DATABASE_PORT:-3306}:3306"
    environment:
      MYSQL_ROOT_PASSWORD: "${DB_PASSWORD-rootpass}"
      MYSQL_DATABASE: "${DB_DATABASE-database}"
      MYSQL_USER: "${DB_USERNAME-user}"
      MYSQL_PASSWORD: "${DB_PASSWORD-pass}"
      MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
    volumes:
      - database:/var/lib/mysql:delegated
    networks:
      - kool_local
    healthcheck:
      test: ["CMD", "mysqladmin", "ping"]

volumes:
  database:

scripts:
  mysql: kool exec -e MYSQL_PWD=$DB_PASSWORD database mysql -u $DB_USERNAME $DB_DATABASE
`,
		"mysql8-adonis.yml": `services:
  database:
    image: mysql:8.0
    command: --default-authentication-plugin=mysql_native_password
    ports:
      - "${KOOL_DATABASE_PORT:-3306}:3306"
    environment:
      MYSQL_ROOT_PASSWORD: "${DB_PASSWORD-rootpass}"
      MYSQL_DATABASE: "${DB_DATABASE-database}"
      MYSQL_USER: "${DB_USER-user}"
      MYSQL_PASSWORD: "${DB_PASSWORD-pass}"
      MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
    volumes:
      - database:/var/lib/mysql:delegated
    networks:
      - kool_local
    healthcheck:
      test: ["CMD", "mysqladmin", "ping"]

volumes:
  database:

scripts:
  mysql: kool exec -e MYSQL_PWD=$DB_PASSWORD database mysql -u $DB_USER $DB_DATABASE
`,
		"mysql8.yml": `services:
  database:
    image: mysql:8.0
    command: --default-authentication-plugin=mysql_native_password
    ports:
      - "${KOOL_DATABASE_PORT:-3306}:3306"
    environment:
      MYSQL_ROOT_PASSWORD: "${DB_PASSWORD-rootpass}"
      MYSQL_DATABASE: "${DB_DATABASE-database}"
      MYSQL_USER: "${DB_USERNAME-user}"
      MYSQL_PASSWORD: "${DB_PASSWORD-pass}"
      MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
    volumes:
      - database:/var/lib/mysql:delegated
    networks:
      - kool_local
    healthcheck:
      test: ["CMD", "mysqladmin", "ping"]

volumes:
  database:

scripts:
  mysql: kool exec -e MYSQL_PWD=$DB_PASSWORD database mysql -u $DB_USERNAME $DB_DATABASE
`,
		"postgresql13-adonis.yml": `services:
  database:
    image: postgres:13-alpine
    ports:
      - "${KOOL_DATABASE_PORT:-5432}:5432"
    environment:
      POSTGRES_DB: "${DB_DATABASE-database}"
      POSTGRES_USER: "${DB_USER-user}"
      POSTGRES_PASSWORD: "${DB_PASSWORD-pass}"
      POSTGRES_HOST_AUTH_METHOD: "trust"
    volumes:
      - database:/var/lib/postgresql/data:delegated
    networks:
      - kool_local
    healthcheck:
      test: ["CMD", "pg_isready", "-q", "-d", "$DB_DATABASE", "-U", "$DB_USER"]

volumes:
  database:

scripts:
  psql: kool exec -e PGPASSWORD=$DB_PASSWORD database psql -U $DB_USER $DB_DATABASE
`,
		"postgresql13.yml": `services:
  database:
    image: postgres:13-alpine
    ports:
      - "${KOOL_DATABASE_PORT:-5432}:5432"
    environment:
      POSTGRES_DB: "${DB_DATABASE-database}"
      POSTGRES_USER: "${DB_USERNAME-user}"
      POSTGRES_PASSWORD: "${DB_PASSWORD-pass}"
      POSTGRES_HOST_AUTH_METHOD: "trust"
    volumes:
      - database:/var/lib/postgresql/data:delegated
    networks:
      - kool_local
    healthcheck:
      test: ["CMD", "pg_isready", "-q", "-d", "$DB_DATABASE", "-U", "$DB_USERNAME"]

volumes:
  database:

scripts:
  psql: kool exec -e PGPASSWORD=$DB_PASSWORD database psql -U $DB_USERNAME $DB_DATABASE
`,
	}
	templates["scripts"] = map[string]string{
		"composer.yml": `scripts:
  composer: kool exec app composer
`,
		"composer2.yml": `scripts:
  composer: kool exec app composer2
`,
		"hugo.yml": `scripts:
  hugo: kool docker -p 1313:1313 klakegg/hugo
  dev: kool run hugo server -D

  setup:
    - kool start
    - kool run dev
`,
		"npm-adonis.yml": `scripts:
  adonis: kool exec app adonis
  npm: kool exec app npm
  npx: kool exec app npx

  setup:
    - kool docker kooldev/node:14 npm install
    - kool start
`,
		"npm-laravel.yml": `scripts:
  artisan: kool exec app php artisan
  npm: kool exec app npm
  npx: kool exec app npx

  node-setup:
    - kool run npm install
    - kool run npm run dev

  setup:
    - cp .env.example .env
    - kool start
    - kool run composer install
    - kool run artisan key:generate
    - kool run node-setup

  reset:
    - kool run composer install
    - kool run artisan migrate:fresh --seed
    - kool run node-setup
`,
		"npm-nestjs.yml": `scripts:
  nest: kool exec app nest
  npm: kool exec app npm
  npx: kool exec app npx

  setup:
    - kool docker kooldev/node:14 npm install
    - kool start
`,
		"npm-nextjs.yml": `scripts:
  npm: kool exec app npm
  npx: kool exec app npx

  setup:
    - kool docker kooldev/node:14 npm install
    - kool start
`,
		"npm-nodejs.yml": `scripts:
  node: kool exec app node
  npm: kool exec app npm
  npx: kool exec app npx

  setup:
    - kool start
    # - add more setup commands
`,
		"npm-nuxtjs.yml": `scripts:
  npm: kool exec app npm
  npx: kool exec app npx

  setup:
    - kool docker kooldev/node:14 npm install
    - kool start
`,
		"npm.yml": `scripts:
  npm: kool exec app npm
  npx: kool exec app npx
`,
		"php.yml": `scripts:
  php: kool exec app php

  setup:
    - kool start
    # - add more setup commands
`,
		"symfony.yml": `scripts:
  console: kool exec app php ./bin/console
  phpunit: kool exec app php ./bin/phpunit

  setup:
    - kool start
    - kool run composer install
    #- kool run console doctrine:migrations:migrate -n
`,
		"wordpress.yml": `scripts:
  php: kool exec app php
  wp: kool exec app wp
`,
		"yarn-adonis.yml": `scripts:
  adonis: kool exec app adonis
  yarn: kool exec app yarn

  setup:
    - kool docker kooldev/node:14 yarn install
    - kool start
`,
		"yarn-laravel.yml": `scripts:
  artisan: kool exec app php artisan
  yarn: kool exec app yarn

  node-setup:
    - kool run yarn install
    - kool run yarn dev

  setup:
    - cp .env.example .env
    - kool start
    - kool run composer install
    - kool run artisan key:generate
    - kool run node-setup

  reset:
    - kool run composer install
    - kool run artisan migrate:fresh --seed
    - kool run node-setup
`,
		"yarn-nestjs.yml": `scripts:
  nest: kool exec app nest
  yarn: kool exec app yarn

  setup:
    - kool docker kooldev/node:14 yarn install
    - kool start
`,
		"yarn-nextjs.yml": `scripts:
  yarn: kool exec app yarn

  setup:
    - kool docker kooldev/node:14 yarn install
    - kool start
`,
		"yarn-nodejs.yml": `scripts:
  node: kool exec app node
  yarn: kool exec app yarn

  setup:
    - kool start
    # - add more setup commands
`,
		"yarn-nuxtjs.yml": `scripts:
  yarn: kool exec app yarn

  setup:
    - kool docker kooldev/node:14 yarn install
    - kool start
`,
		"yarn.yml": `scripts:
  yarn: kool exec app yarn
`,
	}
	return templates
}
