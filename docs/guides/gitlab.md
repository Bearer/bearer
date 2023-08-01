---
title: Using GitLab CI/CD
---
{% renderTemplate "md" %}

# Using GitLab CI/CD


Running Bearer from the CLI is great, but if you want it integrated directly with your Git workflow there's nothing easier than a GitLab CI/CD integration. If you're unfamiliar with GitLab CI/CD, here's a [primer available from GitLab CI/CD](https://docs.gitlab.com/ee/ci/). You can also see how the integration works directly on our [Bear Publishing example app](https://gitlab.com/bearer/bear-publishing/-/blob/main/.gitlab-ci.yml).

## Getting started

To integrate Bearer CLI with GitLab CI/CD, we recommend using the docker entrypoint method. Edit your existing `.gitlab-ci.yml` file or add one to your repository root, then add the following lines:

```yml
bearer:
  image:
    name: bearer/bearer
    entrypoint: [ "" ]

  script: bearer scan .
```

This tells GitLab to use the `bearer/bearer` docker image. You can adjust the `script` key to [customize the scan](/guides/configure-scan/) with flags the same way as a local installation. An example of this file is available in [our example GitLab repo](https://gitlab.com/bearer/bear-publishing/-/tree/main).

GitLab's guide on [Running CI/CD jobs in Docker containers](https://docs.gitlab.com/ee/ci/docker/using_docker_images.html) provides additional context on configuring the CI in this way.

### Enable GitLab security scanning integration

GitLab offers an integrated security scanner that can take results from Bearer CLI's scan and add them to your repository's Security and Compliance page.

![Bearer CLI security report in GitLab security results](/assets/img/gitlab-code-scanning.jpg)

To take advantage of this, you'll need a GitLab plan that supports it. Then, you can configure your `.gitlab-ci.yml` file with Bearer CLI's special format type.

```yml
bearer:
  image:
    name: bearer/bearer
    entrypoint: [ "" ]
  script:
    - bearer scan . --format gitlab-sast --output gl-sast-report.json

  artifacts:
    reports:
      sast: gl-sast-report.json
```

These changes set the format to `gitlab-sast` and write an artifact that GitLab can use. Once run, the results of the security scan will display in the Security and Compliance section of the repository.

### Gitlab Merge Request Comments

Bearer CLI supports [Reviewdog](https://github.com/reviewdog/reviewdog) rdjson format so you can get direct feedback on your merge requests.

![Bearer CLI results in Gitlab MR](/assets/img/gl-mr-review.png)

To keep the thing in one job we download each binary then run the two commands individually.

```yml
pr_check:
  variables:
    SHA: $CI_COMMIT_SHA
    CURRENT_BRANCH: $CI_COMMIT_REF_NAME
    DEFAULT_BRANCH: $CI_DEFAULT_BRANCH
  script:
    - curl -sfL https://raw.githubusercontent.com/Bearer/bearer/main/contrib/install.sh | sh -s -- -b /usr/local/bin
    - curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh | sh -s -- -b /usr/local/bin
    - bearer scan . --format=rdjson --output=rd.json
    - cat rd.json | reviewdog -f=rdjson -reporter=gitlab-mr-discussion
```

[Don't forget](https://github.com/reviewdog/reviewdog#reporter-gitlab-mergerequest-discussions--reportergitlab-mr-discussion) to set `REVIEWDOG_GITLAB_API_TOKEN` in your project environment variables with a personal API access token.

## Make the most of Bearer

For more ways to use Bearer, check out the different [report types](/explanations/reports/), [available rules](/reference/rules/), [supported data types](/reference/datatypes/).

Have a question or need help? Join our [Discord community](https://discord.gg/eaHZBJUXRF) or [open an issue on GitHub](https://github.com/Bearer/bearer/issues).
{% endrenderTemplate %}
