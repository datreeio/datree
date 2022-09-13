#!/bin/bash
set -ex

release_tag=$RELEASE_VERSION
v_release_tag=v$release_tag

git checkout "$release_tag-rc"


git tag $release_tag -a -m "Generated tag from manual GH action production build $GITHUB_ACTION_RUN_ID"
git tag $v_release_tag -a -m "Generated tag with v from manual GH action for production build $GITHUB_ACTION_RUN_ID"

git push origin $release_tag
git push origin $v_release_tag

export DATREE_BUILD_VERSION=$release_tag
echo $DATREE_BUILD_VERSION

bash ./scripts/custom_changelog.sh

#maybe will be needed as in deploy rc
#git restore ./scripts/release.sh
curl -sL https://git.io/goreleaser | GORELEASER_CURRENT_TAG=$DATREE_BUILD_VERSION GO_BUILD_TAG=main VERSION=v$GORELEASER_VERSION bash -s -- --rm-dist --release-notes=changelog.txt

bash ./scripts/upload_install_scripts.sh

git checkout -b "release/${release_tag}"
git push origin "release/${release_tag}"
