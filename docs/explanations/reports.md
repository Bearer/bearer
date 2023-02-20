---
title: Report Types
---

# Report Types

Bearer can generate two types of reports about your codebase, all from the same underlying scan.

## Summary Report

Rules are the core feature of Bearer. The summary report allows you to quickly see data security vulnerabilities in your codebase. The report breaks down rule violations by severity level: Critical, High, Medium, Low, Warning. These levels help you prioritize issues and fix the most important vulnerabilities. The summary report is Bearer's default report type.

For each violation, the summary report includes the affected file and, when possible, the line of code and a snippet of the surrounding code. Here's an excerpt from the summary report run on our [example publishing app](https://github.com/Bearer/bear-publishing):

```txt
...
CRITICAL: Do not store sensitive data in JWTs. [DSR-3]
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

=====================================

24 checks, 18 failures, 0 warnings

CRITICAL: 15 (DSR-1, DSR-2, DSR-3, DSR-5)
HIGH: 0
MEDIUM: 0
LOW: 3 (DSR-2, DSR-7)
WARNING: 0

exit status 1
```

Note that if only warning-level violations are found, the report does not return an exit status of 1.

```
24 checks, 2 warnings

CRITICAL: 0
HIGH: 0
MEDIUM: 0
LOW: 0
WARNING: 2 (DSR-10)

```

Summary reports are [currently available](/reference/supported-languages/) for Ruby projects, with Javascript/Typescript coming soon and more languages to follow.

To run your first summary report, run `bearer scan .` on your project directory. By default, the summary report is output in a human-readable format, but you can also output it as YAML or JSON by using the `--format yaml` or `--format json` flags.

## Privacy Report

The privacy report provides useful information about how your codebase uses sensitive data, with an emphasis on data subjects and third parties.

The data subjects portion displays information about each detected subject and any data types associated with it. It also provides statistics on total detections and any counts of rule failures associated with the data type. In the example below, the report detects 14 instances of user telephone numbers with no rule failures.

_Note: These examples use JSON for readability, but the default format for the privacy report is CSV._

```json
"Subjects": [
  {
    "subject_name": "User",
    "name": "Telephone Number",
    "detection_count": 14,
    "critical_risk_failure_count": 0,
    "high_risk_failure_count": 0,
    "medium_risk_failure_count": 0,
    "low_risk_failure_count": 0,
    "rules_passed_count": 11
  }
]
```


The third parties portion displays data subjects and types that are sent to or processed by known third-party services. In the example below, Bearer detects a user email address sent to Sentry via the Sentry SDK and notes that a critical-risk-level rule has failed associated with this data point.

```json
"ThirdParty": [
  {
    "third_party": "Sentry",
    "subject_name": "User",
    "data_types": [
      "Email Address"
    ],
    "critical_risk_failure_count": 1,
    "high_risk_failure_count": 0,
    "medium_risk_failure_count": 0,
    "low_risk_failure_count": 0,
    "rules_passed_count": 0
  }
]
```

To run the privacy report, run `bearer scan` with the `--report privacy` flag. By default, the privacy report is output in CSV format. To format as JSON, use the `--format json` flag.

### Customizing data subjects

By default, Bearer maps all subjects to “User”, but you can override this by supplying Bearer with custom mappings. This is done by passing the path to a JSON file with the `--data-subject-mapping` flag when you run the privacy report. For example:

```bash
bearer scan . --report=privacy --data-subject-mapping=/path/to/mappings.json
```

The custom map file should follow the format used by [subject_mapping.json]({{meta.sourcePath}}/blob/main/pkg/classification/db/subject_mapping.json). Replace a key’s value with the higher-level subject you’d like to associate it with. Some examples might include Customer, Employee, Client, Patient, etc. Bearer will use your replacement file instead of the default, so make sure to include any and all subjects you want reported.
