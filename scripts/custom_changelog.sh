#!/bin/bash
set -ex

latestRelease=$(curl --silent "https://api.github.com/repos/datreeio/datree/releases/latest" | jq -r '.tag_name' )
customChangeLog=$(git log --pretty="%h %N %s%n" --no-merges --decorate=full ${latestRelease//-rc}..HEAD > ./dist/changelog.txt)

