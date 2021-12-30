#!/bin/bash
set -ex

latestRelease=$(git tag --sort=-version:refname | grep "\-rc$" | head -n 1)
customChangeLog=$(git log --pretty="%h %N %s" --decorate=full ${latestRelease//-rc}..HEAD)

echo $customChangeLog

