#!/usr/bin/env bash

# BRANCH=$(git rev-parse --abbrev-ref HEAD)
# MASTER="master"
# echo $BRANCH

# if ["$MASTER" != "$BRANCH" ]; then
#   echo "You are not on master branch please switch"
#   exit 1
# fi

if [ ! -f ~/.npmrc ]; then
  echo "Missing ~/.npmrc file"
  exit 1
fi

if [ ! -f ~/.gitconfig ]; then
  echo "Missing ~/.gitconfig file"
  exit 1
fi

docker build \
  --build-arg EMAIL="$(git config user.email)" \
  --build-arg NAME="$(git config user.name)" \
  -t bearer-publish-docker \
  .

docker run -ti \
  -v $(dirname $SSH_AUTH_SOCK):$(dirname $SSH_AUTH_SOCK) \
  -e SSH_AUTH_SOCK=$SSH_AUTH_SOCK \
  -v ~/.ssh/id_rsa:/root/.ssh/id_rsa \
  -v ~/.npmrc:/root/.npmrc \
  --rm bearer-publish-docker