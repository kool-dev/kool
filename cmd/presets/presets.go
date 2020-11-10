package presets

// auto generated file

// GetAll get all presets
func GetAll() map[string]map[string]string {
	var presets = make(map[string]map[string]string)
	presets["adonis"] = map[string]string{
		".dockerignore": `/node_modules
`,
		"preset_language": "javascript",
		"Dockerfile.build": `FROM kooldev/node:14-adonis AS build

COPY . /app

RUN npm install

FROM kooldev/node:14-adonis

COPY --from=build --chown=kool:kool /app /app

EXPOSE 3333

CMD [ "npm", "start" ]
`,
		"docker-compose.yml": `version: "3.7"
services:
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
    #  - $HOME/.ssh:/home/kool/.ssh:delegated
    networks:
     - kool_local
     - kool_global
#   database:
#     image: mysql:8.0 # possibly change to: mysql:5.7
#     command: --default-authentication-plugin=mysql_native_password
#     ports:
#      - "${KOOL_DATABASE_PORT:-3306}:3306"
#     environment:
#       MYSQL_ROOT_PASSWORD: "${DB_PASSWORD:-rootpass}"
#       MYSQL_DATABASE: "${DB_DATABASE:-database}"
#       MYSQL_USER: "${DB_USERNAME:-user}"
#       MYSQL_PASSWORD: "${DB_PASSWORD:-pass}"
#       MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
#     volumes:
#      - db:/var/lib/mysql:delegated
#     networks:
#      - kool_local
#   cache:
#     image: redis:6-alpine
#     volumes:
#      - cache:/data:delegated
#
# volumes:
#   db:
#   cache:

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"
`,
		"kool.yml": `scripts:
  node: kool exec app node
  npm: kool exec app npm # can change to: yarn,pnpm
  adonis: kool exec app adonis

  setup:
    - kool docker kooldev/node:14 npm install # can change to: yarn,pnpm
    - kool start
`,
	}
	presets["golang-cli"] = map[string]string{
		"preset_language": "golang",
		"kool.yml": `scripts:
  # Helper for local development - compiling and installing locally
  dev:
    - kool run compile
    - kool run install

  # Runs go CLI with proper version for kool development
  go: kool docker --volume=gopath:/go --env='GOOS=$GOOS' golang:1.15.0 go

  # Compiling cli itself. In case you are on MacOS make sure to have your .env
  # file properly setting GOOS=darwin so you will be able to use the binary.
  compile: kool run go build -o my-cli
  install:
    - mv my-cli /usr/local/bin/my-cli
  fmt: kool run go fmt
  lint: kool docker --volume=gopath:/go golangci/golangci-lint:v1.31.0 golangci-lint run -v
`,
	}
	presets["laravel"] = map[string]string{
		".dockerignore": `/node_modules
/vendor
`,
		"preset_language": "php",
		"preset_app_template": "php74.yml",
		"preset_ask_database": "MySQL 8.0,MySQL 5.7,ProstgreSQL,none",
		"preset_ask_cache": "Redis 6.0,Memcached 1.6,none",
		"Dockerfile.build": `FROM kooldev/php:7.4 AS composer

COPY . /app
RUN composer install --no-interaction --prefer-dist --optimize-autoloader --quiet

FROM kooldev/node:14 AS node

COPY --from=composer /app /app
RUN yarn install && yarn prod

FROM kooldev/php:7.4-nginx

COPY --from=node --chown=kool:kool /app /app
`,
		"docker-compose.yml": `version: "3.7"
services:
  app:
    image: kooldev/php:7.4-nginx
    ports:
     - "${KOOL_APP_PORT:-80}:80"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:delegated
    #  - $HOME/.ssh:/home/kool/.ssh:delegated
    networks:
     - kool_local
     - kool_global
  database:
    image: mysql:8.0 # can change to: mysql:5.7
    command: --default-authentication-plugin=mysql_native_password # remove this line if you change to: mysql:5.7
    ports:
     - "${KOOL_DATABASE_PORT:-3306}:3306"
    environment:
      MYSQL_ROOT_PASSWORD: "${DB_PASSWORD:-rootpass}"
      MYSQL_DATABASE: "${DB_DATABASE:-database}"
      MYSQL_USER: "${DB_USERNAME:-user}"
      MYSQL_PASSWORD: "${DB_PASSWORD:-pass}"
      MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
    volumes:
     - db:/var/lib/mysql:delegated
    networks:
     - kool_local
  cache:
    image: redis:6-alpine
    volumes:
     - cache:/data:delegated
    networks:
     - kool_local

volumes:
  db:
  cache:

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"
`,
		"kool.yml": `scripts:
  artisan: kool exec app php artisan
  composer: kool exec app composer

  node: kool docker kooldev/node:14 node
  npm: kool docker kooldev/node:14 npm # can change to: yarn,pnpm

  mysql: kool exec database mysql -uroot -p$DB_PASSWORD

  setup:
    - cp .env.example .env
    - kool start
    - kool run composer install
    - kool run artisan key:generate
    - kool run npm install
    - kool run npm run dev

  reset:
    - kool run composer install
    - kool run artisan migrate:fresh --seed
    - kool run npm install
    - kool run npm run dev
`,
	}
	presets["nestjs"] = map[string]string{
		".dockerignore": `/node_modules
`,
		"preset_language": "javascript",
		"Dockerfile.build": `FROM kooldev/node:14-nest AS build

COPY . /app

RUN npm install

FROM kooldev/node:14-nest

COPY --from=build --chown=kool:kool /app /app

EXPOSE 3000

CMD [ "npm", "start" ]
`,
		"docker-compose.yml": `version: "3.7"
services:
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
    #  - $HOME/.ssh:/home/kool/.ssh:delegated
    networks:
     - kool_local
     - kool_global
#   database:
#     image: mysql:8.0
#     ports:
#      - "${KOOL_DATABASE_PORT:-3306}:3306"
#     environment:
#       MYSQL_ROOT_PASSWORD: "${DB_PASSWORD:-rootpass}"
#       MYSQL_DATABASE: "${DB_DATABASE:-database}"
#       MYSQL_USER: "${DB_USERNAME:-user}"
#       MYSQL_PASSWORD: "${DB_PASSWORD:-pass}"
#       MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
#     volumes:
#      - db:/var/lib/mysql:delegated
#     networks:
#      - kool_local
#   cache:
#     image: redis:6-alpine
#     volumes:
#      - cache:/data:delegated
#   mongo:
#     image: mongo:4
#     ports:
#      - "${KOOL_MONGO_PORT:-27017}:27017"
#     environment:
#       MONGO_INITDB_ROOT_USERNAME: "${MONGO_USERNAME:-root}"
#       MONGO_INITDB_ROOT_PASSWORD: "${MONGO_PASSWORD:-rootpass}"
#       MONGO_INITDB_DATABASE: "${MONGO_DATABASE:-database}"
#     volumes:
#      - mongo:/data/db:delegated
#     networks:
#      - kool_local
#
# volumes:
#   db:
#   cache:
#   mongo:

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"
`,
		"kool.yml": `scripts:
  node: kool exec app node
  npm: kool exec app npm # can change to: yarn,pnpm
  nest: kool exec app nest

  mysql: kool exec database mysql -uroot -prootpass

  mongo: kool exec mongo mongo -uroot -prootpass

  setup:
    - kool docker kooldev/node:14 npm install # can change to: yarn,pnpm
    - kool start
`,
	}
	presets["nextjs"] = map[string]string{
		".dockerignore": `/.next
/out
/build
/node_modules
`,
		"preset_language": "javascript",
		"Dockerfile.build": `FROM kooldev/node:14 AS build

COPY . /app

RUN npm install && npm run build

FROM kooldev/node:14

COPY --from=build --chown=kool:kool /app /app

EXPOSE 3000

CMD [ "npm", "start" ]
`,
		"docker-compose.yml": `version: "3.7"
services:
  app:
    image: kooldev/node:14
    command: ["npm", "run", "dev"]
    ports:
     - "${KOOL_APP_PORT:-3000}:3000"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:delegated
    #  - $HOME/.ssh:/home/kool/.ssh:delegated
    networks:
     - kool_local
     - kool_global

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"
`,
		"kool.yml": `scripts:
  node: kool exec app node
  npm: kool exec app npm # can change to: yarn,pnpm

  setup:
    - kool docker kooldev/node:14 npm install # can change to: yarn,pnpm
    - kool start
`,
	}
	presets["nextjs-static"] = map[string]string{
		".dockerignore": `/.next
/out
/build
/node_modules
`,
		"preset_language": "javascript",
		"Dockerfile.build": `FROM kooldev/node:14 AS node

COPY . /app

RUN npm install && npm run build && npm run export

FROM kooldev/http:static

ENV ROOT=/app/out

COPY --from=node /app/out /app/out
`,
		"docker-compose.yml": `version: "3.7"
services:
  app:
    image: kooldev/node:14
    command: ["npm", "run", "dev"]
    ports:
     - "${KOOL_APP_PORT:-3000}:3000"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:delegated
    #  - $HOME/.ssh:/home/kool/.ssh:delegated
    networks:
     - kool_local
     - kool_global

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"
`,
		"kool.yml": `scripts:
  node: kool exec app node
  npm: kool exec app npm # can change to: yarn,pnpm

  setup:
    - kool docker kooldev/node:14 npm install # can change to: yarn,pnpm
    - kool start
`,
	}
	presets["nuxtjs"] = map[string]string{
		".dockerignore": `/.nuxt
/dist
/node_modules
`,
		"preset_language": "javascript",
		"Dockerfile.build": `FROM kooldev/node:14 AS build

COPY . /app

RUN npm install && npm run build

FROM kooldev/node:14

COPY --from=build --chown=kool:kool /app /app

EXPOSE 3000

CMD [ "npm", "start" ]
`,
		"docker-compose.yml": `version: "3.7"
services:
  app:
    image: kooldev/node:14
    command: ["npm", "run", "dev"]
    ports:
     - "${KOOL_APP_PORT:-3000}:3000"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:delegated
    #  - $HOME/.ssh:/home/kool/.ssh:delegated
    networks:
     - kool_local
     - kool_global

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"
`,
		"kool.nuxt.config.js": `export default {
  server: {
    host: '0.0.0.0',
  }
}
`,
		"kool.yml": `scripts:
  node: kool exec app node
  npm: kool exec app npm # can change to: yarn,pnpm

  setup:
    - kool docker kooldev/node:14 npm install # can change to: yarn,pnpm
    - kool start
`,
	}
	presets["nuxtjs-static"] = map[string]string{
		".dockerignore": `/.nuxt
/dist
/node_modules
`,
		"preset_language": "javascript",
		"Dockerfile.build": `FROM kooldev/node:14 AS node

COPY . /app

RUN npm install && npm run build && npm run export

FROM kooldev/http:static

ENV ROOT=/app/dist

COPY --from=node /app/dist /app/dist
`,
		"docker-compose.yml": `version: "3.7"
services:
  app:
    image: kooldev/node:14
    command: ["npm", "run", "dev"]
    ports:
     - "${KOOL_APP_PORT:-3000}:3000"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:delegated
    #  - $HOME/.ssh:/home/kool/.ssh:delegated
    networks:
     - kool_local
     - kool_global

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"
`,
		"kool.nuxt.config.js": `export default {
  server: {
    host: '0.0.0.0',
  }
}
`,
		"kool.yml": `scripts:
  node: kool exec app node
  npm: kool exec app npm # can change to: yarn,pnpm

  setup:
    - kool docker kooldev/node:14 npm install # can change to: yarn,pnpm
    - kool start
`,
	}
	presets["symfony"] = map[string]string{
		".dockerignore": `/node_modules
/vendor
`,
		"preset_language": "php",
		"Dockerfile.build": `FROM kooldev/php:7.4 AS composer

COPY . /app
RUN composer install --no-interaction --prefer-dist --optimize-autoloader --quiet

FROM kooldev/node:14 AS node

COPY --from=composer /app /app
RUN yarn install && yarn prod

FROM kooldev/php:7.4-nginx

COPY --from=node --chown=kool:kool /app /app
`,
		"docker-compose.yml": `version: "3.7"
services:
  app:
    image: kooldev/php:7.4-nginx
    ports:
     - "${KOOL_APP_PORT:-80}:80"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:delegated
    #  - $HOME/.ssh:/home/kool/.ssh:delegated
    networks:
     - kool_local
     - kool_global
  database:
    image: mysql:8.0 # can change to: mysql:5.7
    command: --default-authentication-plugin=mysql_native_password # remove this line if you change to: mysql:5.7
    ports:
     - "${KOOL_DATABASE_PORT:-3306}:3306"
    environment:
      MYSQL_ROOT_PASSWORD: "${DB_PASSWORD:-rootpass}"
      MYSQL_DATABASE: "${DB_DATABASE:-database}"
      MYSQL_USER: "${DB_USERNAME:-user}"
      MYSQL_PASSWORD: "${DB_PASSWORD:-pass}"
      MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
    volumes:
     - db:/var/lib/mysql:delegated
    networks:
     - kool_local
  cache:
    image: redis:6-alpine
    volumes:
     - cache:/data:delegated
    networks:
     - kool_local

volumes:
  db:
  cache:

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"
`,
		"kool.yml": `scripts:
  console: kool exec app php ./bin/console
  phpunit: kool exec app php ./bin/phpunit
  composer: kool exec app composer

  node: kool docker kooldev/node:14 node
  npm: kool docker kooldev/node:14 npm # can change to: yarn,pnpm

  mysql: kool exec database mysql -uroot -prootpass

  setup:
    - kool start
    - cp .env.example .env
    - kool run composer install
`,
	}
	presets["wordpress"] = map[string]string{
		"preset_language": "php",
		"docker-compose.yml": `version: "3.7"
services:
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
  database:
    image: mysql:8.0 # can change to: mysql:5.7
    command: --default-authentication-plugin=mysql_native_password # remove this line if you change to: mysql:5.7
    ports:
     - "${KOOL_DATABASE_PORT:-3306}:3306"
    environment:
      MYSQL_ROOT_PASSWORD: "${DB_PASSWORD:-rootpass}"
      MYSQL_DATABASE: "${DB_DATABASE:-database}"
      MYSQL_USER: "${DB_USERNAME:-user}"
      MYSQL_PASSWORD: "${DB_PASSWORD:-pass}"
      MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
    volumes:
     - db:/var/lib/mysql:delegated
    networks:
     - kool_local
  cache:
    image: redis:6-alpine
    volumes:
     - cache:/data:delegated
    networks:
     - kool_local

volumes:
  db:
  cache:

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"
`,
		"kool.yml": `scripts:
  php: kool exec app php
  wp: kool exec app wp

  mysql: kool exec database mysql -uroot -p$DB_PASSWORD
`,
	}
	return presets
}
// GetTemplates get all templates
func GetTemplates() map[string]map[string]string {
	var templates = make(map[string]map[string]string)
	templates["app"] = map[string]string{
		"php74.yml": `image: kooldev/php:7.4-nginx
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
		"memcached16.yml": `image: memcached:1.6-alpine
volumes:
  - cache:/data:delegated
networks:
  - kool_local
`,
		"redis60.yml": `image: redis:6-alpine
volumes:
  - cache:/data:delegated
networks:
  - kool_local
`,
	}
	templates["database"] = map[string]string{
		"mysql57.yml": `image: mysql:5.7
ports:
  - "${KOOL_DATABASE_PORT:-3306}:3306"
environment:
  MYSQL_ROOT_PASSWORD: "${DB_PASSWORD:-rootpass}"
  MYSQL_DATABASE: "${DB_DATABASE:-database}"
  MYSQL_USER: "${DB_USERNAME:-user}"
  MYSQL_PASSWORD: "${DB_PASSWORD:-pass}"
  MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
volumes:
 - db:/var/lib/mysql:delegated
networks:
 - kool_local
`,
		"mysql80.yml": `image: mysql:8.0
command: --default-authentication-plugin=mysql_native_password
ports:
  - "${KOOL_DATABASE_PORT:-3306}:3306"
environment:
  MYSQL_ROOT_PASSWORD: "${DB_PASSWORD:-rootpass}"
  MYSQL_DATABASE: "${DB_DATABASE:-database}"
  MYSQL_USER: "${DB_USERNAME:-user}"
  MYSQL_PASSWORD: "${DB_PASSWORD:-pass}"
  MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
volumes:
 - db:/var/lib/mysql:delegated
networks:
 - kool_local
`,
		"prostgresql.yml": `image: postgres
ports:
  - "${KOOL_DATABASE_PORT:-3306}:3306"
environment:
  POSTGRES_DB: "${DB_DATABASE:-database}"
  POSTGRES_USER: "${DB_USERNAME:-user}"
  POSTGRES_PASSWORD: "${DB_PASSWORD:-pass}"
  POSTGRES_HOST_AUTH_METHOD: "trust"
volumes:
 - db:/var/lib/postgresql:data:delegated
networks:
 - kool_local
`,
	}
	templates["shared"] = map[string]string{
		"networks.yml": `kool_local:
kool_global:
  external: true
  name: "${KOOL_GLOBAL_NETWORK:-kool_global}"
`,
	}
	return templates
}
