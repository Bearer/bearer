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

The easiest way to try out Bearer is with our example project, [Bear Publishing](https://github.com/Bearer/bear-publishing). It simulates a realistic Ruby application with common security flaws. Clone or download it to a convenient location to get started.

```bash
git clone https://github.com/Bearer/bear-publishing.git
```

Now, run the scan command with `bearer scan` on the project directory:

```bash
bearer scan bear-publishing
```

A progress bar will display the status of the scan.

Once the scan is complete, Bearer will output a security report with details of any rules failures, as well as where in the codebase the infractions happened.

By default the `scan` command uses the SAST scanner; other [scanner types](/explanations/scanners) are also available.

## Analyze the report

The security report is an easily digestible view of the security problems detected by Bearer. A report is made up of:

- The list of [rules](/reference/rules/) run against your code.
- Each detected failure, containing the file location and lines that triggered the rules failure.
- A stat section with a summary of rules checks, failures and warnings.

The [Bear Publishing](https://github.com/Bearer/bear-publishing) example application will trigger rule failures and output a full report. Here's a section of the output containing a failure snippet and the final summary:

```text
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

The security report is just one [report type](/explanations/reports/) available in Bearer.

Ready for the next step? Additional options for using and configuring the `scan` command can be found in [configuring the scan command](/guides/configure-scan/).
