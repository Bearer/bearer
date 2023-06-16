---
title: Connect to Bearer Cloud
templateEngineOverride: md
---

# Connect to Bearer Cloud

If you're looking to manage product and application code security at scale, Bearer Cloud offers at platform for teams that syncs with Bearer CLI's output. [Learn more about Bearer Cloud](https://www.bearer.com/bearer-cloud).

## Generate an API token

To connect Bearer CLI to Bearer Cloud, you'll first need to generate an API token. [Log in to Bearer Cloud](https://my.bearer.sh) and navigate to *Settings > API tokens* by selecting your user account in the top right corner, or from the link in the "Add a project" form.

![API token settings page](/assets/img/api-token.jpg)

## Add the API token to Bearer CLI

Use the API token any place where you run a scan. 

### Local projects

Use the `--api-key` flag with the `scan` command:

```bash
bearer scan project-folder --api-key=XXXXXXXX
```

### GitHub Action

Using the same setup process found in [the GitHub action guide](/guides/github-action/), configure the action to run `with` the `api-key` option. For example:

```yml
# .github/workflows/bearer.yml
name: Bearer
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
        uses: bearer/bearer-action@v2
        with:
          api-key: ${{ secrets.BEARER_TOKEN }}
```

We highly recommend using GitHub's [encrypted secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets). In the example above, the secret is named `BEARER_TOKEN`.

### GitLab CI

Set up the [GitLab CI configuration](/guides/ci-setup/#gitlab), then adjust your settings to include the `--api-key` flag with the `scan` command:

```yml
# .gitlab-ci.yml
image:
  name: bearer/bearer
  entrypoint: [ "" ]

bearer:
  script: bearer scan . --api-key=$BEARER_TOKEN
```

We recommend using [GitLab's CI/CD variables](https://docs.gitlab.com/ee/ci/variables/) to protect your token. In the example above, the variable is named `BEARER_TOKEN`.

## Import your projects

Bearer Cloud automatically captures any scans run with a valid `api-key`. Subsequent scans of the same project will update the existing project entry in the Bearer Cloud dashboard.

![Cloud dashboard](/assets/img/cloud-dashboard.jpg)

## Make the most of Bearer CLI

Looking for more ways to make the most of Bearer CLI? Browse the guides to learn how to [Configure a scan](/guides/configure-scan/), [create a custom rule](/guides/custom-rule), or [run a privacy report](/guides/privacy/).