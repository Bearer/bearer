---
title: Quick Start
layout: layouts/doc.njk
---

# Quick Start

Discover data security risks and vulnerabilities in only a few minutes. In this guide you will install Curio, run a scan on a local project, and view the results of a [policy report](/explanations/reports/#policy-report). Let's get started!

## Installation

Curio is available as a standalone executable binary. The latest release is available on the [releases tab](https://github.com/Bearer/curio/releases/latest), or use one of the methods below.

### Install Script

:warning: **Not working until public**: Use the [Binary](#binary) instructions in the next section :warning:

This script downloads the Curio binary automatically based on your OS and architecture.

```bash
curl -sfL https://raw.githubusercontent.com/Bearer/curio/main/contrib/install.sh | sh
```

_Defaults to `./bin` as a bin directory and to the latest releases_

If you need to customize the options, use the following to pass parameters:

```bash
curl -sfL https://raw.githubusercontent.com/Bearer/curio/main/contrib/install.sh | sh -s -- -b /usr/local/bin
```

### Binary

Download the archive file for your operating system/architecture from [here](https://github.com/Bearer/curio/releases/latest/). Unpack the archive, and put the binary somewhere in your $PATH (on UNIX-y systems, /usr/local/bin or the like). Make sure it has permission to execute:

```bash
chmod +x ./curio
```

## Scan your project

The easiest way to try out Curio is with our example project, [Bear Publishing](https://github.com/Bearer/bear-publishing). It simulates a realistic Ruby application with common data security flaws. Clone or download it to a convenient location to get started.  

```bash
git clone https://github.com/Bearer/bear-publishing.git
```

Now, run the scan command with `curio scan` on the project directory:

```bash
curio scan bear-publishing
```

A progress bar will display the status of the scan.

Once the scan is complete, Curio will output a policy report with details of any policy failures, as well as where in the codebase the infractions happened.

## Analyze the report

The policy report is an easily digestible view of the data security problems detected by Curio. A report is made up of:

- The list of [policies](/reference/policies/) run against your code.
- Each detected failure, containing the file location and lines that triggered the policy failure.
- A summary of the report with the stats for passing and failing policies.

The [Bear Publishing](https://github.com/Bearer/bear-publishing) example application will trigger policy failures and output a full report. Here's a section of the output containing a failure snippet and the final summary:

```text

HIGH: Application level encryption missing policy failure with PHI, PII
Application level encryption missing. Enable application level encryption to reduce the risk of leaking sensitive data.

File: /bear-publishing/db/schema.rb:22

 14 create_table "authors", force: :cascade do |t|
 15     t.string "name"
 16     t.datetime "created_at", null: false
 17     t.datetime "updated_at", null: false
 18   end

=====================================

Policy failures detected

14 policies were run and 12 failures were detected.

CRITICAL: 0
HIGH: 10 (Application level encryption missing, Insecure HTTP with Data Category,
          JWT leaking, Logger leaking, Cookie leaking, Third-party data category exposure)
MEDIUM: 2 (Insecure SMTP, Insecure FTP)
LOW: 0
```

The policy report is just one report type available in Curio. Additional options for using and configuring the `scan` command can be found in the [scan documentation](/reference/commands/#scan). For additional guides and usage tips, [view the docs](https://curio.sh).
