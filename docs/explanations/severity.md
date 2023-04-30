---
title: Dynamic severity levels
---
# Dynamic Severity Levels

In order to help you decide which security findings to tackle first, Bearer CLI dynamically sets the severity of each detection. There are five levels: critical, high, medium, low, and warning. You'll see this in the security report output like in the finding below which has a severity of "High":

```txt
HIGH: Sensitive data stored in HTML local storage detected. [CWE-312]
https://docs.bearer.com/reference/rules/javascript_lang_session
To skip this rule, use the flag --skip-rule=javascript_lang_session

File: frontend/src/app/login/login.component.ts:102

 102       localStorage.setItem('email', this.user.email)
```

Bearer CLI's ability to dynamically set severity means **a single rule can trigger at multiple severity levels** depending on how *severe* or risky the code appears.

This doc explains how Bearer CLI calculates severity and the rationale behind behind the algorithm.

## Base rule severity

Each rule defines a severity level in its YML config:

```yml
severity: low
```

This is the rule's *base* or *default* severity. It is the lowest severity the rule can have, but depending on other variables that we'll get into shortly, that level can go higher. *The only exception is "warning" which is not affected by the dynamic severity system and will never increase.*

The base severity level of default rules comes from an overall assessment of the rule, but is highly influenced by community resources like the [CWE Top 25 Most Dangerous Software Weaknesses](https://cwe.mitre.org/top25/archive/2022/2022_cwe_top25.html) and similar guidance from the industry. If you're building a custom rule, we recommend using a base severity level that best fits your risk assessment for the rule.

## Sensitive data category

Bearer CLI focuses on four main sensitive data categories: Personal Health Information (PHI), Personally Identifiable Information (PII), Personal Data (PD), and Personal Data (Sensitive) (PDS). Our detection and classification engine looks at the [data types](/reference/datatypes/) and then takes the category with the highest risk into account to determine how severe a leak of that data would be.

## Trigger type

The algorithm also considers how the rule is triggered. Rule triggers are split into a few categories:

- **local**: The rule triggers if the pattern in the rule exists.
- **global**: These rules should only trigger when Bearer CLI detects data types anywhere in the codebase.
- **presence**: These rules trigger if the code exists, but isn’t dependent on data types. This could be a bad practice, but not data-specific.
- **absence**:  These rules expect code to exist, but can’t find it. Examples are best practices, must-have settings, etc.

*Note: These categories don’t map directly to configuration options in rule YAML files. See the [custom rule guide](/guides/custom-rule/) for which trigger types to use when writing your own rules.*

## Assigning weights to each variable

With these criteria in mind, we assign weights to each based on how damaging a breach of the data could be.
For base severity, the weights are as follows:

```yml
critical: 8
high: 5
medium: 3
low: 2
warning: 1
```

For sensitive data categories, the weights are as follows:

```yml
PHI: 3
PDS: 3
PD: 2
PII: 1
```

For trigger type, the weights are as follows:

```yml
local: 2
global: 1
presence: 1
absence: 1
```

## Determining final severity

The end result is a formula to assess the severity each time a rule triggers.

```txt
Final severity = base severity + (Data Category * Trigger)
```

As an example, a rule with a base severity of medium that detects PII with a local trigger would result in the following formula:

```txt
# Medium + (PII * local) = High
3 + (1 * 2) = 5
```

## Actionable prioritization

All of this is an effort to help you and developers on your team prioritize which findings to investigate first. We believe by connecting a base severity threshold with key contextual factors, Bearer CLI offers a more accurate roadmap for tackling security improvements.

If you have additional questions on how the severity system works or would like to suggest improvements, [file an issue]({{meta.links.issues}}) or [join our Discord community]({{meta.links.discord}}).
