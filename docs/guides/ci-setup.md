---
title: Set up CI/CD
---

# Set up CI/CD for Bearer CLI

Using Bearer CLI in your CI/CD pipeline works similarly to most other integrations. If you're using GitHub, the easiest method is with the [GitHub Action](/guides/github-action/). If you rely on another setup, see the guides below.

## GitLab

To integrate Bearer CLI with GitLab CI/CD, we recommend using the docker entrypoint method. Edit your existing `.gitlab-ci.yml` file or add one to your repository root, then add the following lines:

```yml
image: 
  name: bearer/bearer
  entrypoint: [ "" ]

bearer:
  script: bearer scan .
```

This tells GitLab to use the `bearer/bearer` docker image. You can adjust the `script` key to [customize the scan](/guides/configure-scan/) with flags the same way as a local installation. An example of this file is available in [our example GitLab repo](https://gitlab.com/cfabianski/bear-publishing/-/tree/main).

GitLab's guide on [Running CI/CD jobs in Docker containers](https://docs.gitlab.com/ee/ci/docker/using_docker_images.html) provides additional context on configuring the CI in this way.

## Universal setup

For other services, we recommend selecting the [installation method](/reference/installation/) that best fits the platform.

Do you have a CI/CD workflow that you'd like to see added to this guide? [Open an issue]({{meta.links.issues}}) or let us know on [discord]({{meta.links.discord}}).