#!/bin/bash
set -ex

if test -z "$RELEASE_CANDIDATE_VERSION"; then
    latestRcTag=$(git tag --sort=-version:refname | grep "\-rc$" | head -n 1)
else
    latestRcTag="$RELEASE_CANDIDATE_VERSION"
fi

if test -z "$latestRcTag"; then
    echo "couldn't find latestRcTag"
    exit 1
fi

git checkout $latestRcTag

release_tag=${latestRcTag%-rc}
release_tag="$release_tag-test-yishay" # TODO: remove for creating real production release
git tag $release_tag -a -m "Generated tag from manual TravisCI for production build $TRAVIS_BUILD_NUMBER"
git push --tags

bash ./scripts/sign_application.sh

export DATREE_BUILD_VERSION=$release_tag
echo $DATREE_BUILD_VERSION

curl -sL https://git.io/goreleaser | GO_BUILD_TAG=main VERSION=v$GORELEASER_VERSION bash

# bash ./scripts/brew_push_formula.sh production $release_tag # TODO: uncomment on prod

git checkout -b "release/${release_tag}"
git push origin "release/${release_tag}"
