# Bearer

[![CircleCI](https://circleci.com/gh/BearerSH/bearer.svg?style=svg&circle-token=a18705e56c8c60cb7a98cc08e1af5be001e4239a)](https://circleci.com/gh/BearerSH/bearer)

This repository contains

- [Bearer core](./packages/core)
- [Bearer CLI](./packages/cli)
- [Bearer ui](./packages/ui)

## How to get Started

```bash
// install dependencies
yarn install

// install lerna
yarn global add lerna

// Install dependencies and link packages together
lerna bootstrap
```

_Now You should be able to go into each packages and run existing command (ex: yarn start)_

We use [Lerna](https://github.com/lerna/lerna) to manage dependencies.

### ⚠️ Publish packages ️⚠️

```bash
yarn lerna-publish
```

For further information on Lerna please read Lerna's documentation: [Lerna](https://github.com/lerna/lerna)

### Development

If you want to use these packages locally (in repositories not contained in this one)

**Link packages**

_assuming we have a package named "@bearer/core"_

```bash
lerna run yarn link --scope='@bearer/core'

// somewhere else on you computer inside a scenario, for example

yarn link '@bearer/core'
```
