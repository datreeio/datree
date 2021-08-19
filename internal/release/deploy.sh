set -ex
if [ $TRAVIS_BRANCH == "main" ]; then 
  export DATREE_BUILD_VERSION=$SEMVER_NUMBER; 
  export DATREE_BREW_REPO_NAME=homebrew-datree
else 
  export DATREE_BUILD_VERSION=$SEMVER_NUMBER-$TRAVIS_BRANCH; 
  export DATREE_BREW_REPO_NAME=homebrew-datree-staging
fi
sed -ie "s/___TAP_NAME/$DATREE_BREW_REPO_NAME/" .goreleaser.yml
cat .goreleaser.yml


export GIT_TAG=$DATREE_BUILD_VERSION
git tag $GIT_TAG -a -m "Generated tag from TravisCI for build $TRAVIS_BUILD_NUMBER"
git push --tags
git stash save --keep-index --include-untracked

# Secure private key
openssl aes-256-cbc -K $encrypted_2dfcdd1dc486_key -iv $encrypted_2dfcdd1dc486_iv -in DatreeCli.p12.enc -out DatreeCli.p12 -d

security create-keychain -p test buildagent.keychain
security default-keychain -s buildagent.keychain
security unlock-keychain -p test buildagent.keychain
security list-keychains -d user -s buildagent.keychain
security import DatreeCli.p12 -k buildagent.keychain -P $P12_PASSWORD -T /usr/bin/codesign
security set-key-partition-list -S "apple-tool:,apple:" -s -k test buildagent.keychain
security find-identity -v
curl -sL https://git.io/goreleaser | VERSION=v$GORELEASER_VERSION bash
