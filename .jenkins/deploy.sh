#! /bin/bash

ARG="---conventional-commits --npm-tag=$LERNA_TAG"

git config --global user.email jenkins@bearer.sh
git config --global user.name   jenkins-br

if [ ! -f ~/.npmrc ]; then
  echo "Missing .npmrc file"
  exit 1
fi

if [ ! -f ~/.gitconfig ]; then
  echo "Missing ~/.gitconfig file"
  exit 1
fi

mkdir -p ~/.ssh
echo $JENKINS_PRIVATE_KEY >> ~/.ssh/id_rsa

cat  ~/.ssh/id_rsa
git config --list

git branch
echo $ARG

echo "Starting publishing"
#yarn lerna-publish $@
