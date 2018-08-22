#! /bin/bash

ARG="${@:---conventional-commits --yes}"


git config --global user.email jenkins@bearer.sh
git config --global user.name jenkins-br

if [ ! -f ~/.npmrc ]; then
  echo "Missing ~/.npmrc file"
  exit 1
fi

if [ ! -f ~/.gitconfig ]; then
  echo "Missing ~/.gitconfig file"
  exit 1
fi

echo $JENKINS_PRIVATE_KEY >> ~/.ssh/id_rsa

cat  ~/.ssh/id_rsa
