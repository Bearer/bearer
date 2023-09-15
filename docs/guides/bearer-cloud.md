---
title: Using Bearer Cloud
---

# Bearer Cloud

If you're looking to manage product and application code security at scale, Bearer Cloud offers a platform for teams that syncs with Bearer CLI's output.

<iframe class="w-full aspect-video" src="https://youtube.com/embed/whPRe9GaY7w" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>

[Learn more about Bearer Cloud](https://www.bearer.com/bearer-cloud).

# Create an account

You can [start creating your free account](https://my.bearer.sh/users/sign_up) right now using your GitHub, GitLab, or Google SSO.

{% callout "info" %}
Bearer Cloud free plan comes with these limits:<br/>
- 1 team member<br/>
- 10 applications<br/>
- Slack integration only<br/>

Need more? <a href="https://www.bearer.com/contact">Contact us</a>.
 {% endcallout %}


# Get started with Bearer Cloud


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

```yaml
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
          api-key: {% raw %}${{ secrets.BEARER_TOKEN }}{% endraw %}
```

We highly recommend using GitHub's [encrypted secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets). In the example above, the secret is named `BEARER_TOKEN`.

### GitLab CI/CD

Set up the [GitLab CI/CD configuration](/guides/gitlab), then adjust your settings to include the `--api-key` flag with the `scan` command:

```yaml
# .gitlab-ci.yml
bearer:
  image:
    name: bearer/bearer
    entrypoint: [ "" ]
  script: bearer scan . --api-key=$BEARER_TOKEN
```

We recommend using [GitLab's CI/CD variables](https://docs.gitlab.com/ee/ci/variables/) to protect your token. In the example above, the variable is named `BEARER_TOKEN`.

## Import your projects

Bearer Cloud automatically captures any scans run with a valid `api-key`. Subsequent scans of the same project will update the existing project entry in the Bearer Cloud dashboard.

![Cloud dashboard](/assets/img/cloud-dashboard.jpg)

### Ignored findings in Bearer Cloud

When a valid `api-key` is present, the very first scan of a project reads ignored fingerprints from the ignore file and subsequently creates ignored findings for these in the Cloud, including status and comments (if present). A finding has "False Positive" status in the Cloud if its corresponding ignore file entry is a false positive (`false_positive: true`); otherwise, it has the status "Allowed".

After the initial scan, the Cloud is taken as the source of truth for ignored fingerprints. If there are new entries added to the ignore file, in most cases, these are sent to the Cloud on subsequent scans, and the corresponding Cloud findings are updated to "False Positive" or "Allowed" status accordingly.

However, it is important to note that the Cloud state is always prioritized over the contents of the ignore file. If a finding is already ignored in the Cloud, and then added to the ignore file, its Cloud status and comments are unchanged by subsequent scans. Similarly, if an ignored finding is re-opened in the Cloud, and then added to the ignore file, its Cloud status remains "Open". That is, re-opened findings can only be re-ignored again from the Cloud.

Furthermore, if an ignored finding is later re-opened in the Cloud, any corresponding ignore entry is not automatically removed. Over time, then, the ignore file may become out-of-sync with the Cloud state. To remedy this, and align the ignore file with what is in the Cloud, use the following action:

```bash
bearer ignore pull project-folder --api-key=XXXXXXXX
```

This action overwrites the current ignore file (including any new additions not yet sent to the Cloud) with all ignored findings from the Cloud, including status, comments, and author information.

## Jira integration

The Jira integration is available on the *Settings > Integrations* page.

To use the integration, you must connect a Jira account and allow access to the required permissions through the OAuth login.

Following your company's best practices, you can provide access to an existing account or set up a new user in Jira specifically for this integration. Whichever option you choose, make sure the account has the access permissions required to create and update tickets in the projects you want to.

You have two ways to use the Jira Integration:
1. Creating a Jira Ticket directly from a finding.
![Create Jira Ticket](/assets/img/jira-integration/create.png)

2. Link a finding to an existing Jira ticket.
![Link Jira Ticket](/assets/img/jira-integration/link.png)

Once a finding is associated with a Jira ticket, you can quickly see it in the interface, view the ticket status and go to the ticket.

![View Jira Ticket](/assets/img/jira-integration/view.png)


{% callout "warn" %}
Findings on Bearer Cloud are only marked resolved when the associated code is fixed. If the associated Jira ticket is closed, but no code fix has been applied, the finding will stay open. The source of truth is always the code.
{% endcallout %}


## Slack integration

The Slack integration is available on the *Settings > Integrations* page.

To use the integration, you must connect a Slack account and allow access to the required permissions through the OAuth login, then select a default channel where you want to receive notifications on new findings.

Below an example of a Slack notification triggered by a new finding:
![Slack notification](/assets/img/slack-integration/notification.png)


## Need help?

Get in touch with our team directly on [Discord](https://discord.com/invite/eaHZBJUXRF) or [book a demo](https://www.bearer.com/demo) with one of our engineer.
