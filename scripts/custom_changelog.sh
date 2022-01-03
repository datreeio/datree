#!/bin/bash
set -ex

latestRelease=$(curl --silent "https://api.github.com/repos/datreeio/datree/releases/latest" | jq -r '.tag_name' )
customChangeLog=$(git log --pretty="%h %N %s%n" --decorate=full --no-merges ${latestRelease//-rc}..HEAD)

echo $customChangeLog
