package cmd

// auto generated file

var presets = make(map[string]map[string]string)
func init() {
	presets["adonis-nextjs"] = map[string]string{
		"Dockerfile.adonis.build": `FROM fireworkweb/node:14-adonis

COPY . /app

RUN npm install

EXPOSE 3333

CMD [ "npm", "adonis:start" ]`,
		"Dockerfile.nextjs.build": `FROM fireworkweb/node:14

COPY . /app

RUN npm install && npm run build

EXPOSE 3000

CMD [ "npm", "nextjs:start" ]`,
		"docker-compose.yml": `version: "3.8"
services:
  adonis:
    image: fireworkweb/node:14-adonis
    command: ["adonis", "serve", "--dev"]
    ports:
     - "${PORT:-3333}:3333"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:cached
    #  - $HOME/.ssh:/home/fwd/.ssh:cached
    networks:
     - kool_local
     - kool_global
  nextjs:
    image: fireworkweb/node:14
    command: ["npm", "run", "nextjs:dev"]
    ports:
     - "${KOOL_NEXTJS_PORT:-3000}:3000"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:cached
    #  - $HOME/.ssh:/home/fwd/.ssh:cached
    networks:
     - kool_local
     - kool_global

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
#     networks:
#      - kool_local
#
# volumes:
#   db:
#   cache:

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"`,
		"kool.env": `HOST=0.0.0.0
PORT=3333
NODE_ENV=development
APP_URL=http://localhost:${PORT}`,
		"kool.package.json": `{
  "name": "adonis-nextjs",
  "private": true,
  "description": "AdonisJS and NextJS",
  "scripts": {
    "adonis:start": "node server.js",
    "adonis:test": "node ace test",
    "nextjs:dev": "next",
    "nextjs:build": "next build",
    "nextjs:start": "next start"
  },
  "dependencies": {
    "@adonisjs/ace": "^5.0.8",
    "@adonisjs/auth": "^3.0.7",
    "@adonisjs/bodyparser": "^2.0.5",
    "@adonisjs/cors": "^1.0.7",
    "@adonisjs/fold": "^4.0.9",
    "@adonisjs/framework": "^5.0.9",
    "@adonisjs/ignitor": "^2.0.8",
    "@adonisjs/lucid": "^6.1.3",
    "next": "^9.4.4",
    "react": "^16.13.1",
    "react-dom": "^16.13.1"
  },
  "devDependencies": {},
  "autoload": {
    "App": "./app"
  }
}`,
		"kool.yml": `scripts:
  node: kool exec app node
  npm: kool exec app npm # can change to: yarn,pnpm
  adonis: kool exec adonis adonis

  install:
    - cp .env.example .env
    - kool run --docker fireworkweb/node:14 npm install # can change to: yarn,pnpm
    - kool start`,
	}
	presets["adonis"] = map[string]string{
		"Dockerfile.build": `FROM fireworkweb/node:14-adonis

COPY . /app

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
    #  - $HOME/.ssh:/home/fwd/.ssh:cached
    networks:
     - kool_local
     - kool_global

##########################################################################
# optionally you can enable mysql, redis or other services to your stack #
##########################################################################
#   database:
#     image: mysql:8.0 # possibly change to: mysql:5.7
#     ports:
#      - "${KOOL_DATABASE_PORT:-3306}:3306"
#     environment:
#       MYSQL_ROOT_PASSWORD: "${DB_PASSWORD:-rootpass}"
#       MYSQL_DATABASE: "${DB_DATABASE:-database}"
#       MYSQL_USER: "${DB_USERNAME:-user}"
#       MYSQL_PASSWORD: "${DB_PASSWORD:-}"
#       MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
#     volumes:
#      - db:/var/lib/mysql:cached
#     networks:
#      - kool_local
#   cache:
#     image: redis:alpine
#     volumes:
#      - cache:/data:cached
#
# volumes:
#   db:
#   cache:

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"`,
		"kool.env": `HOST=0.0.0.0
PORT=3333
NODE_ENV=development
APP_URL=http://localhost:${PORT}`,
		"kool.yml": `scripts:
  node: kool exec app node
  npm: kool exec app npm # can change to: yarn,pnpm
  adonis: kool exec app adonis

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
     - "${KOOL_HTTP_PORT:-80}:80"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:cached
    #  - $HOME/.ssh:/home/fwd/.ssh:cached
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
      MYSQL_PASSWORD: "${DB_PASSWORD:-}"
      MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
    volumes:
     - db:/var/lib/mysql:cached
    networks:
     - kool_local
     - kool_global
  cache:
    image: redis:alpine
    volumes:
     - cache:/data:cached

volumes:
  db:
  cache:

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"`,
		"kool.yml": `scripts:
  php: kool exec app php
  composer: kool exec app composer

  node: kool run --docker fireworkweb/node:14 node
  npm: kool run --docker fireworkweb/node:14 npm # can change to: yarn,pnpm

  install:
    - kool start
    - cp .env.example .env
    - kool run composer install
    - kool run php artisan key:generate
    - kool run npm install
    - kool run npm run dev

  reset:
    - kool run composer install
    - kool run php artisan migrate:fresh --seed
    - kool run npm install
    - kool run npm run dev`,
	}
	presets["nextjs-static"] = map[string]string{
		"Dockerfile.build": `FROM fireworkweb/node:14 as node

COPY . /app

RUN npm install && npm run build && npm run export

FROM fireworkweb/http:static

ENV ROOT=/app/out

COPY --from=node /app/out /app/out`,
		"docker-compose.yml": `version: "3.8"
services:
  app:
    image: fireworkweb/node:14
    command: ["npm", "run", "dev"]
    ports:
     - "${KOOL_APP_PORT:-3000}:3000"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:cached
    #  - $HOME/.ssh:/home/fwd/.ssh:cached
    networks:
     - kool_local
     - kool_global

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"`,
		"kool.yml": `scripts:
  node: kool exec app node
  npm: kool exec app npm # can change to: yarn,pnpm

  install:
    - kool run --docker fireworkweb/node:14 npm install # can change to: yarn,pnpm
    - kool start`,
	}
	presets["nextjs"] = map[string]string{
		"Dockerfile.build": `FROM fireworkweb/node:14

COPY . /app

RUN npm install && npm run build

EXPOSE 3000

CMD [ "npm", "start" ]`,
		"docker-compose.yml": `version: "3.8"
services:
  app:
    image: fireworkweb/node:14
    command: ["npm", "run", "dev"]
    ports:
     - "${KOOL_APP_PORT:-3000}:3000"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:cached
    #  - $HOME/.ssh:/home/fwd/.ssh:cached
    networks:
     - kool_local
     - kool_global

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"`,
		"kool.yml": `scripts:
  node: kool exec app node
  npm: kool exec app npm # can change to: yarn,pnpm

  install:
    - kool run --docker fireworkweb/node:14 npm install # can change to: yarn,pnpm
    - kool start`,
	}
	presets["nuxtjs-static"] = map[string]string{
		"Dockerfile.build": `FROM fireworkweb/node:14 as node

COPY . /app

RUN npm install && npm run build && npm run export

FROM fireworkweb/http:static

ENV ROOT=/app/dist

COPY --from=node /app/dist /app/dist`,
		"docker-compose.yml": `version: "3.8"
services:
  app:
    image: fireworkweb/node:14
    command: ["npm", "run", "dev"]
    ports:
     - "${KOOL_APP_PORT:-3000}:3000"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:cached
    #  - $HOME/.ssh:/home/fwd/.ssh:cached
    networks:
     - kool_local
     - kool_global

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"`,
		"kool.nuxt.config.js": `export default {
  server: {
    host: '0.0.0.0',
  }
}`,
		"kool.yml": `scripts:
  node: kool exec app node
  npm: kool exec app npm # can change to: yarn,pnpm

  install:
    - kool run --docker fireworkweb/node:14 npm install # can change to: yarn,pnpm
    - kool start`,
	}
	presets["nuxtjs"] = map[string]string{
		"Dockerfile.build": `FROM fireworkweb/node:14

COPY . /app

RUN npm install && npm run build

EXPOSE 3000

CMD [ "npm", "start" ]`,
		"docker-compose.yml": `version: "3.8"
services:
  app:
    image: fireworkweb/node:14
    command: ["npm", "run", "dev"]
    ports:
     - "${KOOL_APP_PORT:-3000}:3000"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:cached
    #  - $HOME/.ssh:/home/fwd/.ssh:cached
    networks:
     - kool_local
     - kool_global

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"`,
		"kool.nuxt.config.js": `export default {
  server: {
    host: '0.0.0.0',
  }
}`,
		"kool.yml": `scripts:
  node: kool exec app node
  npm: kool exec app npm # can change to: yarn,pnpm

  install:
    - kool run --docker fireworkweb/node:14 npm install # can change to: yarn,pnpm
    - kool start`,
	}
}
