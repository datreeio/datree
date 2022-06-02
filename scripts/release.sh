#!/bin/bash
set -ex

release_tag=$RELEASE_VERSION
v_release_tag=v$release_tag

git checkout "$release_tag-rc"


git tag $release_tag -a -m "Generated tag from manual TravisCI for production build $TRAVIS_BUILD_NUMBER"
git tag $v_release_tag -a -m "Generated tag with v from manual TravisCI for production build $TRAVIS_BUILD_NUMBER"

git push origin $release_tag
git push origin $v_release_tag

export DATREE_BUILD_VERSION=$release_tag
echo $DATREE_BUILD_VERSION

bash ./scripts/custom_changelog.sh

curl -sL https://git.io/goreleaser | GORELEASER_CURRENT_TAG=$DATREE_BUILD_VERSION GO_BUILD_TAG=main VERSION=v$GORELEASER_VERSION bash -s -- --rm-dist --release-notes=changelog.txt

bash ./scripts/upload_install_scripts.sh

bash ./scripts/brew_push_formula.sh production $DATREE_BUILD_VERSION

git checkout -b "release/${release_tag}"
git push origin "release/${release_tag}"
