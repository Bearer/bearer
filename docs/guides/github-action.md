---
title: Using GitHub Action
---

# Using GitHub Action

Running Bearer from the CLI is great, but if you want it integrated directly with your Git workflow there's nothing easier than a GitHub action. If you're unfamiliar with GitHub actions, here's a [primer available from GitHub](https://github.com/features/actions). You can also see how the action works directly on our [Bear Publishing example app](https://github.com/Bearer/bear-publishing/actions/workflows/bearer.yml).

## Getting started

You can [view the action here]({{meta.links.action}}), or follow along below.

Actions live in the `.github/workflows/` directory within your repository. Start by creating a `bearer.yml` file in the workflows directory.

We recommend the following config in `.github/workflows/bearer.yml` to run Bearer's security report:

{% yamlExample "ci/github/basic" %}

This will run the [security report](/explanations/reports), show the report in the job log, and flag the action as pass or fail based on whether Bearer's default rules pass or fail.

## Further configuration

Just as with the CLI app, you can configure the action to meet the needs of your project. Set custom inputs and outputs using the `with` key. Here's an example using the `config-file`, `skip-path`, and `only-rule` flags:

{% yamlExample "ci/github/basic-with-options" %}

### Inputs

{% githubAction bearerAction.inputs %}

### Outputs

If you want to process the output of the cli we recommend using the `output` input above to write a file that can be used elsewhere, but we also provide some basic outputs you can use if needed:

{% githubAction bearerAction.outputs %}

## Configure GitHub code scanning

Bearer CLI supports [GitHub code scanning](https://docs.github.com/en/code-security/code-scanning/automatically-scanning-your-code-for-vulnerabilities-and-errors/about-code-scanning). By using the SARIF output format, you can display [security report](/explanations/reports/#security-report) findings directly in the Security tab of your repository.

![Bearer CLI results in GitHub security tab](/assets/img/gh-code-scanning.jpg)

To enable this feature, update your action configuration to include new permissions, new format and outputs, and an additional step. Here's an example configuration:

{% yamlExample "ci/github/sarif" %}

By setting the format and output path, and adding a new upload step, the action will upload SARIF-formatted findings to GitHub's code scanner.

## Pull Request Diff

When the Bearer action is being used to check a pull request, you can tell the
action to only report findings introduced within the pull request by setting
the `diff` input parameter to `true`.

{% yamlExample "ci/github/diff" %}

See our guide on [configuring a scan](/guides/configure-scan#only-report-new-findings-on-a-branch)
for more information on differential scans.

## Code Review Comments

Bearer CLI supports [Reviewdog](https://github.com/reviewdog/reviewdog) rdjson format so you can use any of the reviewdog reporters to quickly add bearer feedback directly to your pull requests.

![Bearer CLI results in Github PR](/assets/img/gh-pr-review.png)

{% yamlExample "ci/github/diff-reviewdog" %}

## Integrate with Defect Dojo

We can monitor findings with [Defect Dojo](https://github.com/DefectDojo/django-DefectDojo) by using the `gitlab-sast` format and the v2 API. Make sure to update the instance url and set the necessary secrets.

{% yamlExample "ci/github/defect-dojo" %}

## Make the most of Bearer

For more ways to use Bearer, check out the different [report types](/explanations/reports/), [available rules](/reference/rules/), [supported data types](/reference/datatypes/).

Have a question or need help? Please [open an issue on GitHub](https://github.com/Bearer/bearer/issues).
