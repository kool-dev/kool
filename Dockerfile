FROM docker/compose:alpine-1.28.2 AS docker-compose
FROM golang:1.16 AS build

ARG BUILD_VERSION=0.0.0-auto

WORKDIR /app

COPY . /app

RUN go build -a \
	-tags 'osusergo netgo static_build' \
	-ldflags '-X kool-dev/kool/cmd.version='$BUILD_VERSION' -extldflags "-static"' \
	-o kool

FROM alpine:3.12

ENV DOCKER_HOST=tcp://docker:2375

COPY --from=docker-compose /usr/local/bin/docker /usr/local/bin/docker
COPY --from=docker-compose /usr/local/bin/docker-compose /usr/local/bin/docker-compose
COPY --from=build /app/kool /usr/local/bin/kool

RUN apk add --no-cache git bash

CMD [ "kool" ]
