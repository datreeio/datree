#!/bin/bash

if [ $# -lt 2 ]
  then
    echo "Not enough arguments supplied"
    exit
fi

VERSION=$1
DESTINATION_FOLDER=$2

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
    sha256 "79f1b4e98c2f9341496cbbf5557d8774dbbc8e9e96a38a6f32da039050d3679d"
  end
  if OS.linux? && Hardware::CPU.intel?
    url "https://github.com/datreeio/datree/releases/download/$VERSION/datree-cli_${VERSION}_Linux_x86_64.zip"
    sha256 "73aae467ef3b090758246881938b83bd9332276d58f1c6079838504825202710"
  end
  if OS.linux? && Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
    url "https://github.com/datreeio/datree/releases/download/$VERSION/datree-cli_${VERSION}_Linux_arm64.zip"
    sha256 "464567d0a500d1ad63662b5a1d14636aa408fb4d0ac3ad612e3110f24e5b2495"
  end

  def install
    bin.install "datree"
  end
end
EOF
