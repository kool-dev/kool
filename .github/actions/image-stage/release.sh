#!/usr/bin/env bash
# Shared image release boundary. Requires Bash, jq and Docker; never rebuilds.
set -euo pipefail

die() { echo "::error::$*" >&2; exit 1; }

scope() {
  [[ ${GITHUB_RUN_ID:-} =~ ^[1-9][0-9]*$ && ${GITHUB_RUN_ATTEMPT:-} =~ ^[1-9][0-9]*$ ]] ||
    die "Invalid run/attempt"
  [[ ${RELEASE_REPOSITORY:-} =~ ^[a-z0-9.-]+(:[0-9]+)?/[a-z0-9._/-]+$ ]] ||
    die "Invalid registry repository"
}

guard() {
  scope
  [[ -n ${RELEASE_CANONICAL_REPOSITORY:-} && -n ${RELEASE_BRANCH:-} &&
     ${GITHUB_REPOSITORY:-} == "$RELEASE_CANONICAL_REPOSITORY" &&
     ${GITHUB_REF:-} == "refs/heads/$RELEASE_BRANCH" ]] ||
    die "Publication requires the canonical repository's default branch"
  case ${GITHUB_EVENT_NAME:-} in
    push|schedule|workflow_dispatch|repository_dispatch) ;;
    *) die "This event cannot publish images" ;;
  esac
}

stage() {
  guard
  [[ ${RELEASE_ARCH:-} == amd64 || ${RELEASE_ARCH:-} == arm64 ]] ||
    die "Unsupported native architecture"
  # Validate the entire staging batch before pushing its first image.
  jq -e '
    type == "array" and length > 0 and
    (length == (unique_by(.tag) | length)) and
    all(.[]; type == "object" and (keys == ["id", "tag"]) and
      (.tag | type == "string" and test("^[A-Za-z0-9_][A-Za-z0-9_.-]{0,127}$")) and
      (.id | type == "string" and test("^sha256:[0-9a-f]{64}$")))
  ' <<< "${RELEASE_IMAGES:-}" >/dev/null || die "Invalid tested image batch"
  local tag id info staging output digest directory=${RELEASE_METADATA:-release-metadata}
  while IFS=$'\t' read -r tag id; do
    info=$(docker image inspect "$id")
    jq -e --arg id "$id" --arg arch "$RELEASE_ARCH" '
      length == 1 and .[0].Id == $id and .[0].Os == "linux" and .[0].Architecture == $arch
    ' <<< "$info" >/dev/null || die "Tested image ID/platform mismatch: $tag"
  done < <(jq -r '.[] | [.tag, .id] | @tsv' <<< "$RELEASE_IMAGES")
  while IFS=$'\t' read -r tag id; do
    staging="ci-${GITHUB_REPOSITORY##*/}-$GITHUB_RUN_ID-$GITHUB_RUN_ATTEMPT-$tag-$RELEASE_ARCH"
    ((${#staging} <= 128)) || die "Staging tag is too long"
    docker tag "$id" "$RELEASE_REPOSITORY:$staging"
    output=$(docker push "$RELEASE_REPOSITORY:$staging")
    printf '%s\n' "$output"
    # Docker's successful push response is authoritative; no remote inspection.
    digest=$(sed -nE 's/.*digest: (sha256:[0-9a-f]{64})([[:space:]].*)?$/\1/p' <<< "$output")
    [[ $digest =~ ^sha256:[0-9a-f]{64}$ ]] || die "Push did not report exactly one digest"
    mkdir -p "$directory"
    jq -n --arg repository "$RELEASE_REPOSITORY" --arg tag "$tag" --arg arch "$RELEASE_ARCH" \
      --arg id "$id" --arg digest "$digest" --arg run "$GITHUB_RUN_ID" --arg attempt "$GITHUB_RUN_ATTEMPT" \
      '{schema:1, repository:$repository, tag:$tag, arch:$arch, os:"linux",
        image_id:$id, digest:$digest, run_id:$run, run_attempt:$attempt}' \
      > "$directory/$tag-$RELEASE_ARCH.json"
  done < <(jq -r '.[] | [.tag, .id] | @tsv' <<< "$RELEASE_IMAGES")
}

validate() {
  guard
  jq -en --argjson tags "${RELEASE_TAGS:-null}" --argjson arches "${RELEASE_ARCHES:-null}" '
    ($tags | type == "array" and length > 0 and length == (unique | length) and
      all(.[]; type == "string" and test("^[A-Za-z0-9_][A-Za-z0-9_.-]{0,127}$"))) and
    ($arches | type == "array" and length > 0 and length == (unique | length) and
      all(.[]; . == "amd64" or . == "arm64"))
  ' >/dev/null || die "Invalid expected tag/architecture matrix"
  local directory=${RELEASE_METADATA:-release-metadata}
  local -a files=()
  [[ -d $directory ]] || die "Missing release metadata directory"
  mapfile -d '' -t files < <(find "$directory" -type f -name '*.json' -print0 | sort -z)
  ((${#files[@]})) || die "Missing release metadata"
  jq -es --arg repository "$RELEASE_REPOSITORY" --arg run "$GITHUB_RUN_ID" \
    --arg attempt "$GITHUB_RUN_ATTEMPT" --argjson tags "$RELEASE_TAGS" --argjson arches "$RELEASE_ARCHES" '
    . as $records |
    if all(.[];
      type == "object" and
      keys == ["arch","digest","image_id","os","repository","run_attempt","run_id","schema","tag"] and
      .schema == 1 and .repository == $repository and .run_id == $run and .run_attempt == $attempt and
      .os == "linux" and (.tag as $t | $tags | index($t) != null) and
      (.arch as $a | $arches | index($a) != null) and
      (.digest | type == "string" and test("^sha256:[0-9a-f]{64}$")) and
      (.image_id | type == "string" and test("^sha256:[0-9a-f]{64}$"))) and
      (length == ($tags | length) * ($arches | length)) and
      (length == (unique_by([.tag,.arch]) | length)) and
      (length == (unique_by(.digest) | length))
    then $records | sort_by(.tag,.arch)
    else error("Incomplete, duplicate, stale or malformed release metadata") end
  ' "${files[@]}"
}

publish() {
  local records tag arch digest suffix aliases=${RELEASE_ARCH_ALIASES:-false}
  records=$(validate) # Entire set checked before the FIRST stable-tag write.
  [[ $aliases == false || $aliases == true ]] || die "Invalid architecture alias setting"
  local -a sources=()
  while IFS= read -r tag; do
    mapfile -t sources < <(jq -r --arg tag "$tag" --arg repo "$RELEASE_REPOSITORY" \
      '.[] | select(.tag == $tag) | $repo + "@" + .digest' <<< "$records")
    docker buildx imagetools create --tag "$RELEASE_REPOSITORY:$tag" "${sources[@]}"
    # Node's existing architecture tags remain compatibility aliases, but are
    # updated only after the complete release set passes validation.
    if [[ $aliases == true ]]; then
      while IFS=$'\t' read -r arch digest; do
        case $arch in amd64) suffix=x86_64 ;; arm64) suffix=aarch64 ;; esac
        docker buildx imagetools create --tag "$RELEASE_REPOSITORY:$tag-$suffix" \
          "$RELEASE_REPOSITORY@$digest"
      done < <(jq -r --arg tag "$tag" '.[] | select(.tag == $tag) | [.arch,.digest] | @tsv' <<< "$records")
    fi
  done < <(jq -r '.[]' <<< "$RELEASE_TAGS")
}

parent() {
  local base_ref=$1 base_id=$2 derived_id=$3 base derived
  [[ $base_id =~ ^sha256:[0-9a-f]{64}$ && $derived_id =~ ^sha256:[0-9a-f]{64}$ ]] ||
    die "Invalid tested parent/derived ID"
  base=$(docker image inspect "$base_ref")
  derived=$(docker image inspect "$derived_id")
  jq -en --arg id "$base_id" --arg derived_id "$derived_id" --argjson base "$base" --argjson derived "$derived" '
    $base[0].Id == $id and $derived[0].Id == $derived_id and
    ($base[0].RootFS.Layers | length) > 0 and
    $derived[0].RootFS.Layers[0:($base[0].RootFS.Layers | length)] == $base[0].RootFS.Layers
  ' >/dev/null || die "Derived image did not retain the tested same-job parent"
}

case ${1:-} in
  guard) guard ;;
  stage) stage ;;
  validate) validate ;;
  publish) publish ;;
  parent) shift; [[ $# == 3 ]] || die "parent requires BASE_REF BASE_ID DERIVED_ID"; parent "$@" ;;
  *) die "Usage: release.sh guard|stage|validate|publish|parent" ;;
esac
