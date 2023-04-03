#!/bin/bash
set -ex

git checkout main
git pull

latestCliReleaseFromGithub=$(curl --silent --header "Authorization: token ${GITHUB_TOKEN}" "https://api.github.com/repos/datreeio/datree/releases/latest")
# echo $latestCliReleaseFromGithub
latestRelease=$(echo $latestCliReleaseFromGithub | jq -r '.tag_name' )
git log --pretty='%h %N %s %n' --no-merges --decorate=full 1.8.47..HEAD > ./changelog.txt
git checkout - 