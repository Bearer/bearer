#! /bin/bash

echo "//registry.npmjs.org/:_authToken=$NPM_TOKEN" > ~/repo/.npmrc
yarn install --frozen-lockfile
yarn run lerna bootstrap -- --froken-lockfile