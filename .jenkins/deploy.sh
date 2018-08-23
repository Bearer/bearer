#! /bin/bash

ARG="---conventional-commits --npm-tag=$LERNA_TAG"

mkdir -p ~/.ssh
cat $JENKINS_PRIVATE_KEY >> ~/.ssh/id_rsa
chmod 600 ~/.ssh/id_rsa
ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

git config --global user.email jenkins@bearer.sh
git config --global user.name   jenkins-br

git_url=$(git config --get remote.origin.url | sed "s/https:\/\/github\.com\//git@github\.com:/")
git remote set-url origin $git_url

echo $git_url

if [ ! -f ~/.npmrc ]; then
  echo "Missing .npmrc file"
  exit 1
fi

if [ ! -f ~/.gitconfig ]; then
  echo "Missing ~/.gitconfig file"
  exit 1
fi

echo "Starting publishing"
yarn lerna-publish-cicd $ARG
