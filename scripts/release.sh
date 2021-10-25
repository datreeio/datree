#!/bin/bash
set -ex

latestRcTag=$(git tag --sort=-version:refname | grep "\-rc$" | head -n 1)

if [ $latestRcTag == "" ]; then
    echo "couldn't find latestRcTag"
    exit 1
fi

git checkout $latestRcTag

release_tag=${latestRcTag%-rc}
git tag $release_tag -a -m "Generated tag from manual TravisCI for build $TRAVIS_BUILD_NUMBER"
git push --tags

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

bash ./scripts/brew_push_formula.sh production $release_tag

git checkout -b "release/${release_tag}"
git push origin "release/${release_tag}"
