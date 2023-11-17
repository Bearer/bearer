---
title: Run a data flow report
---

# Run a data flow report

Bearer CLI's [data flow report type](/explanations/reports/#data-flow-report) provides additional details about the underlying detections and classifications Bearer CLI makes during a scan. It's a good next step if you need more details than the security report offers.

The data flow report is made up of four sections:

- Data types: For any data type Bearer detects, the report will display their categories, each instance within the codebase, and associated model and field information if available.
- Risks: For languages that support security report rules, the report will include a `risks` section with information about any findings and their associated location within the codebase.
- Components: Resources such as internal data stores, cloud data stores, and third party APIs are displayed alongside their detection location.
- Dependencies: A list of dependencies detected in package managers, along with their linked version.

## Getting started

If you haven't already, install Bearer CLI using the instructions on the [installation page](/reference/installation/) or the [quick start](/quickstart/).

To run your first data flow report, navigate to the project root and use the `bearer scan` command with the `--report dataflow` flag:

```bash
bearer scan . --report dataflow
```

## Data flow report output

To get a better sense of what's included in the report, lets look at the output from the `juice-shop` example project. If you're following along, run the following command.

```bash
bearer scan . --report dataflow
```

If you have [jq](https://stedolan.github.io/jq/) installed, it can format the output for your to make it easier to browse.

```bash
bearer scan . --report dataflow | jq
```

Below is an excerpt of the `data_types` portion of the data flow report. It includes the _Email address_ data type with specifics on where Bearer CLI detected it.

```json
{
  "data_types": [
    {
      "category_name": "Contact",
      "category_groups": ["PII", "Personal Data"],
      "name": "Email Address",
      "detectors": [
        {
          "name": "typescript",
          "locations": [
            {
              "filename": "data/types.ts",
              "full_filename": "data/types.ts",
              "start_line_number": 22,
              "start_column_number": 3,
              "end_column_number": 8,
              "field_name": "email",
              "object_name": "User",
              "subject_name": "User"
            }
          ]
        }
      ]
    }
  ]
}
```

Next, if the codebase [supports rules](/reference/supported-languages/), the report includes risk information based on rule findings. It includes the rule id (`detector_id`), details about the code location, and the content that triggered the rule.

```json
{
  "risks": [
    {
      "detector_id": "javascript_express_exposed_dir_listing",
      "locations": [
        {
          "full_filename": "data/static/codefixes/accessLogDisclosureChallenge_1_correct.ts",
          "filename": "data/static/codefixes/accessLogDisclosureChallenge_1_correct.ts",
          "start_line_number": 2,
          "start_column_number": 3,
          "end_line_number": 2,
          "end_column_number": 76,
          "source": {
            "start_line_number": 2,
            "start_column_number": 3,
            "end_line_number": 2,
            "end_column_number": 76,
            "content": "app.use('/ftp', serveIndexMiddleware, serveIndex('ftp', { icons: true }))"
          },
          "presence_matches": [
            {
              "name": "app.use($\u003c...\u003eserveIndex()$\u003c...\u003e)\n"
            }
          ]
        }
      ]
    }
  ]
}
```

Next, the `components` section contains third party services and local or external datastores. Below, we can see that the project uses an AWS API. Bearer CLI detected it in a configuration file within the codebase.

```json
{
  "components": [
    {
      "name": "Amazon AWS APIs",
      "type": "external_service",
      "sub_type": "third_party",
      "locations": [
        {
          "detector": "yaml_config",
          "full_filename": "config/oss.yml",
          "filename": "config/oss.yml",
          "line_number": 36
        }
      ]
    }
    // ...
  ]
}
```

Finally, the `dependencies` portion analyzes package details for each language and collects details about a project's installed dependencies.

```json
{
  "dependencies": [
    {
      "name": "@angular-devkit/build-angular",
      "version": "15.0.4",
      "filename": "frontend/package.json",
      "detector": "package-json"
    }
    // ...
  ]
}
```

## Write the report to a file

You may find the need to write the data flow report to a file for parsing or consumption by other tools. Use the `--output` flag to write the report to a local path.

```bash
bearer scan . --report dataflow --output output/path/data-flow-report.json
```

## Next steps

For more ways to make the most of our Bearer CLI, see our guide on [configuring the scan](/guides/configure-scan/) and the [commands reference](/reference/commands/). Need additional help? [Open an issue]({{meta.links.issues}}) or join our [Discord community]({{meta.links.discord}}).
