FROM docker/compose:alpine-1.26.2 AS docker-compose
FROM golang:alpine3.12 AS build

WORKDIR /app

COPY . /app

RUN go build -o kool

FROM alpine:3.12

COPY --from=docker-compose /usr/local/bin/docker /usr/local/bin/docker
COPY --from=docker-compose /usr/local/bin/docker-compose /usr/local/bin/docker-compose
COPY --from=build /app/kool /usr/local/bin/kool

CMD [ "kool" ]
