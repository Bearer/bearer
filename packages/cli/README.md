# @bearer/cli

Bearer CLI

[![oclif](https://img.shields.io/badge/cli-oclif-brightgreen.svg)](https://oclif.io)
[![Version](https://img.shields.io/npm/v/@bearer/cli.svg)](https://npmjs.org/package/@bearer/cli)
![node (scoped)](https://img.shields.io/node/v/@bearer/cli.svg)
![npm bundle size (scoped)](https://img.shields.io/bundlephobia/minzip/@bearer/cli.svg)
[![Downloads/week](https://img.shields.io/npm/dw/@bearer/cli.svg)](https://npmjs.org/package/@bearer/cli)
[![License](https://img.shields.io/npm/l/@bearer/cli.svg)](https://github.com/Bearer/bearer/packages/cli/blob/master/package.json)

<!-- toc -->
* [@bearer/cli](#bearercli)
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
@bearer/cli/1.6.0 darwin-x64 node-v10.15.3
$ bearer --help [COMMAND]
USAGE
  $ bearer COMMAND
...
```
<!-- usagestop -->

# Commands

<!-- commands -->
* [`bearer autocomplete [SHELL]`](#bearer-autocomplete-shell)
* [`bearer generate`](#bearer-generate)
* [`bearer generate:component [NAME]`](#bearer-generatecomponent-name)
* [`bearer generate:function [NAME]`](#bearer-generatefunction-name)
* [`bearer help [COMMAND]`](#bearer-help-command)
* [`bearer integrations`](#bearer-integrations)
* [`bearer integrations:create`](#bearer-integrationscreate)
* [`bearer invoke FUNCTION_NAME`](#bearer-invoke-function_name)
* [`bearer link [INTEGRATION_IDENTIFIER]`](#bearer-link-integration_identifier)
* [`bearer login`](#bearer-login)
* [`bearer new [INTEGRATIONNAME]`](#bearer-new-integrationname)
* [`bearer push`](#bearer-push)
* [`bearer setup:auth [CREDENTIALS]`](#bearer-setupauth-credentials)
* [`bearer start`](#bearer-start)

## `bearer autocomplete [SHELL]`

display autocomplete installation instructions

```
USAGE
  $ bearer autocomplete [SHELL]

ARGUMENTS
  SHELL  shell type

OPTIONS
  -r, --refresh-cache  Refresh cache (ignores displaying instructions)

EXAMPLES
  $ bearer autocomplete
  $ bearer autocomplete bash
  $ bearer autocomplete zsh
  $ bearer autocomplete --refresh-cache
```

_See code: [@oclif/plugin-autocomplete](https://github.com/oclif/plugin-autocomplete/blob/v0.1.0/src/commands/autocomplete/index.ts)_

## `bearer generate`

generate function

```
USAGE
  $ bearer generate

OPTIONS
  -h, --help   show CLI help
  --path=path
  --silent

ALIASES
  $ bearer g
```

_See code: [src/commands/generate/index.ts](https://github.com/Bearer/bearer/blob/v1.6.0/src/commands/generate/index.ts)_

## `bearer generate:component [NAME]`

Generate a Bearer component

```
USAGE
  $ bearer generate:component [NAME]

OPTIONS
  -h, --help                        show CLI help
  -t, --type=blank|collection|root
  --path=path
  --silent

ALIASES
  $ bearer g:c
```

_See code: [src/commands/generate/component.ts](https://github.com/Bearer/bearer/blob/v1.6.0/src/commands/generate/component.ts)_

## `bearer generate:function [NAME]`

Generate a Bearer Function

```
USAGE
  $ bearer generate:function [NAME]

OPTIONS
  -h, --help   show CLI help
  --path=path
  --silent

ALIASES
  $ bearer g:f
```

_See code: [src/commands/generate/function.ts](https://github.com/Bearer/bearer/blob/v1.6.0/src/commands/generate/function.ts)_

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

_See code: [@oclif/plugin-help](https://github.com/oclif/plugin-help/blob/v2.1.6/src/commands/help.ts)_

## `bearer integrations`

list deployed integrations

```
USAGE
  $ bearer integrations

OPTIONS
  -h, --help   show CLI help
  --path=path
  --silent
```

_See code: [src/commands/integrations/index.ts](https://github.com/Bearer/bearer/blob/v1.6.0/src/commands/integrations/index.ts)_

## `bearer integrations:create`

create a new Integration

```
USAGE
  $ bearer integrations:create

OPTIONS
  -d, --description=description
  -h, --help                     show CLI help
  -l, --skipLink
  -n, --name=name
  --path=path
  --silent
```

_See code: [src/commands/integrations/create.ts](https://github.com/Bearer/bearer/blob/v1.6.0/src/commands/integrations/create.ts)_

## `bearer invoke FUNCTION_NAME`

invoke function locally

```
USAGE
  $ bearer invoke FUNCTION_NAME

OPTIONS
  -d, --data=data
  -f, --file=file
  -h, --help       show CLI help
  --path=path
  --silent
```

_See code: [src/commands/invoke.ts](https://github.com/Bearer/bearer/blob/v1.6.0/src/commands/invoke.ts)_

## `bearer link [INTEGRATION_IDENTIFIER]`

link to remote Bearer integration

```
USAGE
  $ bearer link [INTEGRATION_IDENTIFIER]

OPTIONS
  -h, --help   show CLI help
  --path=path
  --silent
```

_See code: [src/commands/link.ts](https://github.com/Bearer/bearer/blob/v1.6.0/src/commands/link.ts)_

## `bearer login`

login using Bearer credentials

```
USAGE
  $ bearer login

OPTIONS
  -h, --help   show CLI help
  --path=path
  --silent
```

_See code: [src/commands/login.ts](https://github.com/Bearer/bearer/blob/v1.6.0/src/commands/login.ts)_

## `bearer new [INTEGRATIONNAME]`

generate integration boilerplate

```
USAGE
  $ bearer new [INTEGRATIONNAME]

OPTIONS
  -a, --authType=OAUTH1|OAUTH2|BASIC|APIKEY|NONE|CUSTOM  Authorization type
  -d, --directory=directory                              Select a directory as source of the integration
  -f, --force                                            Force copying files
  -h, --help                                             show CLI help
  -t, --template=template                                Generate an integration from a template (git url)
  --path=path
  --silent
  --withViews                                            Experimental - generate views
```

_See code: [src/commands/new.ts](https://github.com/Bearer/bearer/blob/v1.6.0/src/commands/new.ts)_

## `bearer push`

deploy integration to Bearer

```
USAGE
  $ bearer push

OPTIONS
  -h, --help   show CLI help
  --path=path
  --silent
```

_See code: [src/commands/push.ts](https://github.com/Bearer/bearer/blob/v1.6.0/src/commands/push.ts)_

## `bearer setup:auth [CREDENTIALS]`

setup API credentials for local development.

```
USAGE
  $ bearer setup:auth [CREDENTIALS]

ARGUMENTS
  CREDENTIALS  Provide inline credentials

OPTIONS
  -h, --help   show CLI help
  --path=path
  --silent

DESCRIPTION
  If you would like to bypass the prompt, you can either:
  	* pass credentials as argument (see description later)
  	* use environment variables
  see examples

EXAMPLES
  With argument
  	bearer setup:auth CLIENT_ID:CLIENT_SECRET
  	bearer setup:auth CONSUMER_KEY:CONSUMER_SECRET
  	bearer setup:auth USERNAME:PASSWORD
  	bearer setup:auth APIKEY
  With environment variables
  	BEARER_AUTH_CLIENT_ID=CLIENT_ID BEARER_AUTH_CLIENT_SECRET=CLIENT_SECRET bearer setup:auth
  	BEARER_AUTH_CONSUMER_KEY=CONSUMER_KEY BEARER_AUTH_CONSUMER_SECRET=CONSUMER_SECRET bearer setup:auth
  	BEARER_AUTH_USERNAME=USERNAME BEARER_AUTH_CONSUMER_SECRET=PASSWORD bearer setup:auth
  	BEARER_AUTH_APIKEY=APIKEY bearer setup:auth
```

_See code: [src/commands/setup/auth.ts](https://github.com/Bearer/bearer/blob/v1.6.0/src/commands/setup/auth.ts)_

## `bearer start`

start local development environment

```
USAGE
  $ bearer start

OPTIONS
  -f, --force   Start using random available port
  -h, --help    show CLI help
  --no-install
  --no-open
```

_See code: [src/commands/start.ts](https://github.com/Bearer/bearer/blob/v1.6.0/src/commands/start.ts)_
<!-- commandsstop -->
