---
title: Quick Start
layout: layouts/doc.njk
---

# Quick Start

Discover data security risks and vulnerabilities in only a few minutes. In this guide you will install Bearer, run a scan on a local project, and view the results of a [summary report](/explanations/reports/#summary-report). Let's get started!

## Installation

The quickest way to install Bearer is with the install script. It will auto-select the best build for your architecture. _Defaults installation to `./bin` and to the latest release version_:

```bash
curl -sfL https://raw.githubusercontent.com/Bearer/bearer/main/contrib/install.sh | sh
```

Or, if your platform supports it, with [Homebrew](https://brew.sh/) using [Bearer's official Homebrew package](https://github.com/Bearer/homebrew-bearer):

```bash
brew install Bearer/bearer/bearer
```

If you need more control or another way to install Bearer, we offer more [advanced installation options](https://github.com/Bearer/bearer#gear-additional-installation-options).

## Scan your project

The easiest way to try out Bearer is with our example project, [Bear Publishing](https://github.com/Bearer/bear-publishing). It simulates a realistic Ruby application with common data security flaws. Clone or download it to a convenient location to get started.  

```bash
git clone https://github.com/Bearer/bear-publishing.git
```

Now, run the scan command with `bearer scan` on the project directory:

```bash
bearer scan bear-publishing
```

A progress bar will display the status of the scan.

Once the scan is complete, Bearer will output a summary report with details of any rules failures, as well as where in the codebase the infractions happened.

## Analyze the report

The summary report is an easily digestible view of the data security problems detected by Bearer. A report is made up of:

- The list of [rules](/reference/rules/) run against your code.
- Each detected failure, containing the file location and lines that triggered the rules failure.
- A summary of the report with the stats for passing and failing rules.

The [Bear Publishing](https://github.com/Bearer/bear-publishing) example application will trigger rule failures and output a full report. Here's a section of the output containing a failure snippet and the final summary:

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

The summary report is just one report type available in Bearer. Additional options for using and configuring the `scan` command can be found in the [scan documentation](/reference/commands/#scan). For additional guides and usage tips, [view the docs](/).
