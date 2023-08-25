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

This tells GitLab to use the `bearer/bearer` docker image. You can adjust the `script` key to [customize the scan](/guides/configure-scan/) with flags the same way as a local installation. An example of this file is available in [our example GitLab repo](https://gitlab.com/bearer/bear-publishing/-/tree/main).

GitLab's guide on [Running CI/CD jobs in Docker containers](https://docs.gitlab.com/ee/ci/docker/using_docker_images.html) provides additional context on configuring the CI in this way.

For more details and additional configuration, see our [guide to using GitLab](/guides/gitlab/).

## CircleCI

To integrate with CircleCI, you can add the following job to your `.circleci/config.yml`

```yml
version: 2.1

jobs:
  bearer:
    machine:
      image: ubuntu-2204:2023.07.2
    environment:
      # Set to default branch of your repo
      DEFAULT_BRANCH: main
    steps:
      - checkout
      - run: curl -sfL https://raw.githubusercontent.com/Bearer/bearer/main/contrib/install.sh | sh -s -- -b /tmp
      - run: CURRENT_BRANCH=$CIRCLE_BRANCH SHA=$CIRCLE_SHA1 /tmp/bearer scan .

workflows:
  test:
    jobs:
      - bearer
```

A more advanced example using a Github repository and reviewdog for PR comments:

```yml
version: 2.1

jobs:
  bearer:
    machine:
      image: ubuntu-2204:2023.07.2
    environment:
      # Set to default branch of your repo
      DEFAULT_BRANCH: main
    steps:
      - checkout
      - run: curl -sfL https://raw.githubusercontent.com/Bearer/bearer/main/contrib/install.sh | sh -s -- -b /tmp
      - run: curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh | sh -s -- -b /tmp
      - run: |
          CURRENT_BRANCH=$CIRCLE_BRANCH SHA=$CIRCLE_SHA1 /tmp/bearer scan . --format=rdjson --output=rd.json || export BEARER_EXIT=$?
          cat rd.json | REVIEWDOG_GITHUB_API_TOKEN=$GITHUB_TOKEN /tmp/reviewdog -f=rdjson -reporter=github-pr-review
          exit $BEARER_EXIT

workflows:
  test:
    jobs:
      - bearer:
          filters:
            branches:
              # No need to run a check on default branch
              ignore: main
          context:
            - bearer
            # make sure to set GITHUB_TOKEN in your context

```

The `GITHUB_TOKEN` in this case just requires read and write access to pull requests for the repository.

{% callout "warn" %}
Currently DEFAULT_BRANCH is hard coded and diff scanning is not supported because base branch information is not available in Circle CI.
In the future we hope to support diff scanning in Circle CI by having the CLI call the Github API for the details.
{% endcallout %}

## Universal setup

For other services, we recommend selecting the [installation method](/reference/installation/) that best fits the platform.

Do you have a CI/CD workflow that you'd like to see added to this guide? [Open an issue]({{meta.links.issues}}) or let us know on [discord]({{meta.links.discord}}).
