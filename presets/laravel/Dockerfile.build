FROM fireworkweb/php:7.4 as composer

COPY . /app
RUN composer install --no-interaction --prefer-dist --optimize-autoloader --quiet

FROM fireworkweb/node:14 as node

COPY --from=composer /app /app
RUN yarn install && yarn prod

FROM fireworkweb/php:7.4-nginx

COPY --chown=fwd:fwd --from=node /app /app
