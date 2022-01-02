#!/bin/bash
set -ex

latestRelease=$(git tag -l --sort=-v:refname | grep -v "\-rc$" | grep -v "pull"  | head -n 1)
customChangeLog=$(git log --pretty="%h %N %s%n" --decorate=full ${latestRelease//-rc}..HEAD)

echo $customChangeLog
