FROM docker/compose:alpine-1.26.2 AS docker-compose
FROM golang:1.14 AS build

WORKDIR /app

COPY . /app

RUN go build -tags 'osusergo netgo static_build' -ldflags '-extldflags "-static"' -o kool

FROM alpine:3.12

ENV DOCKER_HOST=tcp://docker:2375

COPY --from=docker-compose /usr/local/bin/docker /usr/local/bin/docker
COPY --from=docker-compose /usr/local/bin/docker-compose /usr/local/bin/docker-compose
COPY --from=build /app/kool /usr/local/bin/kool

CMD [ "kool" ]
