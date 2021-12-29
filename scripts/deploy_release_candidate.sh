#!/bin/bash
set -ex

MAJOR_VERSION=0
MINOR_VERSION=14

latestRcTag=$(git tag --sort=-version:refname | grep "^${MAJOR_VERSION}.${MINOR_VERSION}.\d\+\-rc" | head -n 1 | grep --only-matching "^${MAJOR_VERSION}.${MINOR_VERSION}.\d\+" || true)

if [ "$latestRcTag" == "" ]; then
    nextVersion=$MAJOR_VERSION.$MINOR_VERSION.0
else
    nextVersion=$(echo $latestRcTag | awk -F. '{$NF = $NF + 1;} 1' | sed 's/ /./g')
fi

export DATREE_BUILD_VERSION=$nextVersion-rc-test-changelog
echo $DATREE_BUILD_VERSION

git tag $DATREE_BUILD_VERSION -a -m "Generated tag from TravisCI for build $TRAVIS_BUILD_NUMBER"
git push origin $DATREE_BUILD_VERSION

# bash ./scripts/sign_application.sh
# curl -sL https://git.io/goreleaser | GORELEASER_CURRENT_TAG=$DATREE_BUILD_VERSION GO_BUILD_TAG=staging VERSION=v$GORELEASER_VERSION bash
GORELEASER_CURRENT_TAG=$DATREE_BUILD_VERSION GO_BUILD_TAG=staging goreleaser --rm-dist

# bash ./scripts/brew_push_formula.sh staging $DATREE_BUILD_VERSION
