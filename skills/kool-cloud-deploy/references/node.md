# Node API (Express, NestJS) on Kool Cloud

The TypeScript build runs natively. The runtime stage installs only production dependencies, for the cluster's arch, and runs as the non-root `node` user.

```yaml
# kool.cloud.yml
services:
  app:
    build: .
    expose: 3000
    public: true
    size: small
    environment:
      NODE_ENV: production
      PORT: "3000"                           # the app must listen on it, on 0.0.0.0
      DATABASE_URL: "{{DATABASE_URL}}"       # set DATABASE_URL in the Kool panel
    readinessProbe:                          # optional; any route that answers 200 fast
      httpGet: { path: /health, port: 3000 }
    after:
      - [npx, prisma, migrate, deploy]       # or your ORM's migrate command; drop if none
    daemons:
      - command: [node, dist/worker.js]      # optional background worker, same image
```

```dockerfile
# Dockerfile
# build: native arch, because the compiled JS is arch-independent.
FROM --platform=$BUILDPLATFORM node:24-bookworm-slim AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build                            # nest build / tsc -> dist/

# runtime: the cluster's arch. npm ci here gets the native modules right.
FROM node:24-bookworm-slim
WORKDIR /app
ENV NODE_ENV=production
COPY package.json package-lock.json ./
RUN npm ci --omit=dev && npm cache clean --force
COPY --from=build /app/dist ./dist
USER node
EXPOSE 3000
CMD ["node", "dist/main.js"]
```
- **Plain JavaScript Express, with no build step:** drop the first stage, replace the `COPY --from=build` line with `COPY --chown=node:node . .`, and use `CMD ["node", "server.js"]`.
- **Migration CLIs** (`prisma`, `knex`, `typeorm`, `drizzle-kit`) must be in `dependencies`, because `--omit=dev` drops devDependencies. Prisma also needs `COPY prisma ./prisma` and `RUN npx prisma generate` in the runtime stage.

```gitignore
# .dockerignore
.git
.env*
node_modules
dist
coverage
```

## Notes
- **Listening:** use `app.listen(process.env.PORT ?? 3000)`, or `'0.0.0.0'` as the host. A host of `localhost` or `127.0.0.1` fails the liveness probe, and the deploy rolls back after 5 minutes. Nest's default, `await app.listen(3000)`, is fine.
- **ENTRYPOINT:** keep the node image's `docker-entrypoint.sh` and use `CMD`. Hooks, daemons and cron commands are passed to it as arguments, and it runs them.
- **Databases:** `cloud: true` gives managed MySQL 8 only (see laravel.md). For Postgres or anything else, use an external `DATABASE_URL`, kept as a panel variable through `"{{DATABASE_URL}}"` or in a git-ignored env file (`environment: kool.deploy.env`).
- **Persisted paths:** a fresh `persists` path may mount root-owned. If the app writes there while running as `USER node`, check that it can write (`kool cloud exec ... touch <path>/x`), or drop `USER` for that service, as the proven Next.js setup does.
- **Workers and scheduled jobs:** a queue consumer is a daemon (`command: [node, dist/worker.js]`, logs with `-c daemon0`). A scheduled job is a `cron` entry that runs a script and exits; don't keep a `node-cron` loop inside the web process if `replicas` is above 1.
- **Verify:** `curl -sS https://<domain>/health`, then one real endpoint with its expected JSON, for example `curl -sS -X POST -H 'content-type: application/json' -d '{...}' https://<domain>/api/...`.
