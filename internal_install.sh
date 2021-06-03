
#!/bin/bash

osName=$(uname -s)
DOWNLOAD_URL=$(curl --silent "https://api.github.com/repos/datreeio/datree/releases" | grep -o -m 1 "browser_download_url.*internal_${osName}_x86_64.zip")
echo $DOWNLOAD_URL
exit 0

DOWNLOAD_URL=${DOWNLOAD_URL//\"}
DOWNLOAD_URL=${DOWNLOAD_URL/browser_download_url: /}

OUTPUT_BASENAME=datree-latest
OUTPUT_BASENAME_WITH_POSTFIX=$OUTPUT_BASENAME.zip

curl -L $DOWNLOAD_URL -o $OUTPUT_BASENAME_WITH_POSTFIX
unzip $OUTPUT_BASENAME_WITH_POSTFIX -d $OUTPUT_BASENAME

mkdir -p ~/.datree

if [[ $osName == "Linux" ]]; then

    CONFIG_FILE=/usr/local/bin/datree/config.yaml
    if ! [test -f "$CONFIG_FILE"]; then
        sudo rm -f `ls | grep -v "config.yaml"`
    else 
        sudo rm -f /usr/local/bin/datree
    fi

    # copy but without overwrite
    sudo cp -n $OUTPUT_BASENAME/datree /usr/local/bin
else
    rm -f /usr/local/bin/datree
    cp $OUTPUT_BASENAME/datree /usr/local/bin
fi

rm $OUTPUT_BASENAME_WITH_POSTFIX
rm -rf $OUTPUT_BASENAME

curl https://get.datree.io/k8s-demo.yaml > ~/.datree/k8s-demo.yaml

