set -ex

export DATREE_INTERNAL=$SEMVER_NUMBER-internal
sed -ie "s/homebrew-datree/homebrew-datree-internal/" .goreleaser.yml
git add -A
git commit -m "release $DATREE_INTERNAL"
git tag $DATREE_INTERNAL -a -m "Generaed internal tag from TravisCI for build $TRAVIS_BUILD_NUMBER"
git push --tags

export DATREE_BUILD_VERSION=$DATREE_INTERNAL
curl -sL https://git.io/goreleaser | VERSION=v$GORELEASER_VERSION sh -s -- --rm-dist
