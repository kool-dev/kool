#!/usr/bin/env bash
# Real local-only pushes and promotion; scratch configs do not execute ARM code.
set -euo pipefail
script=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)/release.sh
temporary=$(mktemp -d)
container=
cleanup() {
  [[ -z $container ]] || docker rm -f "$container" > /dev/null
  rm -rf -- "$temporary"
}
trap cleanup EXIT
container=$(docker create -p 127.0.0.1::5000 registry:2)
docker start "$container" > /dev/null
address=$(docker port "$container" 5000/tcp)
curl --fail --silent --retry 20 --retry-all-errors --retry-delay 1 "http://$address/v2/" > /dev/null
export GITHUB_REPOSITORY=kool-dev/docker-php GITHUB_REF=refs/heads/master GITHUB_EVENT_NAME=push
export GITHUB_RUN_ID=123 GITHUB_RUN_ATTEMPT=2
export RELEASE_CANONICAL_REPOSITORY=kool-dev/docker-php RELEASE_BRANCH=master
export RELEASE_REPOSITORY="$address/test/php" RELEASE_TAGS='["v1","v2"]'
export RELEASE_ARCHES='["amd64","arm64"]' RELEASE_METADATA="$temporary/metadata"
printf 'FROM scratch\nCOPY payload /payload\n' > "$temporary/Dockerfile"
for tag in v1 v2; do
  for arch in amd64 arm64; do
    printf '%s-%s\n' "$tag" "$arch" > "$temporary/payload"
    docker build --builder default --platform "linux/$arch" --provenance=false \
      --iidfile "$temporary/id" "$temporary" > /dev/null
    id=$(cat "$temporary/id")
    docker image save "$id" -o "$temporary/image.tar"
    config=$(tar -xOf "$temporary/image.tar" manifest.json | jq -r '.[0].Config')
    printf 'sha256:%s\n' "$(tar -xOf "$temporary/image.tar" "$config" | sha256sum | cut -d ' ' -f1)" \
      > "$temporary/config-$tag-$arch"
    export RELEASE_ARCH=$arch RELEASE_IMAGES
    RELEASE_IMAGES=$(jq -cn --arg tag "$tag" --arg id "$id" '[{tag:$tag,id:$id}]')
    bash "$script" stage
  done
done
bash "$script" publish
for tag in v1 v2; do
  manifest=$(docker buildx imagetools inspect --raw "$RELEASE_REPOSITORY:$tag")
  expected=$(jq -s 'sort_by(.arch) | map({digest,platform:{architecture:.arch,os:.os}})' \
    "$RELEASE_METADATA/$tag-amd64.json" "$RELEASE_METADATA/$tag-arm64.json")
  jq -e --argjson expected "$expected" '
    [.manifests[] | {digest,platform}] | sort_by(.platform.architecture) == $expected
  ' <<< "$manifest" > /dev/null
  for arch in amd64 arm64; do
    digest=$(jq -r .digest "$RELEASE_METADATA/$tag-$arch.json")
    config=$(docker buildx imagetools inspect --raw "$RELEASE_REPOSITORY@$digest" | jq -r .config.digest)
    [[ $config == "$(cat "$temporary/config-$tag-$arch")" ]]
  done
done
# Single-architecture repositories use the same index path, with one descriptor.
export RELEASE_TAGS='["v1"]' RELEASE_ARCHES='["amd64"]' RELEASE_METADATA="$temporary/single"
mkdir "$RELEASE_METADATA"
cp "$temporary/metadata/v1-amd64.json" "$RELEASE_METADATA/"
bash "$script" publish
manifest=$(docker buildx imagetools inspect --raw "$RELEASE_REPOSITORY:v1")
digest=$(jq -r .digest "$RELEASE_METADATA/v1-amd64.json")
jq -e --arg digest "$digest" '
  .manifests | length == 1 and .[0].digest == $digest and .[0].platform.architecture == "amd64"
' <<< "$manifest" > /dev/null

# Historical Node architecture aliases also contain exactly their tested image.
export RELEASE_TAGS='["v1","v2"]' RELEASE_ARCHES='["amd64","arm64"]'
export RELEASE_METADATA="$temporary/metadata" RELEASE_ARCH_ALIASES=true
bash "$script" publish
for suffix in x86_64 aarch64; do
  arch=amd64
  [[ $suffix != aarch64 ]] || arch=arm64
  manifest=$(docker buildx imagetools inspect --raw "$RELEASE_REPOSITORY:v1-$suffix")
  digest=$(jq -r .digest "$RELEASE_METADATA/v1-$arch.json")
  jq -e --arg digest "$digest" --arg arch "$arch" '
    .manifests | length == 1 and .[0].digest == $digest and .[0].platform.architecture == $arch
  ' <<< "$manifest" > /dev/null
done
echo "PASS: real staging, tested config IDs, multi/single-arch indexes and compatibility aliases"
