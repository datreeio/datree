#!/bin/bash

ssh-keyscan github.com >> githubKey

RSA_KEY_GENERATED_FULL=$(ssh-keygen -lf githubKey)
RSA_KEY_GENERATED=${RSA_KEY_GENERATED_FULL:12:43}

RSA_KEY_GITHUB_FULL=$(curl -s https://api.github.com/meta | jq ."ssh_key_fingerprints" | grep RSA)
RSA_KEY_GITHUB=${RSA_KEY_GITHUB_FULL:17:43}

if [ $RSA_KEY_GITHUB != $RSA_KEY_GENERATED ]; then
  echo "failed to match rsa key - potential man in the middle - aborting"
  exit 1
fi

cat githubKey >> ~/.ssh/known_hosts
