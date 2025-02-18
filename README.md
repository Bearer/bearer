<div align="center">
  <a href="https://cycode.com/cygives/" alt="Bearer is part of Cygives, the community hub for free & open developer security tools."/>
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="./docs/assets/img/Cygives-darkmode.svg">
      <img alt="Cygives Banner" src="./docs/assets/img/Cygives-lightmode.svg">
    </picture>
  </a>
  <br/>
  <br/>
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./docs/assets/img/bearer-logo-dark.svg">
    <img alt="Bearer" src="./docs/assets/img/bearer-logo-light.svg">
  </picture>
  <br />
  <hr/>
    Scan your source code against top <strong>security</strong> and <strong>privacy</strong> risks.
  <br /><br />
  Bearer is a static application security testing (SAST) tool designed to scan your source code and analyze data flows to identify, filter, and prioritize security and privacy risks.
  <br/><br/>
  Bearer offers a free, open solution, Bearer CLI, and a commercial solution, Bearer Pro, available through <a href="https://cycode.com/">Cycode</a>.
  <br /><br />

  [Getting Started](#rocket-getting-started) - [FAQ](#question-faqs) - [Documentation](https://docs.bearer.com) - [Report a Bug](https://github.com/Bearer/bearer/issues/new/choose) - [Discord Community][discord]

  [![GitHub Release][release-img]][release]
  [![Test][test-img]][test]
  [![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)
  [![Discord](https://img.shields.io/discord/1042147477765242973?label=discord)][discord]
</div>


## Language Support
**Bearer CLI**: Go • Java • JavaScript • TypeScript • PHP • Python • Ruby <br/>
**Bearer Pro by Cycode**: C# • Kotlin • Elixir *+ everything in Bearer OSS*

<a href="https://docs.bearer.com/reference/supported-languages/">Learn more about language suppport</a>

## Developer friendly static code analysis for security and privacy
<hr/>

<https://user-images.githubusercontent.com/1649672/230438696-9bb0fd35-2aa9-4273-9970-733189d01ff1.mp4>

Bearer CLI scans your source code for:
* **Security risks and vulnerabilities** using [built-in rules](https://docs.bearer.com/reference/rules/) covering the [OWASP Top 10](https://owasp.org/www-project-top-ten/) and [CWE Top 25](https://cwe.mitre.org/top25/archive/2023/2023_top25_list.html), such as:
  * A01: Access control (e.g. Path Traversal, Open Redirect, Exposure of Sensitive Information).
  * A02: Cryptographic Failures (e.g. Weak Algorithm, Insecure Communication).
  * A03: Injection (e.g. SQL Injection, Input Validation, XSS, XPath).
  * A04: Design (e.g. Missing Encryption of Sensitive Data, Persistent Cookies Containing Sensitive Information).
  * A05: Security Misconfiguration (e.g. Cleartext Storage of Sensitive Information in a Cookie or JWT).
  * A07: Identification and Authentication Failures (e.g. Use of Hard-coded Password, Improper Certificate Validation).
  * A08: Data Integrity Failures (e.g. Deserialization of Untrusted Data).
  * A09: Security Logging and Monitoring Failures (e.g. Insertion of Sensitive Information into Log File).
  * A10: Server-Side Request Forgery (SSRF).

  *Note: all the rules and their code patterns are accessible through the [documentation](https://docs.bearer.com/reference/rules/).*

* **Privacy risks** with the ability to detect [sensitive data flow](https://docs.bearer.com/explanations/discovery-and-classification/) such as the use of PII, PHI in your app, and [components](https://docs.bearer.com/reference/recipes/) processing sensitive data (e.g. databases like pgSQL, third-party APIs such as OpenAI, Sentry, etc.). This helps generate a [privacy report](https://docs.bearer.com/guides/privacy/) relevant for:
  * Privacy Impact Assessment (PIA).
  * Data Protection Impact Assessment (DPIA).
  * Records of Processing Activities (RoPA) input for GDPR compliance reporting.

## :rocket: Getting started

Discover your most critical security risks and vulnerabilities in only a few minutes. In this guide, you will install Bearer CLI, run a security scan on a local project, and view the results. Let's get started!

### Install Bearer CLI

The quickest way to install Bearer CLI is with the install script. It will auto-select the best build for your architecture. _Defaults installation to `./bin` and to the latest release version_:

```bash
curl -sfL https://raw.githubusercontent.com/Bearer/bearer/main/contrib/install.sh | sh
```

#### Other install options

<details>
  <summary>Homebrew</summary>

  Using [Bearer CLI's official Homebrew tap](https://github.com/Bearer/homebrew-tap):

  ```bash
  brew install bearer/tap/bearer
  ```

  Update an existing installation with the following:

  ```bash
  brew update && brew upgrade bearer/tap/bearer
  ```

</details>

<details>
  <summary>Debian/Ubuntu</summary>

  ```shell
  sudo apt-get install apt-transport-https
  echo "deb [trusted=yes] https://apt.fury.io/bearer/ /" | sudo tee -a /etc/apt/sources.list.d/fury.list
  sudo apt-get update
  sudo apt-get install bearer
  ```

  Update an existing installation with the following:

  ```bash
  sudo apt-get update
  sudo apt-get install bearer
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
    sudo yum -y update
    sudo yum -y install bearer
  ```

  Update an existing installation with the following:

  ```bash
  sudo yum -y update bearer
  ```

</details>

<details>
  <summary>Docker</summary>

  Bearer CLI is also available as a Docker image on [Docker Hub](https://hub.docker.com/r/bearer/bearer) and [ghcr.io](https://github.com/bearer/bearer/pkgs/container/bearer).

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

  Then, run the `docker compose run` command to run Bearer CLI with any specified flags:

  ```bash
  docker compose run bearer scan /tmp/scan --debug
  ```

  The Docker configurations above will always use the latest release.
</details>

<details>
  <summary>Binary</summary>

  Download the archive file for your operating system/architecture from [here](https://github.com/Bearer/bearer/releases/latest/).

  Unpack the archive, and put the binary somewhere in your $PATH (on UNIX-y systems, /usr/local/bin or the like). Make sure it has permission to execute.

  To update Bearer CLI when using the binary, download the latest release and overwrite your existing installation location.
</details>
<br/>

### Scan your project

The easiest way to try out Bearer CLI is with the OWASP [Juice Shop](https://github.com/juice-shop/juice-shop) example project. It simulates a realistic JavaScript application with common security flaws. Clone or download it to a convenient location to get started.

```bash
git clone https://github.com/juice-shop/juice-shop.git
```

Now, run the scan command with `bearer scan` on the project directory:

```bash
bearer scan juice-shop
```

A progress bar will display the status of the scan.

Once the scan is complete, Bearer CLI will output, by default, a security report with details of any rule findings, as well as where in the codebase the infractions happened and why.

By default the `scan` command use the SAST scanner, other [scanner types](https://docs.bearer.com/explanations/scanners) are available.

### Analyze the report

The security report is an easily digestible view of the security issues detected by Bearer CLI. A report is made up of:

* The list of [rules](https://docs.bearer.com/reference/rules/) run against your code.
* Each detected finding, containing the file location and lines that triggered the rule finding.
* A stat section with a summary of rules checks, findings and warnings.

The [OWASP Juice Shop](https://github.com/juice-shop/juice-shop) example application will trigger rule findings and output a full report. Here's a section of the output:

```text
...
HIGH: Sensitive data stored in HTML local storage detected. [CWE-312]
https://docs.bearer.com/reference/rules/javascript_lang_session
To skip this rule, use the flag --skip-rule=javascript_lang_session

File: juice-shop/frontend/src/app/login/login.component.ts:102

 102       localStorage.setItem('email', this.user.email)


=====================================

59 checks, 40 findings

CRITICAL: 0
HIGH: 16 (CWE-22, CWE-312, CWE-798, CWE-89)
MEDIUM: 24 (CWE-327, CWE-548, CWE-79)
LOW: 0
WARNING: 0
```

In addition of the security report, you can also run a [privacy report](https://docs.bearer.com/explanations/reports/#privacy-report).

Ready for the next step? Additional options for using and configuring the `scan` command can be found in [configuring the scan command](https://docs.bearer.com/guides/configure-scan/).

For more guides and usage tips, [view the docs](https://docs.bearer.com/).

## :question: FAQs

### What makes Bearer CLI different from any other SAST tools?

SAST tools are known to bury security teams and developers under hundreds of issues with little context and no sense of priority, often requiring security analysts to triage issues manually.

The most vulnerable asset today is sensitive data, so we start there and [prioritize](https://github.com/Bearer/bearer/issues/728) findings by assessing sensitive data flows to highlight what is more critical, and what is not. This unique ability allows us to provide you with a privacy scanner too.

We believe that by linking security issues with a clear business impact and risk of a data breach, or data leak, we can build better and more robust software, at no extra cost.

In addition, by being Free and Open, extendable by design, and built with a great developer UX in mind, we bet you will see the difference for yourself.

### What is the privacy scanner?

In addition of detecting security flaws in your code, Bearer CLI allows you to automate the evidence gathering process needed to generate a privacy report for your compliance team.

When you run Bearer CLI on your codebase, it discovers and classifies data by identifying patterns in the source code. Specifically, it looks for data types and matches against them. Most importantly, it never views the actual values—it just can’t—but only the code itself. If you want to learn more, here is the [longer explanation](https://docs.bearer.com/explanations/discovery-and-classification/).

Bearer CLI is able to identify over 120+ data types from sensitive data categories such as Personal Data (PD), Sensitive PD, Personally identifiable information (PII), and Personal Health Information (PHI). You can view the full list in the [supported data types documentation](https://docs.bearer.com/reference/datatypes/).

Finally, Bearer CLI also lets you detect components storing and processing sensitive data such as databases, internal APIs, and third-party APIs. See the [recipe list](https://docs.bearer.com/reference/recipes/) for a complete list of components.

### Supported Language

Bearer CLI currently supports:
<table>
  <tr>
    <td>GA</td>
    <td>JavaScript/TypeScript, Ruby, PHP, Java, Go, Python</td>
  </tr>
  <tr>
    <td>Beta</td>
    <td>-</td>
  </tr>
  <tr>
    <td>Alpha</td>
    <td>-</td>
  </tr>
</table>

[Learn more](https://docs.bearer.com/reference/supported-languages/) about language support.

### How long does it take to scan my code? Is it fast?

It depends on the size of your applications. It can take as little as 20 seconds, up to a few minutes for an extremely large code base.

As a rule of thumb, Bearer CLI should never take more time than running your test suite.

In the case of CI integration, we provide a diff scan solution to make it even faster. [Learn more](https://docs.bearer.com/guides/configure-scan/#only-report-new-findings-on-a-branch).

### What about false positives?

If you’re familiar with SAST tools, false positives are always a possibility.

By using the most modern static code analysis techniques and providing a native filtering and prioritizing solution on the most important issues, we believe we have dramatically improved the overall SAST experience.

We strive to provide the best possible experience for our users. [Learn more](https://docs.bearer.com/reference/supported-languages/#how-do-we-evaluate-language-support%3F) about how we achieve this.

### When and where to use Bearer CLI?

We recommend running Bearer CLI in your CI to check new PRs automatically for security issues, so your development team has a direct feedback loop to fix issues immediately.

You can also integrate Bearer CLI in your CD, though we recommend setting it to only fail on high criticality issues, as the impact for your organization might be important.

In addition, running Bearer CLI as a scheduled job is a great way to keep track of your security posture and make sure new security issues are found even in projects with low activity.

Make sure to read our [integration strategy guide](https://docs.bearer.com/guides/integration-strategy/) for more information.

## :raised_hand: Get in touch

Thanks for using Bearer CLI. Still have questions?

* Start with the [documentation](https://docs.bearer.com).
* Have a question or need some help? Find the Bearer team on [Discord][discord].
* Got a feature request or found a bug? [Open a new issue](https://github.com/Bearer/bearer/issues/new/choose).
* Found a security issue? Check out our [Security Policy](https://github.com/Bearer/bearer/security/policy) for reporting details.
* Find out more at [Bearer.com](https://www.bearer.com)

## :handshake: Contributing

Interested in contributing? We're here for it! For details on how to contribute, setting up your development environment, and our processes, review the [contribution guide](CONTRIBUTING.md).

## :rotating_light: Code of conduct

Everyone interacting with this project is expected to follow the guidelines of our [code of conduct](CODE_OF_CONDUCT.md).

## :shield: Security

To report a vulnerability or suspected vulnerability, [see our security policy](https://github.com/Bearer/bearer/security/policy). For any questions, concerns or other security matters, feel free to [open an issue](https://github.com/Bearer/bearer/issues/new/choose) or join the [Discord Community][discord].

## :mortar_board: License

Bearer CLI code is licensed under the terms of the [Elastic License 2.0](LICENSE.txt) (ELv2), which means you can use it freely inside your organization to protect your applications without any commercial requirements.

You are not allowed to provide Bearer CLI to third parties as a hosted or managed service without the explicit approval of Bearer Inc.

---

[test]: https://github.com/Bearer/bearer/actions/workflows/test.yml
[test-img]: https://github.com/Bearer/bearer/actions/workflows/test.yml/badge.svg
[release]: https://github.com/Bearer/bearer/releases
[release-img]: https://img.shields.io/github/release/Bearer/bearer.svg?logo=github
[discord]: https://discord.gg/eaHZBJUXRF
