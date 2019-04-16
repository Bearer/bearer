# Bearer - The API Integration Framework

<p align="center">
  <a href="https://www.bearer.sh">
    <img alt="Bearer Documentation" src="https://bearer-hub-staging.netlify.com/static/bearer-api-integration-0fc65950490996a35829334804035862.jpg" width="500">
  </a>

  <p align="center">

Bearer provides all of the tools to build, run and manage API
<br/>
<a href="https://www.bearer.sh/?utm_source=github&utm_campaign=repository">Learn more</a>

  </p>
</p>

---

[![Version][version-svg]][package-url]
[![License][license-image]][license-url]
[![Build Status][ci-svg]][ci-url]

<details>
  <summary><strong>Table of contents</strong></summary>

- [Why](#why)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Documentation](#documentation)
- [Command References](#command-references)
- [Contributing](#contributing)
- [License](#license)
  </details>

## Why

You should use Bearer if you want to:

- Consume any API in minutes
- Map API endpoints to your app model
- Integrate into your code with one line
- Deploy and Scale without fuss
- Monitor every API call
- Manage your integrations

## Installation

```bash
# Using yarn
yarn create bearer helloWorld

# Using npm
npm init bearer helloWorld

# Using npx
npx create-bearer helloWorld
```

## Quick Start

Follow the [quick start](http://docs.bearer.sh/building-integration/quick-start) to create your first integration.

## Documentation

The documentation is available on the Bearer [doc center](http://docs.bearer.sh).

## Command References

The Framework comes with a Command Line Interface (CLI) providing many commands:

```bash
$ yarn bearer -h

Bearer CLI

VERSION
  @bearer/cli/rc-1

USAGE
  $ bearer [COMMAND]

COMMANDS
  autocomplete  display autocomplete installation instructions
  generate      generate function
  help          display help for bearer
  integrations  list deployed integrations
  invoke        invoke function locally
  link          link to remote Bearer integration
  login         login using Bearer credentials
  new           generate integration boilerplate
  push          deploy integration to Bearer
  setup         setup API credentials for local development
  start         start local development environment
```

Find below some explanation of commands by order of importance.

### New

As an alternative to the `yarn create bearer` command, you can also use the new command:

```bash
yarn bearer new
```

### Setup

This command help you setup your \(test\) API authentication credentials in local development:

```bash
yarn bearer setup:auth
```

### Generate

The `generate` command provides a fast and easy way to bootstrap a **function**, taking care of much of the boilerplate code:

```text
yarn bearer generate:function FunctionName
```

### Invoke

The `invoke` command let you run a function directly from the console:

```bash
yarn bearer invoke FunctionName
```

It's a great way to test and debug a function.

You can also pass arguments, as a JSON file, to the `invoke` command, which are translated as params to the function:

```bash
yarn bearer invoke FunctionName -p tests/function.json
```

```json
{ "params": { "fullName": "Bearer/bearer" } }
```

_NB: The file path is relative to the integration folder._

### Start

The `start` command provides a local Web Server to test the integration:

```bash
yarn bearer start
```

Behind the scene, it builds the integration (transpiling code, etc...) and manage things like live reload.

### Push

The `push` command let you deploy your integration to the platform 🚀

```bash
yarn bearer push
```

### Login

The login command let you connect to the platform from the CLI:

```bash
yarn bearer login
```

It's a required step in order to `push` your integration.

### Link

The `link` command let you link an integration registered on your bearer's dashboard with your current integration code:

```bash
yarn bearer link
```

_NB: This command is triggered automatically when you deploy an integration for the first time._

### Integrations

The `integrations` command let list your deployed integration and also create a new one:

```bash
yarn bearer integrations
yarn bearer integrations:create
```

## Contributing

We welcome all contributors, from casual to regular 💙

- **Bug report**. Is something not working as expected? [Send a bug report](https://github.com/bearer/bearer/issues/new).
- **Feature request**. Would you like to add something to the framework? [Send a feature request](https://github.com/bearer/bearer/issues/new).
- **Documentation**. Did you find a typo in the doc? [Open an issue](https://github.com/bearer/bearer/issues/new) and we'll take care of it.
- **Development**. If you don't know where to start, you can check the [open issues](https://github.com/bearer/bearer/issues?q=is%3Aissue+is%3Aopen).

To start contributing to code, you need to:

1. [Fork the project](https://help.github.com/articles/fork-a-repo/)
2. [Clone the repository](https://help.github.com/articles/cloning-a-repository/)
3. Checkout the [Getting Started](GETTING_STARTED.md)

## License

Bearer is [MIT licensed][license-url].

<!-- Badges -->

[version-svg]: https://img.shields.io/npm/v/@bearer/react.svg?style=flat-square
[package-url]: https://npmjs.org/package/@bearer/cli
[license-image]: http://img.shields.io/badge/license-MIT-green.svg?style=flat-square
[ci-svg]: https://jenkins.bearer.tech/buildStatus/icon?job=Bearer%2Fbearer%2Fmaster
[ci-url]: https://jenkins.bearer.tech/job/Bearer/job/bearer/job/master/
[license-url]: LICENSE

<!-- Links -->

[bearer-website]: https://www.bearer.sh/?utm_source=github&utm_campaign=repository
