FROM golang:1.21 AS build

ARG BUILD_VERSION=0.0.0-auto

WORKDIR /app

COPY . /app

RUN go build -a \
	-tags 'osusergo netgo static_build' \
	-ldflags '-X kool-dev/kool/commands.version='$BUILD_VERSION' -extldflags "-static"' \
	-o kool

FROM docker:20.10.21-cli

ENV DOCKER_HOST=tcp://docker:2375

COPY --from=build /app/kool /usr/local/bin/kool

RUN apk add --no-cache git bash \
	&& rm -rf /var/cache/apk/* /tmp/*

CMD [ "kool" ]
