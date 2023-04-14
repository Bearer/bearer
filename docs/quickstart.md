---
title: Quick Start
layout: layouts/doc.njk
---

# Quick Start

Discover your application security risks and vulnerabilities in only a few minutes. In this guide you will install Bearer, run the [SAST scanner](/explanations/scanners) on a local project, and view the results of a [security report](/explanations/reports/#security-report). Let's get started!

## Installation

The quickest way to install Bearer is with the install script. It will auto-select the best build for your architecture. _Defaults installation to `./bin` and to the latest release version_:

```bash
curl -sfL https://raw.githubusercontent.com/Bearer/bearer/main/contrib/install.sh | sh
```

Or, if your platform supports it, with [Homebrew](https://brew.sh/) using [Bearer's official Homebrew tap](https://github.com/Bearer/homebrew-tap):

```bash
brew install Bearer/tap/bearer
```

If you need more control or another way to install Bearer, we offer more [installation options](/reference/installation).

## Scan your project

The easiest way to try out Bearer is with our example project, [OWASP Juice Shop](https://github.com/juice-shop/juice-shop). It simulates a realistic Ruby application with common security flaws. Clone or download it to a convenient location to get started.

```bash
git clone https://github.com/juice-shop/juice-shop.git
```

Now, run the scan command with `bearer scan` on the project directory:

```bash
bearer scan juice-shop
```

A progress bar will display the status of the scan.

Once the scan is complete, Bearer will output a security report with details of any rules findings, as well as where in the codebase the infractions happened.

By default the `scan` command uses the SAST scanner; other [scanner types](/explanations/scanners) are also available.

## Analyze the report

The security report is an easily digestible view of the security problems detected by Bearer. A report is made up of:

- The list of [rules](/reference/rules/) run against your code.
- Each detected finding, containing the file location and lines that triggered the rules finding.
- A stat section with a summary of rules checks, findings and warnings.

The [OWASP Juice Shop](https://github.com/juice-shop/juice-shop) example application will trigger rule findings and output a full report. Here's a section of the output containing a finding snippet and the final summary:

```text
...
HIGH: Sensitive data stored in HTML local storage detected. [CWE-312]
https://docs.bearer.com/reference/rules/javascript_lang_session
To skip this rule, use the flag --skip-rule=javascript_lang_session

File: juice-shop/frontend/src/app/login/login.component.ts:102

 102       localStorage.setItem('email', this.user.email)


=====================================

58 checks, 38 findings

CRITICAL: 15 (CWE-22, CWE-798, CWE-89)
HIGH: 23 (CWE-312, CWE-327, CWE-548)
MEDIUM: 0
LOW: 0
WARNING: 0
```

The security report is just one [report type](/explanations/reports/) available in Bearer.

Ready for the next step? Additional options for using and configuring the `scan` command can be found in [configuring the scan command](/guides/configure-scan/).
