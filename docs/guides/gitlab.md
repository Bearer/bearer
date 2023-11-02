---
title: Using GitLab CI/CD
---

# Using GitLab CI/CD


Running Bearer from the CLI is great, but if you want it integrated directly with your Git workflow there's nothing easier than a GitLab CI/CD integration. If you're unfamiliar with GitLab CI/CD, here's a [primer available from GitLab CI/CD](https://docs.gitlab.com/ee/ci/). You can also see how the integration works directly on our [Bear Publishing example app](https://gitlab.com/bearer/bear-publishing/-/blob/main/.gitlab-ci.yml).

## Getting started

To integrate Bearer CLI with GitLab CI/CD, we recommend using the docker entrypoint method. Edit your existing `.gitlab-ci.yml` file or add one to your repository root, then add the following lines:

{% yamlExample "ci/gitlab/basic" %}

This tells GitLab to use the `bearer/bearer` docker image. You can adjust the `script` key to [customize the scan](/guides/configure-scan/) with flags the same way as a local installation. An example of this file is available in [our example GitLab repo](https://gitlab.com/bearer/bear-publishing/-/tree/main).

GitLab's guide on [Running CI/CD jobs in Docker containers](https://docs.gitlab.com/ee/ci/docker/using_docker_images.html) provides additional context on configuring the CI in this way.

### Enable GitLab security scanning integration

GitLab offers an integrated security scanner that can take results from Bearer CLI's scan and add them to your repository's Security and Compliance page.

![Bearer CLI security report in GitLab security results](/assets/img/gitlab-code-scanning.jpg)

To take advantage of this, you'll need a GitLab plan that supports it. Then, you can configure your `.gitlab-ci.yml` file with Bearer CLI's special format type.

{% yamlExample "ci/gitlab/sast" %}

These changes set the format to `gitlab-sast` and write an artifact that GitLab can use. Once run, the results of the security scan will display in the Security and Compliance section of the repository.

### Gitlab Merge Request Diff

When Bearer CLI is being used to check a merge request, you can tell the Bearer
CLI to only report findings introduced within the merge request by setting the
`DIFF_BASE_BRANCH` variable.

{% yamlExample "ci/gitlab/diff" %}

See our guide on [configuring a scan](/guides/configure-scan#only-report-new-findings-on-a-branch)
for more information on differential scans.

### Gitlab Merge Request Comments

Bearer CLI supports [Reviewdog](https://github.com/reviewdog/reviewdog) rdjson format so you can get direct feedback on your merge requests.

![Bearer CLI results in Gitlab MR](/assets/img/gl-mr-review.png)

To keep the thing in one job we download each binary then run the two commands individually.

{% yamlExample "ci/gitlab/diff-reviewdog" %}

[Don't forget](https://github.com/reviewdog/reviewdog#reporter-gitlab-mergerequest-discussions--reportergitlab-mr-discussion) to set `REVIEWDOG_GITLAB_API_TOKEN` in your project environment variables with a personal API access token.

## Make the most of Bearer

For more ways to use Bearer, check out the different [report types](/explanations/reports/), [available rules](/reference/rules/), [supported data types](/reference/datatypes/).

Have a question or need help? Join our [Discord community](https://discord.gg/eaHZBJUXRF) or [open an issue on GitHub](https://github.com/Bearer/bearer/issues).

