#!/bin/bash
set -ex

git tag --sort=-version:refname >tags.txt
head -n 10 tags.txt
stagingTags=$(cat tags.txt | grep '^0.13.\d')
echo $stagingTags
stagingTag=$(echo $stagingTags | head -n 1)
echo $stagingTag
tag=$(echo $stagingTag | grep --only-matching '^0.13.\d\+')
echo $tag

# latestStagingTag=$(git tag --sort=-version:refname | grep '^0.13.\d\+\-staging' | head -n 1 | grep --only-matching '^0.13.\d\+')
# if [ $TRAVIS_BRANCH == "main" ]; then
#     export SEMVER_NUMBER=$latestStagingTag
# else
#     nextVersion=$(echo $latestStagingTag | awk -F. '{$NF = $NF + 1;} 1' | sed 's/ /./g')
#     export SEMVER_NUMBER=$nextVersion
# fi
