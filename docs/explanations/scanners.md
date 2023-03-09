---
title: Scanner Types
---

# Scanner Types

Bearer comes with two types of security scanners, SAST (default) and Secrets. 

## SAST Scanner

The SAST scanner is the default one if you don't specify any.
This scanner uses the [built-in rules](/explanations/rules) to detect various security risks and vulnerabilities in your code.

The output of the SAST scanner depends on the [report type](/explanations/reports) used, by default the security report will be selected and display the list of rules violations in your terminal.

```txt
$ bearer scan . --scanner=sast
...
CRITICAL: Sensitive data stored in a JWT detected.
https://docs.bearer.com/reference/rules/ruby_lang_jwt
To skip this rule, use the flag --skip-rule=ruby_lang_jwt

File: /bear-publishing/lib/jwt.rb:6

 3     JWT.encode(
 4       {
 5         id: user.id,
 6         email: user.email,
 7         class: user.class,
 8       },
 9       nil,
 	...
 11     )

...
```

You can see a full list of [built-in rules](/reference/rules) or create a [custom rule](/guides/custom-rule/).


## Secrets Scanner

The Secrets scanner type detects hard-coded secrets in your code. It checks for common secret patterns such as keys, tokens, and passwords using the popular [Gitleaks](https://gitleaks.io/) library.

```txt
$ bearer scan . --scanner=secrets
...
CRITICAL: Hard-coded secret detected. [CWE-798]

Detected: Password in URL
File: ../../OWASP/NodeGoat/README.md:59
```

You can see a full list of [built-in patterns](https://github.com/Bearer/bearer/blob/main/pkg/detectors/gitleaks/gitlab_config.toml).

⚠️ Secret detection patterns are not configurable today. If this is something you'd like to see, please open an [issue](https://github.com/Bearer/bearer/issues).

