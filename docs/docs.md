---
title: Documentation
layout: "layouts/doc"
permalink: "/"
---

# Bearer CLI

Welcome to the Bearer CLI documentation. Bearer CLI is a static application security testing (SAST) tool that scans your source code and analyzes your [data flows](/explanations/discovery-and-classification) to discover, filter and prioritize security and privacy risks.

The CLI provides [built-in rules](/reference/rules) that check against a common set of security risks and vulnerabilities, known as [OWASP Top 10](https://owasp.org/www-project-top-ten/), and privacy risks. Here are some practical examples of what those rules look for:

- Non-filtered user input (sql injection, path traversal, etc.)
- Leakage of sensitive data through cookies, internal loggers, third-party logging services, and into analytics environments.
- Usage of weak encryption libraries or misusage of encryption algorithms.
- Unencrypted incoming and outgoing communication (HTTP, FTP, SMTP) of sensitive data.
- Hard-coded secrets and tokens.

And [many more](/reference/rules).

Bearer CLI currently supports **JavaScript / TypeScript** and **Ruby** stacks, and more will follow.

Want a quick rundown? Here's a minute and a half of what you can expect from Bearer CLI:

<iframe class="w-full aspect-video" src="https://www.youtube.com/embed/WeyPmYyP5Nc" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>

## Getting started

New to Bearer CLI? Check out the [quickstart](/quickstart/) to scan your first project. 

Ready to dive in? Bearer CLI's [scanners](/explanations/scanners/) and [reports](/explanations/reports/) are your path to analyzing security risks and vulnerabilities in your application. Check the [command reference](/reference/commands/) to configure Bearer CLI to your needs.

## Guides

Guides help you make the most of Bearer CLI so you can get up and running quickly.

- [Configure the scan command](/guides/configure-scan/)
- [GitHub action integration](/guides/github-action/)
- [Set up CI/CD](/guides/ci-setup/)
- [Create custom rule](/guides/custom-rule/)
- [Run a privacy report](/guides/privacy/)

## Explanations

Explanations dive into the rational behind Bearer CLI and explain some of its heavier concepts.
- [How Bearer CLI works](/explanations/workflow/)
- [Bearer CLI's scanner types](/explanations/scanners/)
- [Bearer CLI's report types](/explanations/reports/)
- [How Bearer CLI discovers and classifies data](/explanations/discovery-and-classification/)
- [How Bearer CLI sets severity levels](/explanations/severity/)

## Reference

Reference documents are where you'll find detailed information about each command, as well as support charges for languages, rules, datatypes, and more.

- [Installation](/reference/installation/)
- [Configuration](/reference/config/)
- [Commands](/reference/commands/)
- [Built-in Rules](/reference/rules/)
- [Supported Data Types](/reference/datatypes/)
- [Recipes](/reference/recipes/)
- [Supported Languages](/reference/supported-languages/)

## Contributing

We'd love to see the impact you can bring to Bearer CLI. Our contributing documentation will help get you started, whether you want to dive deep into the code or simply fix a typo.

- [Get started contributing to Bearer CLI](/contributing/)
- [Set up Bearer CLI locally to contribute code](/contributing/code/)
- [Help improve and apply fixes to the documentation](/contributing/docs/)
- [Add new recipes to Bearer CLI's database](/contributing/recipes/)
