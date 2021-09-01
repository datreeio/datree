#!/bin/bash

if [ $# -lt 2 ]
  then
    echo "Not enough arguments supplied"
    exit
fi

DEPLOYMENT=$1
VERSION=$2
BREW_REPO_NAME=""

if [ $DEPLOYMENT == "production" ]; then
  BREW_REPO_NAME="homebrew-datree"
elif [ $DEPLOYMENT == "internal" ]; then
  BREW_REPO_NAME="homebrew-datree-internal"
elif [ $DEPLOYMENT == "staging" ]; then 
  BREW_REPO_NAME="homebrew-datree-staging"
else
  echo "No such deployment $DEPLOYMENT: Skipping deployment to brew"
  exit
fi

BREW_REPO_URL="git@github.com:datreeio/${BREW_REPO_NAME}.git"

git clone $BREW_REPO_URL
bash ./internal/release/brew_formula_generator.sh $VERSION $BREW_REPO_NAME
cd $BREW_REPO_NAME
git add -A
git commit -m "Brew formula update for datree-cli version $VERSION"
git push
