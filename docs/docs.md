---
title: Documentation
layout: "layouts/doc"
permalink: "/"
---

# Bearer

Welcome to the Bearer documentation. Bearer is a static application security testing (SAST) tool that scans your source code and analyzes your [data flows](/explanations/discovery-and-classification) to discover, filter and prioritize security risks and vulnerabilities leading to [sensitive data](/reference/datatypes/) exposures (PII, PHI, PD).

We provides [built-in rules](/reference/rules) against a common set of security risks and vulnerabilities, known as [OWASP Top 10](https://owasp.org/www-project-top-ten/). Here are some practical examples of what those rules look for:
- Leakage of sensitive data through cookies, internal loggers, third-party logging services, and into analytics environments.
- Usage of weak encryption libraries or misusage of encryption algorithms.
- Unencrypted incoming and outgoing communication (HTTP, FTP, SMTP) of sensitive information.
- Non-filtered user input.
- Hard-coded secrets and tokens.

And [many more](/reference/rules).

Bearer currenty supporting **JavaScript** and **Ruby** stacks, more will follow.

![bearer security scanner gif](/assets/img/Bearer-security-OSS.gif)

## Getting started

New to Bearer? Check out the [quickstart](/quickstart/) to scan your first project. 

Ready to dive in? Bearer's [scanners](/explanations/scanners/) and [reports](/explanations/reports/) are your path to analyzing security risks and vulnerabilities in your application. Check the [command reference](/reference/commands/) to configure Bearer to your needs.

## Guides

Guides help you make the most of Bearer so you can get up and running quickly.

- [GitHub action integration](/guides/github-action/)
- [Create custom rule](/guides/custom-rule/)

## Explanations

Explanations dive into the rational behind Bearer and explain some of its heavier concepts.
- [Bearer's scanner types](/explanations/scanners/)
- [Bearer's report types](/explanations/reports/)
- [How Bearer discovers and classifies data](/explanations/discovery-and-classification/)

## Reference

Reference documents are where you'll find detailed information about each command, as well as support charges for languages, rules, datatypes, and more.

- [Installation](/reference/installation/)
- [Configuration](/reference/config/)
- [Commands](/reference/commands/)
- [Built-in Rules](/reference/rules/)
- [Supported Data Types](/reference/datatypes/)
- [Supported Languages](/reference/supported-languages/)

## Contributing

We'd love to see the impact you can bring to Bearer. Our contributing documentation will help get you started, whether you want to dive deep into the code or simply fix a typo.

- [Get started contributing to Bearer](/contributing/)
- [Set up Bearer locally to contribute code](/contributing/code/)
- [Help improve and apply fixes to the documentation](/contributing/docs/)
- [Add new recipes to Bearer's database](/contributing/recipes/)
  
![Sunglasses Bear](/assets/img/contribute-bearer.png)
