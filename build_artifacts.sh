#!/bin/bash

if [ -f .env ]; then
    source .env
fi

GO_IMAGE=${GO_IMAGE:-golang:1.15.0}

rm -rf dist
mkdir -p dist

echo "Building to GOOS=darwin GOARCH=amd64"

docker run --rm --env GOOS=darwin --env GOARCH=amd64 --env CGO_ENABLED=0 -v $(pwd):/code -w /code $GO_IMAGE go build -a -tags 'osusergo netgo static_build' -ldflags '-extldflags "-static"' -o dist/kool-darwin-x86_64

echo "Building to GOOS=linux GOARCH=amd64"

docker run --rm --env GOOS=linux --env GOARCH=amd64 --env CGO_ENABLED=0 -v $(pwd):/code -w /code $GO_IMAGE go build -a -tags 'osusergo netgo static_build' -ldflags '-extldflags "-static"' -o dist/kool-linux-x86_64

echo "Building to GOOS=linux GOARCH=arm GOARM=6"

docker run --rm --env GOOS=linux --env GOARCH=arm --env GOARM=6 --env CGO_ENABLED=0 -v $(pwd):/code -w /code $GO_IMAGE go build -a -tags 'osusergo netgo static_build' -ldflags '-extldflags "-static"' -o dist/kool-linux-arm6

echo "Building to GOOS=linux GOARCH=arm GOARM=7"

docker run --rm --env GOOS=linux --env GOARCH=arm --env GOARM=7 --env CGO_ENABLED=0 -v $(pwd):/code -w /code $GO_IMAGE go build -a -tags 'osusergo netgo static_build' -ldflags '-extldflags "-static"' -o dist/kool-linux-arm7

echo "Going to generate CHECKSUMS"

for file in dist/*; do
	shasum $file | awk '{print $1}' > $file.checksum
done
