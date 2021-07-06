set -ex
touch temp
git add -A
export DATREE_INTERNAL=0.1.$TRAVIS_BUILD_NUMBER-test-internal
git commit -m "release $DATREE_INTERNAL"
git tag $DATREE_INTERNAL -a -m "Generated internal tag from TravisCI for build $TRAVIS_BUILD_NUMBER"
git push --tags
export DATREE_BUILD_VERSION=$DATREE_INTERNAL
curl -sL https://git.io/goreleaser | VERSION=v$GORELEASER_VERSION sh -s -- --rm-dist
