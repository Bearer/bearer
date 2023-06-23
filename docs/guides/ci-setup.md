---
title: Set up CI/CD
---

# Set up CI/CD for Bearer CLI

Using Bearer CLI in your CI/CD pipeline works similarly to most other integrations. You can choose to run scans as part of the native CI/CD workflows of GitHub or GitLab, or roll your own support for additional third party services.

## GitHub

Bearer offers an official [GitHub Action](https://github.com/marketplace/actions/bearer-security) to connect directly with your repository. To enable it with the default settings, create a `bearer.yml` file in your `.github/workflows` directory and include the following:

```yml
steps:
  - uses: actions/checkout@v3
  - uses: bearer/bearer-action@v2
```

For more details and additional configuration, see our [guide to using the GitHub action](/guides/github-action/). To hook directly into GitHub's code scanning feature, check the [configure GitHub code scanning](/guides/github-action/#configure-github-code-scanning) section of the doc.

## GitLab

To integrate Bearer CLI with GitLab CI/CD, we recommend using the docker entrypoint method. Edit your existing `.gitlab-ci.yml` file or add one to your repository root, then add the following lines:

```yml
bearer:
  image:
    name: bearer/bearer
    entrypoint: [ "" ]

  script: bearer scan .
```

This tells GitLab to use the `bearer/bearer` docker image. You can adjust the `script` key to [customize the scan](/guides/configure-scan/) with flags the same way as a local installation. An example of this file is available in [our example GitLab repo](https://gitlab.com/cfabianski/bear-publishing/-/tree/main).

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


## Universal setup

For other services, we recommend selecting the [installation method](/reference/installation/) that best fits the platform.

Do you have a CI/CD workflow that you'd like to see added to this guide? [Open an issue]({{meta.links.issues}}) or let us know on [discord]({{meta.links.discord}}).