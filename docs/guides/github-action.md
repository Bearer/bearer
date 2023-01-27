---
title: Using the GitHub Action
---

# Using the GitHub Action

Running Curio from the CLI is great, but if you want it integrated directly with your Git workflow there's nothing easier than a GitHub action. If you're unfamiliar with GitHub actions, here's a [primer available from GitHub](https://github.com/features/actions). You can also see how the action works directly on our [Bear Publishing example app](https://github.com/Bearer/bear-publishing/actions/workflows/curio.yml).

## Getting started

You can [view the action here](https://github.com/marketplace/actions/curio-sh), or follow along below.

Actions live in the `.github/workflows/` directory within your repository. Start by creating a `curio.yml` file in the workflows directory.

We recommend the following config in `.github/workflows/curio.yml` to run Curio's summary report:

```yml
name: Curio

on:
  push:
    branches:
      - main

permissions:
  contents: read

jobs:
  rule_check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run Report
        id: report
        uses: bearer/curio-action@v0.1
      - id: summary
        name: Display Summary
        uses: actions/github-script@v6
        with:
          script: |
            // github does not support multiline outputs so report is encoded
            const report = decodeURIComponent(`${{ steps.report.outputs.rule_breaches }}`);
            const passed = `${{ steps.report.outputs.exit_code }}` == "0";
            if(!passed){ core.setFailed(report); }
```

This will run the summary report, display the results to the action summary screen within GitHub, and flag the action as pass or fail based on whether Curio's default rules pass or fail.

## Further configuration

Just as with the CLI app, you can configure the action to meet the needs of your project. Set custom inputs and outputs using the `with` key. Here's an example using the `config-file`, `skip-path`, and `only-rule` flags:

```yml
steps:
  - uses: actions/checkout@v3
  - name: Curio
    uses: bearer/curio-action@v0.1
    with:
      config-file: '/some/path/curio.yml'
      only-rule: 'ruby_lang_cookies,ruby_lang_http_post_insecure_with_data'
      skip-path: 'users/*.go,users/admin.sql'
```

The following are a list of available inputs and outputs:

### Inputs

#### `config-file`

**Optional** Curio configuration file path

#### `only-rule`

**Optional** Specify the comma-separated IDs of the rules to run; skips all other rules.

#### `skip-rule`

**Optional** Specify the comma-separated IDs of the rules to skip; runs all other rules.

#### `only-detector`

**Optional** Specify the comma-separated IDs of the detectors to run; skips all other detectors

#### `skip-detector`

**Optional** Specify the comma-separated IDs of the detectors to skip; runs all other detectors

#### `skip-path`

**Optional** Specify the comma-separated paths to skip. Supports wildcard syntax, e.g. `users/*.go,users/admin.sql`

### Outputs

#### `rule_breaches`

Details of any rule breaches that occur. This is URL encoded to work round GitHub issues with multiline outputs.

#### `exit_code`

Exit code of the curio binary, 0 indicates a pass

## Make the most of Curio

For more ways to use Curio, check out the different [report types](https://curio.sh/explanations/reports/), [available rules](https://curio.sh/reference/rules/), [supported data types](https://curio.sh/reference/datatypes/). 

Have a question or need help? Join our [Discord community](https://discord.gg/eaHZBJUXRF) or [open an issue on GitHub](https://github.com/Bearer/curio/issues).