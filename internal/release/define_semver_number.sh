#!/bin/bash
set -ex

latestStagingTag=$(git tag --sort=-version:refname | grep '^0.13.\d\+\-staging' | head -n 1 | grep --only-matching '^0.13.\d\+')
if [ $TRAVIS_BRANCH == "main" ]; then
    echo "$latestStagingTag"
else
    nextVersion=$(echo $latestStagingTag | awk -F. '{$NF = $NF + 1;} 1' | sed 's/ /./g')
    echo "$nextVersion"
fi
