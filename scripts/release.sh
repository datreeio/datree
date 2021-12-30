#!/bin/bash
set -ex

release_tag=$(bash scripts/define_release_version.sh 0.14.143-rc-test-changelog)
# if [ -z "$release_tag"]
# then
#     release_tag=$(git tag --sort=-version:refname | grep "\-rc$" | head -n 1)
# else
#     echo "\$release_tag"
# fi

# git checkout "$release_tag"

git tag $release_tag -a -m "Generated tag from manual TravisCI for production build $TRAVIS_BUILD_NUMBER"
git push origin $release_tag # TODO: check if goreleaser pushes the tag itself (so no need to push here)

# bash ./scripts/sign_application.sh

export DATREE_BUILD_VERSION=$release_tag
echo $DATREE_BUILD_VERSION

# --skip-publish --rm-dist --release-notes ./scripts/custom_changelog.sh
#

curl -sL https://git.io/goreleaser | GORELEASER_CURRENT_TAG=$DATREE_BUILD_VERSION GO_BUILD_TAG=main VERSION=v$GORELEASER_VERSION bash -- --skip-publish --rm-dist --release-notes ./scripts/custom_changelog.sh
# bash ./scripts/upload_install_scripts.sh

# bash ./scripts/brew_push_formula.sh production $DATREE_BUILD_VERSION

# git checkout -b "release/${release_tag}"
# git push origin "release/${release_tag}"
