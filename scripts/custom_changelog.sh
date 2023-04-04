#!/bin/bash
set -ex

git rev-parse --abbrev-ref HEAD # Show current branch
git status
git log -n 5 --format="%h" > latest5commits.txt 
cat latest5commits.txt 
git log -n 5 > latest5commits.txt 
cat latest5commits.txt 
git checkout main
git pull --unshallow
git status
git log -n 5 --format="%h" > latest5commits.txt
cat latest5commits.txt

latestCliReleaseFromGithub=$(curl --silent --header "Authorization: token ${GITHUB_TOKEN}" "https://api.github.com/repos/datreeio/datree/releases/latest")
# echo $latestCliReleaseFromGithub
latestRelease=$(echo $latestCliReleaseFromGithub | jq -r '.tag_name' )
git log --pretty='%h %N %s %n' --no-merges --decorate=full ${latestRelease//-rc}..HEAD > ./changelog.txt
cat changelog.txt
git checkout - 
