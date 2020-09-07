#!/bin/bash

if [ -f .env ]; then
    source .env
fi

GO_IMAGE=${GO_IMAGE:-golang:1.15.0}

rm -rf dist
mkdir -p dist

BUILD=( \
  "dist/kool-darwin-x86_64|--env GOOS=darwin --env GOARCH=amd64" \
  "dist/kool-linux-x86_64|--env GOOS=linux --env GOARCH=amd64" \
  "dist/kool-linux-arm6|--env GOOS=linux --env GOARCH=arm --env GOARM=6" \
  "dist/kool-linux-arm7|--env GOOS=linux --env GOARCH=arm --env GOARM=7" \
  "dist/kool.exe|--env GOOS=windows --env GOARCH=amd64" \
)

for i in "${!BUILD[@]}"; do
    dist=$(echo ${BUILD[$i]} | cut -d'|' -f1)
    flags=$(echo ${BUILD[$i]} | cut -d'|' -f2)
    echo "Building to ${flags}"
    docker run --rm \
        $flags \
        --env CGO_ENABLED=0 \
        -v $(pwd):/code -w /code $GO_IMAGE \
        go build -a -tags 'osusergo netgo static_build' \
        -ldflags '-extldflags "-static"' \
        -o $dist
done

echo "Building kool-install.exe"

docker run --rm -i \
    -v $(pwd):/work \
    amake/innosetup /dApplicationVersion=1.0.16 inno-setup/kool.iss
mv inno-setup/Output/mysetup.exe dist/kool-install.exe

echo "Going to generate CHECKSUMS"

for file in dist/*; do
    shasum $file | awk '{print $1}' > $file.checksum
done
