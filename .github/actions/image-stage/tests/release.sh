#!/usr/bin/env bash
set -euo pipefail
script=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)/release.sh
temporary=$(mktemp -d)
trap 'rm -rf -- "$temporary"' EXIT
mkdir "$temporary/bin"
export FAKE_LOG="$temporary/docker.log" FAKE_IMAGE="$temporary/image.json"
cat > "$temporary/bin/docker" <<'BASH'
#!/usr/bin/env bash
set -euo pipefail
printf '%s\n' "$*" >> "$FAKE_LOG"
case "$1 $2" in
  "image inspect") cat "$FAKE_IMAGE" ;;
  "push "*)
    case ${PUSH_MODE:-good} in
      missing) echo "no digest" ;;
      ambiguous) printf 'digest: sha256:%064d size: 1\n' 1 2 ;;
      failed) exit 1 ;;
      good) echo "tag: digest: sha256:$(printf %s "$2" | sha256sum | cut -d ' ' -f1) size: 1" ;;
    esac ;;
esac
BASH
chmod +x "$temporary/bin/docker"
export PATH="$temporary/bin:$PATH"
export GITHUB_REPOSITORY=kool-dev/docker-php GITHUB_REF=refs/heads/master GITHUB_EVENT_NAME=push
export GITHUB_RUN_ID=123 GITHUB_RUN_ATTEMPT=2
export RELEASE_CANONICAL_REPOSITORY=kool-dev/docker-php RELEASE_BRANCH=master
export RELEASE_REPOSITORY=kooldev/php RELEASE_TAGS='["v1","v2"]' RELEASE_ARCHES='["amd64","arm64"]'
export RELEASE_METADATA="$temporary/metadata"
id=$(printf 'sha256:%064d' 1)
export RELEASE_IMAGES
RELEASE_IMAGES=$(jq -cn --arg id "$id" '[{tag:"v1",id:$id},{tag:"v2",id:$id}]')

expect_fail() {
  if "$@" > "$temporary/error.log" 2>&1; then
    echo "Expected failure: $*" >&2
    exit 1
  fi
}
for arch in amd64 arm64; do
  export RELEASE_ARCH=$arch
  jq -n --arg id "$id" --arg arch "$arch" \
    '[{Id:$id,Os:"linux",Architecture:$arch,RootFS:{Layers:["base"]}}]' > "$FAKE_IMAGE"
  bash "$script" stage > /dev/null
done
[[ $(find "$RELEASE_METADATA" -type f | wc -l) == 4 ]]
! grep -Eq 'imagetools|pull|build ' "$FAKE_LOG"
cp -r "$RELEASE_METADATA" "$temporary/good"
mkdir -p "$temporary/nested/artifact-one" "$temporary/nested/artifact-two"
cp "$RELEASE_METADATA/"*-amd64.json "$temporary/nested/artifact-one/"
cp "$RELEASE_METADATA/"*-arm64.json "$temporary/nested/artifact-two/"
RELEASE_METADATA="$temporary/nested" bash "$script" validate > /dev/null
restore() {
  rm -rf -- "$RELEASE_METADATA"
  cp -r "$temporary/good" "$RELEASE_METADATA"
  : > "$FAKE_LOG"
}
mutate() {
  jq "$1" "$RELEASE_METADATA/v2-arm64.json" > "$temporary/changed.json"
  mv "$temporary/changed.json" "$RELEASE_METADATA/v2-arm64.json"
}
bash "$script" validate > /dev/null
: > "$FAKE_LOG"
bash "$script" publish
[[ $(wc -l < "$FAKE_LOG") == 2 ]]
grep -q 'kooldev/php:v1 kooldev/php@sha256:' "$FAKE_LOG"
! grep -q 'inspect' "$FAKE_LOG"

for filter in '.repository="foreign/php"' '.run_attempt="1"' '.arch="s390x"' \
  '.os="windows"' '.image_id="bad"' '.digest="sha256:bad"' '.extra=true' \
  '.tag="v1"' 'del(.digest)' '[]'; do
  restore
  mutate "$filter"
  expect_fail bash "$script" publish
  [[ ! -s "$FAKE_LOG" ]]
done
restore
rm "$RELEASE_METADATA/v2-arm64.json"
expect_fail bash "$script" publish
[[ ! -s "$FAKE_LOG" ]]
restore
cp "$RELEASE_METADATA/v2-arm64.json" "$RELEASE_METADATA/duplicate.json"
expect_fail bash "$script" publish
[[ ! -s "$FAKE_LOG" ]]
restore
digest=$(jq -r .digest "$RELEASE_METADATA/v1-amd64.json")
mutate ".digest=\"$digest\""
expect_fail bash "$script" publish
[[ ! -s "$FAKE_LOG" ]]
restore
printf '{' > "$RELEASE_METADATA/v2-arm64.json"
expect_fail bash "$script" publish
[[ ! -s "$FAKE_LOG" ]]
restore
for event in pull_request pull_request_target; do
  expect_fail env GITHUB_EVENT_NAME="$event" bash "$script" stage
  expect_fail env GITHUB_EVENT_NAME="$event" bash "$script" publish
done
expect_fail env GITHUB_REF=refs/heads/feature bash "$script" stage
expect_fail env GITHUB_REPOSITORY=fork/docker-php bash "$script" publish
expect_fail env RELEASE_TAGS='["v1","v1"]' bash "$script" publish
expect_fail env RELEASE_ARCHES='["s390x"]' bash "$script" publish
[[ ! -s "$FAKE_LOG" ]]
for mode in missing ambiguous failed; do
  expect_fail env PUSH_MODE="$mode" RELEASE_METADATA="$temporary/no-metadata-$mode" bash "$script" stage
  [[ ! -d "$temporary/no-metadata-$mode" ]]
done
: > "$FAKE_LOG"
expect_fail env RELEASE_ARCH=amd64 bash "$script" stage # fixture is ARM64
[[ $(wc -l < "$FAKE_LOG") == 1 ]]
expect_fail env RELEASE_IMAGES='[{"tag":"bad;tag","id":"bad"}]' bash "$script" stage

# Parent checks work without release credentials/context, including on PRs.
export GITHUB_EVENT_NAME=pull_request
bash "$script" parent local/base "$id" "$id"
expect_fail bash "$script" parent local/base "$(printf 'sha256:%064d' 2)" "$id"
jq '.[0].RootFS.Layers=[]' "$FAKE_IMAGE" > "$temporary/empty.json"
mv "$temporary/empty.json" "$FAKE_IMAGE"
expect_fail bash "$script" parent local/base "$id" "$id"
echo "PASS: shared Bash release guards, exact IDs, complete matrix and parent contracts"
