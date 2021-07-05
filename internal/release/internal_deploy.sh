touch temp
git add -A
export DATREE_INTERNAL=$GIT_TAG-internal
git commit -m "release $DATREE_INTERNAL"
git tag $DATREE_INTERNAL
git push --tags
export DATREE_BUILD_VERSION=$DATREE_INTERNAL
echo GORELEASER_VERSION ======== $GORELEASER_VERSION
curl -sL https://git.io/goreleaser | VERSION=v$GORELEASER_VERSION bash
