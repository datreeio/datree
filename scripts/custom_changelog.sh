#!/bin/bash
set -ex

latestCliReleaseFromGithub=$(curl --silent --header "Authorization: token ${GITHUB_TOKEN}" "https://api.github.com/repos/datreeio/datree/releases/latest")
# echo $latestCliReleaseFromGithub
latestRelease=$(echo $latestCliReleaseFromGithub | jq -r '.tag_name' )
git log --pretty='%h %N %s %n' --no-merges --decorate=full ${latestRelease//-rc}..HEAD > ./changelog.txt
cat changelog.txt
