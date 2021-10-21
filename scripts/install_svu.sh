#!/bin/bash

set -e

SVU_TAR_FILE="svu_1.6.1_darwin_amd64.tar.gz"
CHECKSUM_AGAINST="626f7b2e13023e9b9135cf8ebd5699f9d29bd881e6bda7d9bd912ef373f6d06c"

curl -sL https://github.com/caarlos0/svu/releases/download/v1.6.1/svu_1.6.1_darwin_amd64.tar.gz -o $SVU_TAR_FILE

shasum -a 256 $SVU_TAR_FILE | grep "$CHECKSUM_AGAINST"

tar -xvf $SVU_TAR_FILE "svu"
cp ./svu /usr/local/bin

rm ./svu
rm ./$SVU_TAR_FILE
