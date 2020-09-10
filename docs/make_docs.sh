#!/bin/bash

rm -f 4-Commands/*.md

go run docs.go

rm 4-Commands/kool_deploy.md

mv 4-Commands/kool.md 4-Commands/0-kool.md
mv 4-Commands/kool_init.md 4-Commands/kool-init.md
mv 4-Commands/kool_start.md 4-Commands/kool-start.md
mv 4-Commands/kool_stop.md 4-Commands/kool-stop.md
mv 4-Commands/kool_exec.md 4-Commands/kool-exec.md
mv 4-Commands/kool_run.md 4-Commands/kool-run.md
mv 4-Commands/kool_docker.md 4-Commands/kool-docker.md
mv 4-Commands/kool_status.md 4-Commands/kool-status.md
mv 4-Commands/kool_info.md 4-Commands/kool-info.md
mv 4-Commands/kool_self-update.md 4-Commands/kool-self--update.md
mv 4-Commands/kool_logs.md 4-Commands/kool-logs.md
