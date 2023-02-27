<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./docs/assets/img/bearer-logo-dark.svg">
    <img alt="Bearer" src="./docs/assets/img/bearer-logo-light.svg">
  </picture>

  <br />
  <hr/>
    Discover, filter, and prioritize security risks and vulnerabilities impacting your code.
  <br /><br />
  Bearer is a static application security testing (SAST) tool that scans your source code and analyzes your data flows to discover, filter and prioritize security risks and vulnerabilities leading to sensitive data exposures (PII, PHI, PD).
  <br /><br />
  Currently supporting <strong>JavaScript</strong> and <strong>Ruby</strong> stacks.
  <br /><br />

  [Getting Started](#rocket-getting-started) - [FAQ](#question-faqs) - [Documentation](https://docs.bearer.com) - [Report a Bug](https://github.com/Bearer/bearer/issues/new/choose) - [Discord Community][discord]

  [![GitHub Release][release-img]][release]
  [![Test][test-img]][test]
  [![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)
  [![Discord](https://img.shields.io/discord/1042147477765242973?label=discord)][discord]

</div>

## Code security scanner that natively filters and prioritizes security risks using sensitive data flow analysis.
<hr/>

https://user-images.githubusercontent.com/380564/221234857-4d449804-fe20-4ca4-8d4d-145b7997c379.mov


Bearer provides built-in rules against a common set of security risks and vulnerabilities, known as [OWASP Top 10](https://owasp.org/www-project-top-ten/). Here are some practical examples of what those rules look for:
* Non-filtered user input.
* Leakage of sensitive data through cookies, internal loggers, third-party logging services, and into analytics environments.
* Usage of weak encryption libraries or misusage of encryption algorithms.
* Unencrypted incoming and outgoing communication (HTTP, FTP, SMTP) of sensitive information.
* Hard-coded secrets and tokens

And many [more](https://docs.bearer.sh/reference/rules/).

Bearer is Open Source ([*see license*](#mortar_board-license)) and fully customizable, from creating your own rules to component detection (database, API) and data classification.

Bearer also powers our commercial offering, [Bearer Cloud](https://www.bearer.com), allowing security teams to scale and monitor their application security program using the same engine.

## :rocket: Getting started

Discover your most critical security risks and vulnerabilities in only a few minutes. In this guide, you will install Bearer, run a scan on a local project, and view the results. Let's get started!

### Install Bearer

The quickest way to install Bearer is with the install script. It will auto-select the best build for your architecture. _Defaults installation to `./bin` and to the latest release version_:

```bash
curl -sfL https://raw.githubusercontent.com/Bearer/bearer/main/contrib/install.sh | sh
```

#### Other install options
<details>
  <summary>Homebrew</summary>

  Using [Bearer's official Homebrew tap](https://github.com/Bearer/homebrew-tap):

  ```bash
  brew install bearer/tap/bearer
  ```
</details>

<details>
  <summary>Debian/Ubuntu</summary>

  ```shell
  $ sudo apt-get install apt-transport-https
  $ echo "deb [trusted=yes] https://apt.fury.io/bearer/ /" | sudo tee -a /etc/apt/sources.list.d/fury.list
  $ sudo apt-get update
  $ sudo apt-get install bearer
  ```
</details>

<details>
  <summary>RHEL/CentOS</summary>

  Add repository setting:
  ```shell
  $ sudo vim /etc/yum.repos.d/fury.repo
  [fury]
  name=Gemfury Private Repo
  baseurl=https://yum.fury.io/bearer/
  enabled=1
  gpgcheck=0
  ```

  Then install with yum:
  ```shell
    $ sudo yum -y update
    $ sudo yum -y install bearer
  ```
</details>

<details>
  <summary>Docker</summary>

  Bearer is also available as a Docker image on [Docker Hub](https://hub.docker.com/r/bearer/bearer) and [ghcr.io](https://github.com/Bearer/bearer/pkgs/container/bearer).

  With docker installed, you can run the following command with the appropriate paths in place of the examples.

  ```bash
  docker run --rm -v /path/to/repo:/tmp/scan bearer/bearer:latest-amd64 scan /tmp/scan
  ```

  Additionally, you can use docker compose. Add the following to your `docker-compose.yml` file and replace the volumes with the appropriate paths for your project:

  ```yml
  version: "3"
  services:
    bearer:
      platform: linux/amd64
      image: bearer/bearer:latest-amd64
      volumes:
        - /path/to/repo:/tmp/scan
  ```

  Then, run the `docker compose run` command to run Bearer with any specified flags:

  ```bash
  docker compose run bearer scan /tmp/scan --debug
  ```
</details>

<details>
  <summary>Binary</summary>

  Download the archive file for your operating system/architecture from [here](https://github.com/Bearer/bearer/releases/latest/).

  Unpack the archive, and put the binary somewhere in your $PATH (on UNIX-y systems, /usr/local/bin or the like). Make sure it has permission to execute.
</details>
<br/>

### Scan your project

The easiest way to try out Bearer is with our example project, [Bear Publishing](https://github.com/Bearer/bear-publishing). It simulates a realistic Ruby application with common security flaws. Clone or download it to a convenient location to get started.

```bash
git clone https://github.com/Bearer/bear-publishing.git
```

Now, run the scan command with `bearer scan` on the project directory:

```bash
bearer scan bear-publishing
```

A progress bar will display the status of the scan.

Once the scan is complete, Bearer will output a summary report with details of any rule failures, as well as where in the codebase the infractions happened and why.

### Analyze the report

The summary report is an easily digestible view of the security issues detected by Bearer. A report is made up of:

- The list of [rules](https://docs.bearer.sh/reference/rules/) run against your code.
- Each detected failure, containing the file location and lines that triggered the rule failure.
- A summary of the report with the stats for passing and failing rules.

The [Bear Publishing](https://github.com/Bearer/bear-publishing) example application will trigger rule failures and output a full report. Here's a section of the output:

```text
...
CRITICAL: Only communicate using SFTP connections.
https://docs.bearer.sh/reference/rules/ruby_lang_insecure_ftp

File: bear-publishing/app/services/marketing_export.rb:34

 34     Net::FTP.open(
 35       'marketing.example.com',
 36       'marketing',
 37       'password123'
  	...
 41     end


=====================================

56 checks, 10 failures, 6 warnings

CRITICAL: 7
HIGH: 0
MEDIUM: 0
LOW: 3
WARNING: 6
```

The summary report is just one report type available in Bearer. Additional options for using and configuring the `scan` command can be found in the [scan documentation](https://docs.bearer.sh/reference/commands/#scan).

For additional guides and usage tips, [view the docs](https://docs.bearer.sh).

## :question: FAQs

### How do you detect sensitive data flows from the code?

When you run Bearer on your codebase, it discovers and classifies data by identifying patterns in the source code. Specifically, it looks for data types and matches against them. Most importantly, it never views the actual values (it just can’t)—but only the code itself.

Bearer assesses 120+ data types from sensitive data categories such as Personal Data (PD), Sensitive PD, Personally identifiable information (PII), and Personal Health Information (PHI). You can view the full list in the [supported data types documentation](https://docs.bearer.sh/reference/datatypes/).

In a nutshell, our static code analysis is performed on two levels:
Analyzing class names, methods, functions, variables, properties, and attributes. It then ties those together to detected data structures. It does variable reconciliation etc.
Analyzing data structure definitions files such as OpenAPI, SQL, GraphQL, and Protobuf.

Bearer then passes this over to the classification engine we built to support this very particular discovery process.

If you want to learn more, here is the [longer explanation](https://docs.bearer.sh/explanations/discovery-and-classification/).

### When and where to use Bearer?

We recommend running Bearer in your CI to check new PR automatically for security issues, so your development team has a direct feedback loop to fix issues immediately.

You can also integrate Bearer in your CD, though we recommend to only make it fail on high criticality issues only, as the impact for your organization might be important.

In addition, running Bearer on a scheduled job is a great way to keep track of your security posture and make sure new security issues are found even in projects with low activity.

### Supported Language

Bearer currently supports JavaScript and Ruby and their associated most used frameworks and libraries. More languages will follow.

### What makes Bearer different from any other SAST tools?

SAST tools are known to bury security teams and developers under hundreds of issues with little context and no sense of priority, often requiring security analysts to triage issues. Not Bearer.

The most vulnerable asset today is sensitive data, so we start there and prioritize application security risks and vulnerabilities by assessing sensitive data flows in your code to highlight what is urgent, and what is not.

We believe that by linking security issues with a clear business impact and risk of a data breach, or data leak, we can build better and more robust software, at no extra cost.

In addition, by being Open Source, extendable by design, and built with a great developer UX in mind, we bet you will see the difference for yourself.

### How long does it take to scan my code? Is it fast?

It depends on the size of your applications. It can take as little as 20 seconds, up to a few minutes for an extremely large code base. We’ve added an internal caching layer that only looks at delta changes to allow quick, subsequent scans.

Running Bearer should not take more time than running your test suite.

### What about false positives?

If you’re familiar with other SAST tools, false positives are always a possibility.

By using the most modern static code analysis techniques and providing a native filtering and prioritizing solution on the most important issues, we believe this problem won’t be a concern when using Bearer.

## :raised_hand: Get in touch

Thanks for using Bearer. Still have questions?

- Start with the [documentation](https://docs.bearer.sh).
- Have a question or need some help? Find the Bearer team on [Discord][discord].
- Got a feature request or found a bug? [Open a new issue](https://github.com/Bearer/bearer/issues/new/choose).
- Found a security issue? Check out our [Security Policy](https://github.com/Bearer/bearer/security/policy) for reporting details.
- Find out more at [Bearer.com](https://www.bearer.com)

## :handshake: Contributing

Interested in contributing? We're here for it! For details on how to contribute, setting up your development environment, and our processes, review the [contribution guide](CONTRIBUTING.md).

## :rotating_light: Code of conduct

Everyone interacting with this project is expected to follow the guidelines of our [code of conduct](CODE_OF_CONDUCT.md).

## :shield: Security

To report a vulnerability or suspected vulnerability, [see our security policy](https://github.com/Bearer/bearer/security/policy). For any questions, concerns or other security matters, feel free to [open an issue](https://github.com/Bearer/bearer/issues/new/choose) or join the [Discord Community][discord].

## :mortar_board: License

Bearer code is licensed under the terms of the [Elastic License 2.0](LICENSE.txt) (ELv2), which means you can use it freely inside your organization to protect your applications without any commercial requirements.

You are not allowed to provide Bearer to third parties as a hosted or managed service without the explicit approval of Bearer Inc.

---

[test]: https://github.com/Bearer/bearer/actions/workflows/test.yml
[test-img]: https://github.com/Bearer/bearer/actions/workflows/test.yml/badge.svg
[release]: https://github.com/Bearer/bearer/releases
[release-img]: https://img.shields.io/github/release/Bearer/bearer.svg?logo=github
[github-all-releases-img]: https://img.shields.io/github/downloads/Bearer/bearer/total?logo=github
[discord]: https://discord.gg/eaHZBJUXRF
