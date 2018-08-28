# Bearer

[![Build Status](https://jenkins.bearer.tech/buildStatus/icon?job=Bearer/bearer/master)](https://jenkins.bearer.tech/job/Bearer/job/bearer/job/master/)

This repository contains

- [Bearer CLI](./packages/cli)
- [Bearer core](./packages/core)
- [Bearer intents](./packages/intents)
- [Bearer templates](./packages/templates)
- [Bearer transpiler](./packages/transpiler)
- [Bearer ui](./packages/ui)

## How to get Started

---

We use [Lerna](https://github.com/lerna/lerna) to manage dependencies.

```bash
// install dependencies
yarn install

// install lerna
yarn global add lerna

// Install dependencies and link packages together
lerna bootstrap
```

_Now You should be able to go into each packages and run existing command (ex: yarn start)_

## Development

---

### ⚠️ Publish packages ️⚠️

```bash
./scripts/release.sh

// or pass lerna arguments
./scripts/release.sh --force-publish
```

For further information on Lerna please read Lerna's documentation: [Lerna](https://github.com/lerna/lerna)

### Upgrade stencil version dependency

```bash
yarn upgrade-stencil
```

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

Because npm/yarn link behaviours can be really annoying, if you switch multiple time a day between local and distribution version for instance, we recommand this process.

- Start local watcher (automatically rebuild) CLI (because rewritten with TS)

```bash
cd packages/cli
yarn build -w # starts TS compiler with watch mode enabled
```

at this point you can use the CLI as follow

```bash
/path/bearer/repo/packages/cli/lib/bin/index.js start # or whatever bearer command you want
```

Additional tip

```bash
echo 'alias bl="/path/bearer/bearer/packages/cli/lib/bin/index.js"' >> ~/.zshrc # or .bashrc or everything else you use

# within a scenario
bl start
```

**Other packages**

If you want to use these packages in a local environment (in repositories not contained in this one)

_assuming we have a package named "@bearer/core"_

```bash
lerna run yarn link --scope='@bearer/core'

// somewhere else on you computer inside a scenario, for example

yarn link '@bearer/core'
```
