#!/bin/bash
set -ex

BUILD_TAG=staging

if [ "$RELEASE_DATREE_PROD" == "true" ]; then
    BUILD_TAG=main
fi

echo $BUILD_TAG
