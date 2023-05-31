---
title: Run a privacy report
---

# Run a privacy report

Bearer CLI's [privacy report type](/explanations/reports/#privacy-report) allows anyone on your team to generate information about the data types processed by your app, as well as any third-party libraries or APIs that receive data. This enables you to do things like:

- Automate evidence collection to fulfill compliance requirements, such as a record of processing activities (ROPA) for GDPR.
- Create snapshots of data processing over time.
- Compare processing across applications.
- Identify third parties that receive sensitive data.
- Spot unexpected sensitive data use, like PII or HPI, in code that isn't expected to contain it.

## Getting started

If you haven't already, install Bearer CLI using the instructions on the [installation page](/reference/installation/) or the [quick start](/quickstart/). 

To run your first privacy report, navigate to the project root and use the `bearer scan` command with the `--report privacy` flag:

```bash
bearer scan . --report privacy
```

## Privacy report output

To get a better sense of what's included in the report, lets look at the output from the `juice-shop` example project. We've set the `--format` flag to `json` to make it more readable. If you're following along, run the following command.

```bash
bearer scan . --report privacy --format json
```

If you have [jq](https://stedolan.github.io/jq/) installed, it can format the output for your to make it easier to browse.

```bash
bearer scan . --report privacy --format json | jq
```

Here's a portion of the output. Notice that it is broken down into two main sets: *subjects* and *third_party*.

```json
{
  "subjects": [
    // ...
    {
      "subject_name": "User",
      "name": "Physical Address",
      "detection_count": 1,
      "critical_risk_failure_count": 0,
      "high_risk_failure_count": 0,
      "medium_risk_failure_count": 0,
      "low_risk_failure_count": 0,
      "rules_passed_count": 21
    },
    // ...
    {
      "subject_name": "Unknown",
      "name": "Interactions",
      "detection_count": 1,
      "critical_risk_failure_count": 0,
      "high_risk_failure_count": 0,
      "medium_risk_failure_count": 0,
      "low_risk_failure_count": 0,
      "rules_passed_count": 21
    }
  ],
  "third_party": [
    {
      "third_party": "Amazon AWS APIs",
      "subject_name": "Unknown",
      "data_types": [
        "Unknown"
      ],
      "critical_risk_failure_count": 0,
      "high_risk_failure_count": 0,
      "medium_risk_failure_count": 0,
      "low_risk_failure_count": 0,
      "rules_passed_count": 0
    }
  ]
}
```

Subjects are the individuals Bearer CLI associates with [data types](/reference/datatypes/). This information is derived from the [discovery and classification engine](/explanations/discovery-and-classification/). You can adjust this by overriding the [subject mapping](#subject-mapping). In this instance, Bearer CLI found some data types associated with the standard *User* subject, and one datatype with an *Unknown* subject. You can see the `detection_count` for each, as well as statistics for any rule findings linked to the datatype detection.

In addition to subjects, Bearer CLI found one third party. In this case, the application is interacting with AWS. It's not clear if any datatypes are sent to AWS. In instances where the scan makes a clearer detection, the report will include the subject and data types sent to third parties. For example, in the [bear publishing](https://github.com/Bearer/bear-publishing) example app, a scan detects email addresses sent to Sentry.

```json
{
  "third_party": "Sentry",
  "subject_name": "User",
  "data_types": [
    "Email Address"
  ],
  "critical_risk_failure_count": 0,
  "high_risk_failure_count": 1,
  "medium_risk_failure_count": 0,
  "low_risk_failure_count": 0,
  "rules_passed_count": 1
}
```

By default, the privacy report outputs in CSV format to the terminal. If you're handing this off to a member of your privacy or legal team, you'll likely want to output to a file.

```bash
bearer scan . --report privacy --output report.csv
```

This will allow team members to import the report into spreadsheets or their preferred data governance platform.

## Subject mapping

Bearer CLI uses "User" as the default data subject. To override this, you can copy the [subject_mapping.json](https://github.com/bearer/bearer/blob/main/pkg/classification/db/subject_mapping.json) and customize it to your needs. Then, use the `--data-subject-mapping` flag to use your mappings instead. This will use your supplied mapping file instead of the default.

```bash
bearer scan . --report privacy --data-subject-mapping /path/to/mappings.json
```

This is useful when your team has different terms for data subjects, or multiple groups of subjects, such as "customers", "employees", or "patients".

## Next steps

For more ways to make the most of our Bearer CLI, see our guide on [configuring the scan](/guides/configure-scan/) and the [commands reference](/reference/commands/). Need additional help? [Open an issue]({{meta.links.issues}}) or join our [Discord community]({{meta.links.discord}}).