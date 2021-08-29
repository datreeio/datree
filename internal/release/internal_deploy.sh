set -ex

export DATREE_INTERNAL=$SEMVER_NUMBER-internal
sed -ie "s/___TAP_NAME/$DATREE_INTERNAL/" .goreleaser.yml
git add .goreleaser.yml
git commit -m "release $DATREE_INTERNAL"
git stash save --keep-index --include-untracked

git tag $DATREE_INTERNAL -a -m "Generated internal tag from TravisCI for build $TRAVIS_BUILD_NUMBER"
git push --tags

export DATREE_BUILD_VERSION=$DATREE_INTERNAL
curl -sL https://git.io/goreleaser | VERSION=v$GORELEASER_VERSION sh -s -- --rm-dist
