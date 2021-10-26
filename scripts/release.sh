#!/bin/bash
set -ex

latestRcTag=$(git tag --sort=-version:refname | grep "\-rc$" | head -n 1)

if [ $latestRcTag == "" ]; then
    echo "couldn't find latestRcTag"
    exit 1
fi

git checkout $latestRcTag

release_tag=${latestRcTag%-rc}
release_tag="$release_tag-test-yishay" # TODO: remove for creating real production release
git tag $release_tag -a -m "Generated tag from manual TravisCI for production build $TRAVIS_BUILD_NUMBER"
git push --tags

bash ./scripts/sign_application.sh

curl -sL https://git.io/goreleaser | VERSION=v$GORELEASER_VERSION bash

bash ./scripts/brew_push_formula.sh production $release_tag

git checkout -b "release/${release_tag}"
git push origin "release/${release_tag}"
