#!/bin/bash

bash parse_presets.sh

go build -o /usr/local/bin/kool

# docker run --rm --env GOOS=linux --env GOARCH=amd64 -v $(pwd):/code -w /code golang:1.14 go build -tags netgo -ldflags '-extldflags "-static"' -o dist/kool-linux-amd64
