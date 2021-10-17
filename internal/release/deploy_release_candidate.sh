#!/bin/bash
set -ex

SEMVER_NUMBER=$(bash ./internal/release/define_semver_number.sh)
export DATREE_BUILD_VERSION=$SEMVER_NUMBER-rc

export GIT_TAG=$DATREE_BUILD_VERSION
git tag $GIT_TAG -a -m "Generated tag from TravisCI for build $TRAVIS_BUILD_NUMBER"
git push --tags

bash ./internal/release/sign_application.sh
curl -sL https://git.io/goreleaser | VERSION=v$GORELEASER_VERSION bash

bash ./internal/release/brew_push_formula.sh staging $DATREE_BUILD_VERSION
