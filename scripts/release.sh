#!/bin/bash
set -ex

release_tag=$RELEASE_VERSION

git checkout "$release_tag-rc"

git tag $release_tag -a -m "Generated tag from manual TravisCI for production build $TRAVIS_BUILD_NUMBER"
git push origin $release_tag # TODO: check if goreleaser pushes the tag itself (so no need to push here)

bash ./scripts/sign_application.sh

export DATREE_BUILD_VERSION=$release_tag
echo $DATREE_BUILD_VERSION

curl -sL https://git.io/goreleaser | GORELEASER_CURRENT_TAG=$DATREE_BUILD_VERSION GO_BUILD_TAG=main VERSION=v$GORELEASER_VERSION bash

bash ./scripts/brew_push_formula.sh production $DATREE_BUILD_VERSION

git checkout -b "release/${release_tag}"
git push origin "release/${release_tag}"
