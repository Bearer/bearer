---
title: Contributing to documentation
---

# Contribute to Curio's documentation

If you're interested in contributing to Curio's documentation, this guide will help you do it. If you haven't already, review the [contributing overview](/contributing/) for different ways you can contribute.

## Contribute directly on GitHub

If you're making a change to an existing page, or fixing something small like a typo, you can do so directly on GitHub. Select the "Edit this page" link in the right column on any page in the docs to edit the page's source.

## Setting up local development

The documentation is built with [eleventy](https://www.11ty.dev) and [tailwindcss](https://tailwindcss.com/), and written primarily in markdown. To get started, you'll need NodeJS 12 or newer.

Once you've forked and cloned the repo, navigate to the docs directory, and `npm install` to bootstrap the dependencies.

```bash
# github.com/bearer/curio/docs
npm install
```

With dependencies installed, you can start a local dev server by running the `start` command:

```bash
npm start
```

## Getting help

If you're unsure where to start, have questions, or need help contributing to the documentation, join our [community discord]({{meta.links.discord}}) and we'll be happy to help out.