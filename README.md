# We are currently in the process of [renaming the project](https://github.com/Bearer/curio/issues/604) to Bearer. Some of the documentation may be out of date you can see the last stable version of [curio documentation here](https://github.com/Bearer/curio/tree/v0.24.0)

<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./docs/assets/img/curio-logo-dark.svg">
    <img alt="Curio" src="./docs/assets/img/curio-logo-light.svg">
  </picture>

  <br />
  <hr/>
  Open Source SAST to protect sensitive data in your code.
  <br />
  <br />
  Curio is a static code analysis tool (SAST) that scans your source code to discover security risks and vulnerabilities that put your sensitive data at risk (PHI, PD, PII).
  <br /><br />

  [Getting Started](#rocket-getting-started) - [FAQ](#question-faqs) - [Documentation](https://curio.sh) - [Report a Bug](https://github.com/Bearer/curio/issues/new/choose) - [Discord Community][discord]

  [![GitHub Release][release-img]][release]
  [![Test][test-img]][test]
  [![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)
  [![Discord](https://img.shields.io/discord/1042147477765242973?label=discord)][discord]

</div>

## The data-first security scanner to protect sensitive data from leaks and breaches
<hr/>

https://user-images.githubusercontent.com/1649672/208479520-3476f4d2-2fa2-4f05-abad-5ce4d3f5d3ba.mov


Curio helps developers and security teams to:

- Protect their application from leaking sensitive data (*loggers, cookies, third-parties, etc*.)
- Protect their application from having their sensitive data failed (*missing encryption, insecure communication, SQL injection, etc.*)
- Monitor sensitive data flows across every component (*Data stores, internal and external APIs*)

Curio is Open Source ([*see license*](#mortar_board-license)), and is built to be fully customizable, from creating your own rules, to adding custom code detectors up to enriching our data classifiers.

Curio also powers [Bearer](https://www.bearer.com), the developer-first platform to proactively find and fix data security risks and vulnerabilities across your development lifecycle

## :rocket: Getting started

Discover data security risks and vulnerabilities in only a few minutes. In this guide you will install Curio, run a scan on a local project, and view the results of a summary report. Let's get started!

### Install Curio

The quickest way to install Curio is with the install script. It will auto-select the best build for your architecture. _Defaults installation to `./bin` and to the latest release version_:

```bash
curl -sfL https://raw.githubusercontent.com/Bearer/curio/main/contrib/install.sh | sh
```

Or, if your platform supports it, with [Homebrew](https://brew.sh/) using [Curio's official Homebrew package](https://github.com/Bearer/homebrew-curio):

```bash
brew install Bearer/curio/curio
```

#### Debian/Ubuntu

```shell
$ sudo apt-get install apt-transport-https
$ echo "deb [trusted=yes] https://apt.fury.io/bearer/ /" | sudo tee -a /etc/apt/sources.list.d/fury.list
$ sudo apt-get update
$ sudo apt-get install curio
```

#### RHEL/CentOS

Add repository setting

```shell
$ sudo vim /etc/yum.repos.d/fury.repo
[fury]
name=Gemfury Private Repo
baseurl=https://yum.fury.io/bearer/
enabled=1
gpgcheck=0
$ sudo yum -y update
$ sudo yum -y install curio
```

[Additional installation options](#gear-additional-installation-options) are available.

### Scan your project

The easiest way to try out Curio is with our example project, [Bear Publishing](https://github.com/Bearer/bear-publishing). It simulates a realistic Ruby application with common data security flaws. Clone or download it to a convenient location to get started.

```bash
git clone https://github.com/Bearer/bear-publishing.git
```

*Alternatively, you can use your own application. Check the [supported languages](#supported-language) to see if your stack supports a summary report.*

Now, run the scan command with `curio scan` on the project directory:

```bash
curio scan bear-publishing
```

A progress bar will display the status of the scan.

Once the scan is complete, Curio will output a summary report with details of any rule failures, as well as where in the codebase the infractions happened.

### Analyze the report

The summary report is an easily digestible view of the data security problems detected by Curio. A report is made up of:

- The list of [rules](https://curio.sh/reference/rules/) run against your code.
- Each detected failure, containing the file location and lines that triggered the rule failure.
- A summary of the report with the stats for passing and failing rules.

The [Bear Publishing](https://github.com/Bearer/bear-publishing) example application will trigger rule failures and output a full report. Here's a section of the output:

```text

HIGH: Application level encryption missing rule failure with PHI, PII
Application level encryption missing. Enable application level encryption to reduce the risk of leaking sensitive data.

File: /bear-publishing/db/schema.rb:22

 14 create_table "authors", force: :cascade do |t|
 15     t.string "name"
 16     t.datetime "created_at", null: false
 17     t.datetime "updated_at", null: false
 18   end

=====================================

Rule failures detected

14 rules were run and 12 failures were detected.

CRITICAL: 0
HIGH: 10 (Application level encryption missing, Insecure HTTP with Data Category,
          JWT leaking, Logger leaking, Cookie leaking, Third-party data category exposure)
MEDIUM: 2 (Insecure SMTP, Insecure FTP)
LOW: 0
```

The summary report is just one report type available in Curio. Additional options for using and configuring the `scan` command can be found in the [scan documentation](https://curio.sh/reference/commands/#scan). For additional guides and usage tips, [view the docs](https://curio.sh).

## :gear: Additional installation options

### Configurable install script

The default installation script works well for most use cases, but if you need more control you can customize the options by passing additional parameters.

- `-b`: sets the installation directory (defaults to `./bin`)
- `-d`: enables debug logging
- `[tag]`: specifies a tag release (defaults to the latest release)

For example:

```bash
curl -sfL https://raw.githubusercontent.com/Bearer/curio/main/contrib/install.sh | sh -s -- -b /usr/local/bin
```

### Binary

Download the archive file for your operating system/architecture from [here](https://github.com/Bearer/curio/releases/latest/). Unpack the archive, and put the binary somewhere in your $PATH (on UNIX-y systems, /usr/local/bin or the like). Make sure it has permission to execute.

### Docker

Curio is also available as a Docker image on [Docker Hub](https://hub.docker.com/r/bearer/curio) and [ghcr.io](https://github.com/Bearer/curio/pkgs/container/curio).

With docker installed, you can run the following command with the appropriate paths in place of the examples.

```bash
docker run --rm -v /path/to/repo:/tmp/scan bearer/curio:latest-amd64 scan /tmp/scan
```

Additionally, you can use docker compose. Add the following to your `docker-compose.yml` file and replace the volumes with the appropriate paths for your project:

```yml
version: "3"
services:
  curio:
    platform: linux/amd64
    image: bearer/curio:latest-amd64
    volumes:
      - /path/to/repo:/tmp/scan
```

Then, run the `docker compose run` command to run Curio with any specified flags:

```bash
docker compose run curio scan /tmp/scan --debug
```

## :question: FAQs

### What can I do with Curio?

The number one thing you can do with Curio is protect sensitive data processed by your applications from being exposed by identifying associated risks. Add Curio to your CI/CD and catch data breach and data leak risks early. This is the **Rules feature.**

Curio also helps your Compliance & Privacy team assess sensitive data flows processed through your applications, allowing them to comply with the growing number of data regulations across the World (*GDPR, CCPA, CPRA, LGPD, PIPA, HIPAA, etc.*). This is the **Data Flow Report**.

### Supported Language

Curio currently fully supports Ruby, while JavaScript/TypeScript support is coming soon. Additional languages listed below are supported only for the Data Flow Report features for now.

| Languages       | Data Flow Report | Rules        |
| --------------- | ---------------- | ------------ |
| Ruby            | ✅                | ✅            |
| JavaScript / TS | ✅                | Coming soon  |
| PHP             | ✅                | Coming later |
| Java            | ✅                | Coming later |
| C#              | ✅                | Coming later |
| Go              | ✅                | Coming later |
| Python          | ✅                | Coming later |

In addition, Curio also supports these common structured data file formats:

| Format   | Support |
| -------- | ------- |
| SQL      | ✅       |
| OpenAPI  | ✅       |
| GraphQL  | ✅       |
| Protobuf | ✅       |

### What is sensitive data?

Out of the box, Curio supports 120+ data types from sensitive data categories such as Personal Data (PD), Sensitive Personal Data, Personally Identifiable Information (PII), and Protected Health Information (PHI). Please refer to the [data types documentation](https://curio.sh/reference/datatypes/) for the complete list of supported data types and how the classification works.

You can also add new data types, change the overall taxonomy, or update the classifiers to support custom use-cases (guide coming soon).

### When to use Curio?

As we often say, security shouldn’t be an afterthought. Curio is Open Source to allow any company, of any size, to start a more secure journey and protect their customer's data in the best way possible. Today is a good day to start!

### Who is using Curio?

Curio is used by developers and app sec engineers, and has been built as part of
Bearer with some of the best tech companies in the world to help them protect highly sensitive data.


What if you don't process sensitive data? Well, that's pretty rare these days. Even if you discover no usage after a one-time scan, Curio is a great addition to your CI tooling to ensure sensitive data processing isn't introduced into your codebase without your knowledge.

### How is it different from X solution?

Static code analysis tools for security are not new and many great solutions exist including [Snyk Code](https://snyk.io/product/snyk-code/), [CodeQL](https://codeql.github.com/), [Semgrep](https://semgrep.dev/), [SonarQube](https://www.sonarsource.com/products/sonarqube/), [Brakeman](https://github.com/presidentbeef/brakeman), etc.

Curio is different as it emphasizes that the most important asset we need to protect today is sensitive data, not “systems.” At the same time, it’s built with the knowledge that developers are not security experts and can’t dedicate more time to fix a growing number of security issues without a clear prioritization and understanding of the impact.

By putting the **focus on sensitive data first**, Curio is able to trigger security risks and vulnerabilities that are by design top priority, with much less noise.

The spirit of Curio is the less (alerts) is more (security) philosophy.

## :raised_hand: Get in touch

Thanks for using Curio. Still have questions?

- Start with the [documentation](https://curio.sh).
- Have a question or need some help? Find the Curio team on [Discord][discord].
- Got a feature request or found a bug? [Open a new issue](https://github.com/Bearer/curio/issues/new/choose).
- Found a security issue? Check out our [Security Policy](https://github.com/Bearer/curio/security/policy) for reporting details.
- Find out more at [Bearer.com](https://www.bearer.com)

## :handshake: Contributing

Interested in contributing? We're here for it! For details on how to contribute, setting up your development environment, and our processes, review the [contribution guide](CONTRIBUTING.md).

## :rotating_light: Code of conduct

Everyone interacting with this project is expected to follow the guidelines of our [code of conduct](CODE_OF_CONDUCT.md).

## :shield: Security

To report a vulnerability or suspected vulnerability, [see our security policy](https://github.com/Bearer/curio/security/policy). For any questions, concerns or other security matters, feel free to [open an issue](https://github.com/Bearer/curio/issues/new/choose) or join the [Discord Community][discord].

## :mortar_board: License

Curio code is licensed under the terms of the [Elastic License 2.0](LICENSE.txt) (ELv2), which means you can use it freely inside your organization to protect your applications without any commercial requirements.

You are not allowed to provide Curio to third parties as a hosted or managed service without the explicit approval of Bearer Inc.

---

[test]: https://github.com/Bearer/curio/actions/workflows/test.yml
[test-img]: https://github.com/Bearer/curio/actions/workflows/test.yml/badge.svg
[release]: https://github.com/Bearer/curio/releases
[release-img]: https://img.shields.io/github/release/Bearer/curio.svg?logo=github
[github-all-releases-img]: https://img.shields.io/github/downloads/Bearer/curio/total?logo=github
[discord]: https://discord.gg/eaHZBJUXRF
