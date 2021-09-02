set -ex

export DATREE_INTERNAL=$SEMVER_NUMBER-internal

git tag $DATREE_INTERNAL -a -m "Generated internal tag from TravisCI for build $TRAVIS_BUILD_NUMBER"
git push --tags

export DATREE_BUILD_VERSION=$DATREE_INTERNAL
travis_wait curl -sL https://git.io/goreleaser | VERSION=v$GORELEASER_VERSION sh -s -- --rm-dist
bash ./internal/release/brew_push_formula.sh internal $DATREE_INTERNAL
