# Next.js on Kool Cloud (including Prisma + SQLite)

This shape is proven on Kool: Next.js 16 App Router, plus Prisma 7 with SQLite through `@prisma/adapter-better-sqlite3`. Without a database, drop the `prisma` lines and the `persists` and `DATABASE_URL` entries.

```yaml
# kool.cloud.yml
services:
  app:
    build: .
    expose: 3000
    public: true
    size: small                              # micro (256Mi) is tight for `next start`
    environment:
      NODE_ENV: production
      NEXT_TELEMETRY_DISABLED: "1"
      DATABASE_URL: "file:/app/data/prod.db"   # SQLite on the persisted disk
    persists:
      size: 1Gi
      paths:
        - path: /app/data                    # never `sync: true` here: it would overwrite the db
```

```dockerfile
# Dockerfile
# builder: the host's native arch. Turbopack's `next build` segfaults under QEMU on Apple
# Silicon, and .next plus the generated Prisma client are arch-independent.
FROM --platform=$BUILDPLATFORM node:24-bookworm-slim AS builder
RUN apt-get update && apt-get install -y --no-install-recommends python3 make g++ openssl ca-certificates \
 && rm -rf /var/lib/apt/lists/*
WORKDIR /app
ENV NEXT_TELEMETRY_DISABLED=1 DATABASE_URL=file:/tmp/build.db
COPY package.json package-lock.json prisma.config.ts ./
COPY prisma ./prisma
RUN npm ci
COPY . .
RUN npx prisma generate && npm run build

# runner: the cluster's arch (amd64). npm ci here installs native modules
# (better-sqlite3, sharp) for that arch.
FROM node:24-bookworm-slim AS runner
RUN apt-get update && apt-get install -y --no-install-recommends python3 make g++ openssl ca-certificates \
 && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY package.json package-lock.json prisma.config.ts ./
COPY prisma ./prisma
RUN npm ci
COPY . .
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/src/generated ./src/generated
RUN mkdir -p /app/data
ENV NODE_ENV=production
EXPOSE 3000
CMD ["sh", "-c", "npx prisma migrate deploy && exec npx next start -H 0.0.0.0 -p 3000"]
```
- `src/generated` is the Prisma 7 generator `output`; use the path from your schema.
- A config file with another name needs `--config <file>` on the prisma commands.

```gitignore
# .dockerignore
.git
.env*
node_modules
.next
out
*.db
*.db-journal
```
Without this, the host's `node_modules` (macOS/arm64 binaries) and `.env` overwrite the image's.

## Notes
- **Migrations:** they run in `CMD`, before `next start`. The SQLite file lives inside this pod's disk, and `persists` makes each deploy stop-then-start, so there is only ever one writer. `before: [[npx, prisma, migrate, deploy]]` also works, the first deploy included.
- **Replicas:** keep one. SQLite can't be shared across pods; for more replicas, use a server database (Postgres or MySQL) through `DATABASE_URL`.
- **Prisma 7:** `prisma migrate deploy` needs `datasource.url` in `prisma.config.ts` (for example `env("DATABASE_URL")`), and Prisma no longer reads `.env` by itself. Keep `.env` out of the image, or its local `DATABASE_URL` wins at build time.
- **Native build tools:** `better-sqlite3` needs `python3 make g++` when no prebuilt binary matches. Always run `npm ci` inside the image, never copy `node_modules` in.
- **Dependencies in the runner:** it uses a full `npm ci` because the `prisma` CLI is usually a devDependency. Move `prisma` to `dependencies` to use `npm ci --omit=dev`.
- **`NEXT_PUBLIC_*` values are inlined at build time.** Pass them as build args (see spa.md); runtime `environment:` won't change them.
- **Build-time `DATABASE_URL`:** the builder's value is a throwaway that only lets `prisma generate` and `next build` load the config. Pages that query the database at build time must be dynamic.
- **Bind to 0.0.0.0:** `-H 0.0.0.0` keeps `next start` listening on all interfaces whatever `HOSTNAME` holds; Kubernetes may set it to the pod name.
- **Verify:** `curl -sS https://<domain>/ | grep -i '<title>'`. Then, if the app has a write API, POST to it and GET the result back, which proves the database is writable. To look at the file: `kool cloud exec --domain <d> app ls -la /app/data`.
