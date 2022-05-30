#!/bin/bash

create_uninstall_script()
{
    UNINSTALL_SCRIPT="$HOME/.datree/uninstall.sh"
    touch $UNINSTALL_SCRIPT && chmod +x $UNINSTALL_SCRIPT

    cat >> $UNINSTALL_SCRIPT << 'END'
    if [ "$(id -u)" -ne 0 ] ; then
    echo "This script must be executed with root privileges." && exit 1
    fi
    rm -f /usr/local/bin/datree
    rm -rf $HOME/.datree
    echo "Datree was successfully uninstalled."
END
}


set -e

osName=$(uname -s)

osArchitecture=$(uname -m)

if [[ $osArchitecture == *'aarch'* || $osArchitecture == *'arm'* ]]; then
	osArchitecture='arm64'
fi

if ! [[ -d /usr/local/bin/ ]]
then
	mkdir -p "/usr/local/bin/" 2> /dev/null || sudo mkdir -p "/usr/local/bin/"
fi

DOWNLOAD_URL=$(curl --silent "https://api.github.com/repos/datreeio/datree/releases/latest" | grep -o "browser_download_url.*\_${osName}_${osArchitecture}.zip")
DOWNLOAD_URL=${DOWNLOAD_URL//\"}
DOWNLOAD_URL=${DOWNLOAD_URL/browser_download_url: /}


OUTPUT_BASENAME=datree-latest
OUTPUT_BASENAME_WITH_POSTFIX=$OUTPUT_BASENAME.zip

echo "Installing Datree..."
echo

curl -sL $DOWNLOAD_URL -o $OUTPUT_BASENAME_WITH_POSTFIX
echo -e "\033[32m[V] Downloaded Datree\033[0m"

if ! unzip >/dev/null 2>&1;then
    echo -e "\033[31;1m error: unzip command not found \033[0m"
    echo -e "\033[33;1m install unzip command in your system \033[0m"
    exit 1
fi

unzip -qq $OUTPUT_BASENAME_WITH_POSTFIX -d $OUTPUT_BASENAME

mkdir -p ~/.datree

rm -f /usr/local/bin/datree 2> /dev/null || sudo rm -f /usr/local/bin/datree
cp $OUTPUT_BASENAME/datree /usr/local/bin 2> /dev/null || sudo cp $OUTPUT_BASENAME/datree /usr/local/bin

rm $OUTPUT_BASENAME_WITH_POSTFIX
rm -rf $OUTPUT_BASENAME

# download and save demo file
curl -s https://get.datree.io/k8s-demo.yaml > ~/.datree/k8s-demo.yaml

# create uninstall script
create_uninstall_script

echo -e "\033[32m[V] Finished Installation\033[0m"

echo

echo -e "\033[35m Usage: $ datree test ~/.datree/k8s-demo.yaml \033[0m"

echo -e "\033[35m Using Helm? => https://github.com/datreeio/helm-datree \033[0m"

echo -e "\033[35m Using Kustomize? => https://hub.datree.io/kustomize-support \033[0m"

echo -e "\033[35m Run 'datree completion -h' to learn how to generate shell autocompletions \033[0m"

echo
