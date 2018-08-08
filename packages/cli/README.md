# @bearer/cli

Bearer CLI

[![oclif](https://img.shields.io/badge/cli-oclif-brightgreen.svg)](https://oclif.io)
[![Version](https://img.shields.io/npm/v/@bearer/cli.svg)](https://npmjs.org/package/@bearer/cli)

[![CircleCI](https://circleci.com/gh/Bearer/bearer/packages/cli/tree/master.svg?style=shield)](https://circleci.com/gh/Bearer/bearer/packages/cli/tree/master)

[![Appveyor CI](https://ci.appveyor.com/api/projects/status/github/Bearer/bearer/packages/cli?branch=master&svg=true)](https://ci.appveyor.com/project/Bearer/bearer/packages/cli/branch/master)
[![Codecov](https://codecov.io/gh/Bearer/bearer/packages/cli/branch/master/graph/badge.svg)](https://codecov.io/gh/Bearer/bearer/packages/cli)
[![Downloads/week](https://img.shields.io/npm/dw/@bearer/cli.svg)](https://npmjs.org/package/@bearer/cli)
[![License](https://img.shields.io/npm/l/@bearer/cli.svg)](https://github.com/Bearer/bearer/packages/cli/blob/master/package.json)

<!-- toc -->
* [@bearer/cli](#bearer-cli)
* [Usage](#usage)
* [Commands](#commands)
<!-- tocstop -->

# Usage

<!-- usage -->
```sh-session
$ npm install -g @bearer/cli
$ bearer COMMAND
running command...
$ bearer (-v|--version|version)
@bearer/cli/0.37.2 darwin-x64 node-v10.8.0
$ bearer --help [COMMAND]
USAGE
  $ bearer COMMAND
...
```
<!-- usagestop -->

# Commands

<!-- commands -->
* [`bearer deploy`](#bearer-deploy)
* [`bearer generate [NAME]`](#bearer-generate-name)
* [`bearer help [COMMAND]`](#bearer-help-command)
* [`bearer invoke INTENT_NAME`](#bearer-invoke-intent-name)
* [`bearer link SCENARIO_IDENTIFIER`](#bearer-link-scenario-identifier)
* [`bearer login`](#bearer-login)
* [`bearer new SCENARIONAME`](#bearer-new-scenarioname)
* [`bearer start`](#bearer-start)

## `bearer deploy`

Deploy a scenario

```
USAGE
  $ bearer deploy

OPTIONS
  -h, --help          show CLI help
  -i, --intents-only  Deploy intents only
  -s, --views-only    Deploy views only
```

_See code: [src/commands/deploy.ts](https://github.com/Bearer/bearer/blob/v0.37.2/src/commands/deploy.ts)_

## `bearer generate [NAME]`

Generate scenario intents or components

```
USAGE
  $ bearer generate [NAME]

OPTIONS
  -h, --help              show CLI help
  --blank-component
  --collection-component
  --root-group
  --setup
```

_See code: [src/commands/generate.ts](https://github.com/Bearer/bearer/blob/v0.37.2/src/commands/generate.ts)_

## `bearer help [COMMAND]`

display help for bearer

```
USAGE
  $ bearer help [COMMAND]

ARGUMENTS
  COMMAND  command to show help for

OPTIONS
  --all  see all commands in CLI
```

_See code: [@oclif/plugin-help](https://github.com/oclif/plugin-help/blob/v2.0.5/src/commands/help.ts)_

## `bearer invoke INTENT_NAME`

Invoke Intent locally

```
USAGE
  $ bearer invoke INTENT_NAME

OPTIONS
  -h, --help       show CLI help
  -p, --path=path
```

_See code: [src/commands/invoke.ts](https://github.com/Bearer/bearer/blob/v0.37.2/src/commands/invoke.ts)_

## `bearer link SCENARIO_IDENTIFIER`

Link your local scenario to a remote one

```
USAGE
  $ bearer link SCENARIO_IDENTIFIER

OPTIONS
  -h, --help  show CLI help
```

_See code: [src/commands/link.ts](https://github.com/Bearer/bearer/blob/v0.37.2/src/commands/link.ts)_

## `bearer login`

Login to Bearer platform

```
USAGE
  $ bearer login

OPTIONS
  -e, --email=email  (required)
  -h, --help         show CLI help
```

_See code: [src/commands/login.ts](https://github.com/Bearer/bearer/blob/v0.37.2/src/commands/login.ts)_

## `bearer new SCENARIONAME`

Generate a new scenario

```
USAGE
  $ bearer new SCENARIONAME

OPTIONS
  -h, --help  show CLI help
```

_See code: [src/commands/new.ts](https://github.com/Bearer/bearer/blob/v0.37.2/src/commands/new.ts)_

## `bearer start`

Start local development environment

```
USAGE
  $ bearer start

OPTIONS
  -h, --help    show CLI help
  --no-install
  --no-open
```

_See code: [src/commands/start.ts](https://github.com/Bearer/bearer/blob/v0.37.2/src/commands/start.ts)_
<!-- commandsstop -->
