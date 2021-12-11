#!/bin/bash

set -ex

kool run make-docs
kool run fmt
kool run lint
kool run test

if [ ! -z "$(git status -s)" ]; then
    echo "You have uncommited changes; aborting creating release."
    exit 1
fi

read -p "What version do you want to build (0.0.0 semver format): "
if [[ ! $REPLY =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]
then
    echo "Bad version format; expected semver 0.0.0"
    exit 1
fi

export BUILD_VERSION=$REPLY

exec bash build_artifacts.sh

# TODO: create new tag / draft a new release with github CLI
