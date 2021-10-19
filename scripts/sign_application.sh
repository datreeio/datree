#!/bin/bash

# Secure private key
openssl aes-256-cbc -K $encrypted_2dfcdd1dc486_key -iv $encrypted_2dfcdd1dc486_iv -in DatreeCli.p12.enc -out DatreeCli.p12 -d

security create-keychain -p test buildagent.keychain
security default-keychain -s buildagent.keychain
security unlock-keychain -p test buildagent.keychain
security list-keychains -d user -s buildagent.keychain
security import DatreeCli.p12 -k buildagent.keychain -P $P12_PASSWORD -T /usr/bin/codesign
security set-key-partition-list -S "apple-tool:,apple:" -s -k test buildagent.keychain
security find-identity -v
