---
title: Documentation
layout: "layouts/doc"
permalink: "/"
---

{% callout "info" %}Discover <a href="/guides/bearer-cloud">Bearer Cloud</a>, our solution to manage product and application code security at scale.{% endcallout %}

# Bearer CLI

Welcome to the Bearer CLI documentation. Bearer CLI is a static application security testing (SAST) tool that scans your source code and analyzes your data flows to discover, filter and prioritize security and privacy risks.

This includes:
* **Security risks and vulnerabilities** using [built-in rules](https://docs.bearer.com/reference/rules/) covering the [OWASP Top 10](https://owasp.org/www-project-top-ten/) and [CWE Top 25](https://cwe.mitre.org/top25/archive/2023/2023_top25_list.html), such as:
  * A01: Access control (e.g. Path Traversal, Open Redirect, Exposure of Sensitive Information).
  * A02: Cryptographic Failures (e.g. Weak Algorithm, Insecure Communication).
  * A03: Injection (e.g. SQL Injection, Input Validation, XSS, XPath).
  * A04: Design (e.g. Missing Encryption of Sensitive Data, Persistent Cookies Containing Sensitive Information).
  * A05: Security Misconfiguration (e.g. Cleartext Storage of Sensitive Information in a Cookie or JWT).
  * A07: Identification and Authentication Failures (e.g. Use of Hard-coded Password, Improper Certificate Validation).
  * A08: Data Integrity Failures (e.g. Deserialization of Untrusted Data).
  * A09: Security Logging and Monitoring Failures (e.g. Insertion of Sensitive Information into Log File).
  * A10: Server-Side Request Forgery (SSRF).

* **Privacy risks** with the ability to detect [sensitive data flow](https://docs.bearer.com/explanations/discovery-and-classification/) such as the use of PII, PHI in your app, and [components](https://docs.bearer.com/reference/recipes/) processing sensitive data (e.g. databases like pgSQL, third-party APIs such as OpenAI, Sentry, etc.). This helps generate a [privacy report](https://docs.bearer.com/guides/privacy/) relevant for:
  * Privacy Impact Assessment (PIA).
  * Data Protection Impact Assessment (DPIA).
  * Records of Processing Activities (RoPA) input for GDPR compliance reporting.

Bearer CLI currently supports **JavaScript, TypeScript**, **Ruby**, and **Java** stacks, and more will follow.

Want a quick rundown? Here's a minute and a half of what you can expect from Bearer CLI:

<iframe class="w-full aspect-video" src="https://www.youtube.com/embed/WeyPmYyP5Nc" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>

## Getting started

New to Bearer CLI? Check out the [quickstart](/quickstart/) to scan your first project.

Ready to dive in? Bearer CLI's [scanners](/explanations/scanners/) and [reports](/explanations/reports/) are your path to analyzing security risks and vulnerabilities in your application. Check the [command reference](/reference/commands/) to configure Bearer CLI to your needs.

## Guides

Guides help you make the most of Bearer CLI so you can get up and running quickly.

- [Configure the scan command](/guides/configure-scan/)
- [GitHub action integration](/guides/github-action/)
- [GitLab CI/CD](/guides/gitlab/)
- [Set up CI/CD](/guides/ci-setup/)
- [Create custom rule](/guides/custom-rule/)
- [Run a privacy report](/guides/privacy/)
- [Run a data flow report](/guides/dataflow/)
- [Using Bearer Cloud](/guides/bearer-cloud/)
- [Enable Completion Script](/guides/shell-completion)

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
