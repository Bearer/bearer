#! /bin/bash

ARG="--yes --npm-tag=$LERNA_TAG"
FORCE_ARG="--yes --npm-tag=$LERNA_TAG --force-publish"

if [ ! .npmrc ]; then
  echo "Missing .npmrc file"
  exit 1
fi

if [ ! -f ~/.gitconfig ]; then
  echo "Missing ~/.gitconfig file"
  exit 1
fi

echo "Starting publishing"
echo "current value of force $FORCE"

git status
git --no-pager diff
git --no-pager log -1 --format="%H"

if [ $CANARY ==  "true" ]; then
  echo "releasing canary version"
  yarn lerna publish prerelease --yes --preid canary --force-publish --canary
elif [  $FORCE ==  "true"]; then
  echo "force publishing.."
  yarn lerna-publish $FORCE_ARG
else
  yarn lerna-publish $ARG
fi
