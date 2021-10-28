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
echo $latestRcTag

git checkout $latestRcTag

release_tag=${latestRcTag%-rc}
git tag $release_tag -a -m "Generated tag from manual TravisCI for production build $TRAVIS_BUILD_NUMBER"
git push origin $release_tag # TODO: check if goreleaser pushes the tag itself (so no need to push here)

bash ./scripts/sign_application.sh

export DATREE_BUILD_VERSION=$release_tag
echo $DATREE_BUILD_VERSION

curl -sL https://git.io/goreleaser | GORELEASER_CURRENT_TAG=$DATREE_BUILD_VERSION GO_BUILD_TAG=main VERSION=v$GORELEASER_VERSION bash

bash ./scripts/brew_push_formula.sh production $DATREE_BUILD_VERSION

git checkout -b "release/${release_tag}"
git push origin "release/${release_tag}"
