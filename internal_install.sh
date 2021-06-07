
#!/bin/bash

osName=$(uname -s)
DOWNLOAD_URL=$(curl --silent "https://api.github.com/repos/datreeio/datree/releases" | grep -o -m 1 "browser_download_url.*internal_${osName}_x86_64.zip")

DOWNLOAD_URL=${DOWNLOAD_URL//\"}
DOWNLOAD_URL=${DOWNLOAD_URL/browser_download_url: /}


OUTPUT_BASENAME=datree-latest
OUTPUT_BASENAME_WITH_POSTFIX=$OUTPUT_BASENAME.zip

curl -L $DOWNLOAD_URL -o $OUTPUT_BASENAME_WITH_POSTFIX
unzip $OUTPUT_BASENAME_WITH_POSTFIX -d $OUTPUT_BASENAME

DATREE_CONFIG_PATH=~/.datree
mkdir -p $DATREE_CONFIG_PATH

if [[ $osName == "Linux" ]];
then
    sudo rm -f /usr/local/bin/datree
    sudo cp $OUTPUT_BASENAME/datree /usr/local/bin
else
    rm -f /usr/local/bin/datree
    cp $OUTPUT_BASENAME/datree /usr/local/bin
fi

CONFIG_FILE_PATH=$DATREE_CONFIG_PATH/config.yaml
if [ ! -f "$CONFIG_FILE_PATH" ]; then
    echo "token: internal_"$(openssl rand -hex 12) >> $CONFIG_FILE_PATH
fi

rm $OUTPUT_BASENAME_WITH_POSTFIX
rm -rf $OUTPUT_BASENAME

curl https://get.datree.io/k8s-demo.yaml > ~/.datree/k8s-demo.yaml

