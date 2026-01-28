---
name: kool-cli
description: Docker development environment CLI. Use for managing containers (start/stop/restart), executing commands in services, viewing logs, and running project scripts from kool.yml.
---

# Kool CLI

Kool simplifies Docker-based development with commands for container lifecycle, service execution, and custom scripts.

## Quick Reference

```bash
kool start                    # Start all services from docker-compose.yml
kool stop                     # Stop all services
kool restart --rebuild        # Restart and rebuild images
kool status                   # Show running containers
kool exec <service> <cmd>     # Run command in service container
kool logs -f <service>        # Follow service logs
kool run --json               # List available scripts as JSON
kool run <script>             # Run a script from kool.yml
```

## Service Lifecycle

Services are defined in `docker-compose.yml`. Kool wraps docker-compose with simpler commands.

```bash
kool start                    # Start all services
kool start app database       # Start specific services
kool start --rebuild          # Rebuild images before starting
kool start --foreground       # Run in foreground (see logs)
kool start --profile worker   # Enable a docker-compose profile

kool stop                     # Stop all services
kool stop app                 # Stop specific service
kool stop --purge             # Stop and remove volumes (destructive)

kool restart                  # Restart all services
kool restart --rebuild        # Rebuild images on restart
kool restart --purge          # Purge volumes on restart

kool status                   # Show status of all containers
```

## Executing Commands in Containers

Use `exec` to run commands inside running service containers (like SSH).

```bash
kool exec <service> <command>
kool exec app bash                      # Interactive shell
kool exec app php artisan migrate       # Run Laravel migration
kool exec app npm install               # Install npm packages
kool exec -e VAR=value app env          # With environment variable
```

## One-off Docker Containers

Use `docker` to run commands in temporary containers (not services).

```bash
kool docker <image> <command>
kool docker node npm init                           # Run npm in node container
kool docker --volume=$PWD:/app golang go build      # Mount current dir
kool docker --env=DEBUG=1 python python script.py  # With env var
kool docker --publish=3000:3000 node npm start      # Expose port
```

## Viewing Logs

```bash
kool logs                     # Last 25 lines from all services
kool logs app                 # Logs from specific service
kool logs -f                  # Follow logs (live)
kool logs -f app worker       # Follow multiple services
kool logs --tail 100          # Last 100 lines
kool logs --tail 0            # All logs
```

## Project Scripts

Scripts are defined in `kool.yml` and provide project-specific commands.

```bash
kool run --json               # List scripts as JSON [{name, comments, commands}]
kool run                      # List scripts (human-readable)
kool run <script>             # Run a script
kool run <script> -- <args>   # Pass args (single-line scripts only)
kool run -e VAR=1 <script>    # Run with environment variable
```

## Global Options

All commands support:

```bash
-w, --working_dir <path>      # Run from different directory
--verbose                     # Increase output verbosity
```

## Important Rules

- **ALWAYS** run from project root (has `docker-compose.yml` and `kool.yml`) or use `-w`.
- **Service names** come from `docker-compose.yml` service definitions.
- **Script args** only work with single-line scripts; multi-line scripts reject extra args.
- **Scripts** in `kool.yml` are not full bash - use `kool docker <image> bash -c "..."` for pipes/conditionals.
