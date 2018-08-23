#! /bin/bash

echo "init dependencies"
echo "//registry.npmjs.org/:_authToken=$NPM_TOKEN" > ~/.npmrc
yarn install --frozen-lockfile
yarn run lerna bootstrap --froken-lockfile