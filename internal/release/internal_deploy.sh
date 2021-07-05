touch temp
git add -A
export DATREE_INTERNAL=$GIT_TAG-internal
git commit -m "release $DATREE_INTERNAL"
git tag $DATREE_INTERNAL
git push --tags
export DATREE_BUILD_VERSION=$DATREE_INTERNAL

cp ./dist/datree-macos_darwin_amd64/datree ./datree
rm -rf ./dist
mkdir -p ./dist/datree-macos_darwin_amd64
cp ./datree ./dist/datree-macos_darwin_amd64/datree

curl -sL https://git.io/goreleaser | VERSION=v$GORELEASER_VERSION bash
