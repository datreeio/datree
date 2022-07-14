#!/bin/bash
set -ex

MAJOR_VERSION=1
MINOR_VERSION=5

latestRcTag=$(git tag --sort=-version:refname | grep -E "^${MAJOR_VERSION}\.${MINOR_VERSION}.[0-9]+-rc" | head -n 1 | grep --only-matching "^${MAJOR_VERSION}\.${MINOR_VERSION}.[0-9]\+" || true)

if [ "$latestRcTag" == "" ]; then
    nextVersion=$MAJOR_VERSION.$MINOR_VERSION.0
else
    nextVersion=$(echo $latestRcTag | awk -F. '{$NF = $NF + 1;} 1' | sed 's/ /./g')
fi

export DATREE_BUILD_VERSION=$nextVersion-rc
echo $DATREE_BUILD_VERSION

v_release_tag=v$DATREE_BUILD_VERSION

git tag $DATREE_BUILD_VERSION -a -m "Generated tag from TravisCI for build $TRAVIS_BUILD_NUMBER"
git tag $v_release_tag -a -m "Generated tag with v from TravisCI for build $TRAVIS_BUILD_NUMBER"
git push --tags

curl -sL https://git.io/goreleaser | GORELEASER_CURRENT_TAG=$DATREE_BUILD_VERSION GO_BUILD_TAG=staging VERSION=v$GORELEASER_VERSION bash

bash ./scripts/brew_push_formula.sh staging $DATREE_BUILD_VERSION
