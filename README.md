<div align="center">

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./docs/assets/img/curio-logo-dark.svg">
  <img alt="Curio" src="./docs/assets/img/curio-logo-light.svg">
</picture>

<br />

[![GitHub Release][release-img]][release]
[![Test][test-img]][test]
[![GitHub All Releases][github-all-releases-img]][release]
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)
[![Discord](https://img.shields.io/discord/1042147477765242973?label=discord)][discord]

</div>

# Curio: The data-first security scanner that finds risks and vulnerabilities in your code to protect sensitive data

Curio is a static code analysis tool (SAST) that scans your source code to discover security risks and vulnerabilities that put your sensitive data at risk (PHI, PD, PII, Financial data).

[Explore the docs](https://curio.sh) - [Getting Started](#getting-started) - [FAQ](#faq) - [Report Bug](https://github.com/Bearer/curio/issues/new/choose) - [Discord Community][discord]

```
Scanning target... 100% [===============] (7481/7481, 22 files/s) [5m44s]

Policy Report

=====================================
Policy list:

- Logger leaking
- Application level encryption missing

=====================================

SUCCESS

2 policies were run and no breaches were detected.

```

Curio is a static code analysis tool (SAST) dedicated to data security. It scans your source code to discover sensitive data flows (PHI, PII, PD as well as Data stores, internal and external APIs) and data security risks (leaks, missing encryption, third-party sharing, etc).
You can use Curio as a CLI or add it to your CI/CD pipeline.

Curio helps developers and security teams to:

- Protect their application from leaking sensitive data (*loggers, cookies, third-parties, etc*.)
- Protect their application from having their sensitive data breached (*missing encryption, insecure communication, SQL injection, etc.*)
- Monitor sensitive data flows across every component (*Data stores, internal and external APIs*)

Curio is Open Source ([*see license*](#license)), and is built to be fully customizable, from creating your own policies, to adding custom code detectors up to enriching our data classifiers.

Curio also powers [Bearer](https://www.bearer.com), the developer-first platform to proactively find and fix data security risks and vulnerabilities across your development lifecycle

## Getting started

Scan your first project in X minutes or less.

### Installation

Curio is available as a standalone executable binary. The latest release is available on the [releases tab](https://github.com/Bearer/curio/releases/latest), or use one of the methods below.

#### Install Script

:warning: **Not working till public** :warning:

This script downloads the Curio binary automatically based on your OS and architecture.

```bash
curl -sfL https://raw.githubusercontent.com/Bearer/curio/main/contrib/install.sh | sh
```

_Defaults to `./bin` as a bin directory and to the latest releases_

If you need to customize the options, use the following to pass parameters:

```bash
curl -sfL https://raw.githubusercontent.com/Bearer/curio/main/contrib/install.sh | sh -s -- -b /usr/local/bin
```

#### Binary

Download the archive file for your operating system/architecture from [here](https://github.com/Bearer/curio/releases/latest/). Unpack the archive, and put the binary somewhere in your $PATH (on UNIX-y systems, /usr/local/bin or the like). Make sure it has permission to execute:

```bash
chmod +x ./curio
```

### Scan your project

Run `curio scan` on a project directory:

```bash
curio scan /path/to/your_project
```

or a single a file:

```bash
curio scan ./curio-ci-test/Pipfile.lock
```

<!-- TODO: insert sample output or video here -->

Additional options for using and configuring the `scan` command can be found in the [scan documentation](https://curio.sh/reference/commands/#scan).

## FAQs

### What can I do with Curio?

The number one thing you can do with Curio is protect sensitive data processed by your applications from being exposed by identifying risks. Add Curio to your CI/CD and catch data breach and data leak risks early. This is the **Policy feature.**

Curio also helps your Compliance & Privacy team assess sensitive data flows processed through your applications, allowing them to comply with the growing number of data regulations across the World (*GDPR, CCPA, CPRA, LGPD, PIPA, HIPAA, etc.*). This is the **Data Report feature**.

### Supported Language

Curio currently fully supports Ruby, while JavaScript/TypeScript support is coming soon. Additional languages listed below are supported only for the Data Report features for now.

| Languages | Data Report | Policy |
| --- | --- | --- |
| Ruby | ✅ | ✅ |
| JavaScript / TS | ✅ | Coming soon |
| PHP | ✅ | Coming later |
| Java | ✅ | Coming later |
| C# | ✅ | Coming later |
| Go | ✅ | Coming later |
| Python | ✅ | Coming later |

### What is sensitive data?

Out of the box, Curio supports 120+ data types, ranging from Personal Data (PD), to Personal Identifiable Information (PII), Protected Health Information (PHI), and Financial Data. Please refer to the data types documentation for the complete list of supported data types and how the classification works.

You can also add new data types, change the overall taxonomy, or update the classifiers to support custom use-cases (guide coming soon).

### When to use Curio?

As we often say, security shouldn’t be an afterthought. Curio is Open Source to allow any company, of any size, to start a more secure journey and protect their customer's data in the best way possible. Today is a good day to start!

### Who is using Curio?

Curio is used by developers and app sec engineers, and has been built as part of
Bearer with some of the best tech companies in the world to help them protect highly sensitive data.


What if you don't process sensitive data? Well, that's pretty rare these days. Even if you discover no usage after a one-time scan, Curio is a great addition to your CI tooling to ensure sensitive data processing isn't introduced into your codebase.

### How is it different from X solution?

Static code analysis tools for security are not new and many great solutions exist including Snyk Code, CodeQL, Semgrep, SonarQube, Brakeman, etc.

Curio is different as it emphasizes that the most asset we need to protect today is sensitive data, not “systems.” At the same time, it’s built with the knowledge that developers are not security experts and can’t dedicate more time to fix a growing number of security issues without a clear prioritization and understanding of the impact.

By putting the **focus on sensitive data first**, Curio is able to trigger security risks and vulnerabilities that are by design top priority, with much less noise.

The spirit of Curio is the less (alerts) is more (security) philosophy.

## Get in touch

Thanks for using Curio. Still have questions?

- Start with the [documentation](https://curio.sh).
- Have a question or need some help? Find the Curio team on [Discord][discord].
- Got a feature request or found a bug? [Open a new issue](https://github.com/Bearer/curio/issues/new/choose).
- Found a security issue? Check out our [Security Policy](https://github.com/Bearer/curio/security/policy) for reporting details.
- Find out more at [Bearer.com](https://www.bearer.com)

## Contributing

Interested in contributing? We're here for it! For details on how to contribute, setting up your development environment, and our processes, review the [contribution guide](CONTRIBUTING.md).

## Code of conduct

Everyone interacting with this project is expected to follow the guidelines of our [code of conduct](CODE_OF_CONDUCT.md).

## Security

To report a vulnerability or suspected vulnerability, [see our security policy](https://github.com/Bearer/curio/security/policy). For any questions, concerns or other security matters, feel free to [open an issue](https://github.com/Bearer/curio/issues/new/choose) or join the [Discord Community][discord].

## License

Curio code is licensed under the terms of the [Elastic License 2.0](LICENSE.txt) (ELv2), which means you can use it freely inside your organization to protect your applications without any commercial requirements.

You are not allowed to provide Curio to third parties as a hosted or managed service without the explicit approval of Bearer Inc.

---

[test]: https://github.com/Bearer/curio/actions/workflows/test.yml
[test-img]: https://github.com/Bearer/curio/actions/workflows/test.yml/badge.svg
[release]: https://github.com/Bearer/curio/releases
[release-img]: https://img.shields.io/github/release/Bearer/curio.svg?logo=github
[github-all-releases-img]: https://img.shields.io/github/downloads/Bearer/curio/total?logo=github
[discord]: https://discord.gg/eaHZBJUXRF
