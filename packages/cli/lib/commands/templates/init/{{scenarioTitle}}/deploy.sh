#!/usr/bin/bash

set -e
npm install -g yarn
yarn global add @bearer/bearer-cli

cd /scenario/intents && yarn
cd /scenario/screens && yarn

BEARER_ENV=dev bearer deploy
