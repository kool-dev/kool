---
name: kool
description: Discovers and runs project scripts with the kool CLI. Use when the user wants to list custom commands from kool.yml, understand project tasks, or run kool scripts. Lists scripts in JSON for agents and reads kool.yml for context.
---

# Kool CLI

Use kool to discover and run project scripts defined in `kool.yml`.

## Quick Start

```bash
kool scripts --json
kool run <script>
```

## Core Workflow

1. Work from the project root (has `kool.yml`) or use `-w` to point to it.
2. Discover scripts:

```bash
kool scripts --json   # Preferred: returns [{name, comments, commands}]
kool scripts          # Human-readable list
```

3. Run scripts:

```bash
kool run <script>
kool run <script> -- <args>
```

## Important Rules

- **ALWAYS** run commands from the project root or use `-w/--working_dir`.
- **ONLY** pass extra args to single-line scripts; multi-line scripts reject extra args.
- **REMEMBER** `kool.yml` scripts are not full bash (no pipes/conditionals); use `kool docker <image> bash -c "..."` for complex shell logic.
- **CHECK** `kool.yml` or `kool.yaml` in the project root and `~/kool` for shared scripts.
