# Bearer

[![Build Status](https://jenkins.bearer.tech/buildStatus/icon?job=Bearer/bearer/master)](https://jenkins.bearer.tech/job/Bearer/job/bearer/job/master/)

This repository contains

- [@bearer/cli](./packages/cli)
- [@bearer/core](./packages/core)
- [create-bearer](./packages/create-bearer)
- [@bearer/functions](./packages/functions)
- [@bearer/bearer-cli](./packages/legacy-cli)
- [@bearer/logger](./packages/logger)
- [@bearer/node](./packages/node)
- [@bearer/package-init](./packages/package-init)
- [@bearer/react](./packages/react)
- [@bearer/security](./packages/security)
- [@bearer/transpiler](./packages/transpiler)
- [@bearer/tsconfig](./packages/tsconfig)
- [@bearer/tslint-config](./packages/tslint-config)
- [@bearer/types](./packages/types)
- [@bearer/ui](./packages/ui)

## How to get Started

**Install NVM**: [site link](https://github.com/creationix/nvm)

Mac Users:

```bash
brew install nvm
```

_Follow carefully post install instructions_

**Install correct node version**

```bash
nvm install
nvm use
```

---

We use [Lerna](https://github.com/lerna/lerna) to manage dependencies.

```bash
// install dependencies
yarn install


// Install dependencies and link packages together
yarn lerna bootstrap
```

_Now You should be able to go into each packages and run existing command (ex: yarn start)_

## Development

### Conventional commits

As we try to be as clear as possible on changes within bearer packages, we adopted [conventional commits](https://conventionalcommits.org/) as pattern for git commits.

_Learn more https://conventionalcommits.org/_

This repository is [commitizen](https://github.com/commitizen/cz-cli) friendly, which means you can use the below command to generate your commits. (in the root directory).

_Learn more about Commitizen and available solutions https://github.com/commitizen/cz-cli_

```bash
yarn cm
```

Using conventional commits allow us (through lerna) to generate automatically the [CHANGELOG](./CHANGELOG.md) and pick the correct version

### Local package use

Sometimes you need to test locally your changes before publishing anything. Here's some tips you can use.

**CLI**

- Let's link and build the CLI

```
cd packages/cli
yarn link
yarn build -w # starts TS compiler with watch mode enabled
```

- then within your integration

```
yarn link @bearer/cli
yarn bearer start # will use the local version you previously linked
```

At this point you are using your local cli build.

if you want to generate a new integration from you local build then you can run

```bash
/path/bearer/repo/packages/cli/lib/bin/index.js new IntegrationName
```

**Other packages**

If you want to use these packages in a local environment (in repositories not contained in this one)

_assuming we have a package named "@bearer/core"_

```bash
lerna exec yarn link --scope='@bearer/core'

// somewhere else on you computer inside a integration, for example

yarn link '@bearer/core'
```
