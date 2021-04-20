#!/bin/bash

osName=$(uname -s)
DOWNLOAD_URL=$(curl --silent "https://api.github.com/repos/datreeio/datree/releases/latest" | grep -o "browser_download_url.*\_${osName}_x86_64.zip")

DOWNLOAD_URL=${DOWNLOAD_URL//\"}
DOWNLOAD_URL=${DOWNLOAD_URL/browser_download_url: /}

OUTPUT_BASENAME=datree-latest
OUTPUT_BASENAME_WITH_POSTFIX=$OUTPUT_BASENAME.zip

curl -L $DOWNLOAD_URL -o $OUTPUT_BASENAME_WITH_POSTFIX
unzip $OUTPUT_BASENAME_WITH_POSTFIX -d $OUTPUT_BASENAME

mkdir -p ~/.datree

if [[ $osName == "Linux" ]]
then
    sudo cp $OUTPUT_BASENAME/datree /usr/local/bin
else
    cp $OUTPUT_BASENAME/datree /usr/local/bin
fi

rm $OUTPUT_BASENAME_WITH_POSTFIX
rm -rf $OUTPUT_BASENAME

curl https://get.datree.io/k8s-demo.yaml > ~/.datree/k8s-demo.yaml
