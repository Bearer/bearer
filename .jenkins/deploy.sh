#! /bin/bash

ARG="---conventional-commits --npm-tag=$LERNA_TAG"

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

mkdir -p ~/.ssh
echo $JENKINS_PRIVATE_KEY >> ~/.ssh/id_rsa

cat  ~/.ssh/id_rsa
git config --list

git branch
echo $ARG

echo "Fetch all remote tags"
git tag --list | grep 0.5
echo "Starting publishing"
#yarn lerna-publish-cicd $ARG
