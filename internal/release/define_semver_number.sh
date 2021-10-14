#!/bin/bash
set -ex

MINOR_VERSION=14
MAJOR_VERSION=0
latestStagingTag=$(git tag --sort=-version:refname | grep "^${MAJOR_VERSION}.${MINOR_VERSION}.\d\+\-staging" | head -n 1 | grep --only-matching "^${MAJOR_VERSION}.${MINOR_VERSION}.\d\+" || true)

if [ $TRAVIS_BRANCH == "main" ]; then
    if [ "$latestStagingTag" == "" ]; then
        echo "latestStagingTag must be a legit tag version. Failing.."
        exit 1
    fi
    echo "$latestStagingTag"
else
    if [ "$latestStagingTag" == "" ]; then
        nextVersion=$MAJOR_VERSION.$MINOR_VERSION.0
    else
        nextVersion=$(echo $latestStagingTag | awk -F. '{$NF = $NF + 1;} 1' | sed 's/ /./g')
    fi
    echo "$nextVersion"
fi
