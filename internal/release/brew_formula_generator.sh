#!/bin/bash

if [ $# -lt 2 ]
  then
    echo "Not enough arguments supplied"
    exit
fi

VERSION=$1
DESTINATION_FOLDER=$2

cat ./dist/checksums.txt

SHA256_MAC=$(cat ./dist/checksums.txt | grep Darwin_x86_64 | cut -d" " -f1)
SHA256_LINUX_INTEL=$(cat ./dist/checksums.txt | grep Linux_x86_64 | cut -d" " -f1)
SHA256_LINUX_ARM=$(cat ./dist/checksums.txt | grep Linux_arm64 | cut -d" " -f1)

cat > $DESTINATION_FOLDER/datree.rb <<-EOF
# typed: false
# frozen_string_literal: true

class Datree < Formula
  desc ""
  homepage "https://datree.io/"
  version "$VERSION"
  bottle :unneeded

  if OS.mac?
    url "https://github.com/datreeio/datree/releases/download/$VERSION/datree-cli_${VERSION}_Darwin_x86_64.zip"
    sha256 "$SHA256_MAC"
  end
  if OS.linux? && Hardware::CPU.intel?
    url "https://github.com/datreeio/datree/releases/download/$VERSION/datree-cli_${VERSION}_Linux_x86_64.zip"
    sha256 "$SHA256_LINUX_INTEL"
  end
  if OS.linux? && Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
    url "https://github.com/datreeio/datree/releases/download/$VERSION/datree-cli_${VERSION}_Linux_arm64.zip"
    sha256 "$SHA256_LINUX_ARM"
  end

  def install
    bin.install "datree"
  end
end
EOF
