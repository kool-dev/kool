#!/bin/bash

go run docs.go

rm 4-Commands/kool_deploy.md

mv 4-Commands/kool.md 4-Commands/0-kool.md
mv 4-Commands/kool_init.md 4-Commands/1-kool-init.md
mv 4-Commands/kool_start.md 4-Commands/2-kool-start.md
mv 4-Commands/kool_stop.md 4-Commands/3-kool-stop.md
mv 4-Commands/kool_exec.md 4-Commands/4-kool-exec.md
mv 4-Commands/kool_run.md 4-Commands/5-kool-run.md
mv 4-Commands/kool_docker.md 4-Commands/6-kool-docker.md
mv 4-Commands/kool_status.md 4-Commands/7-kool-status.md
mv 4-Commands/kool_info.md 4-Commands/8-kool-info.md
mv 4-Commands/kool_self-update.md 4-Commands/9-kool-self--update.md
