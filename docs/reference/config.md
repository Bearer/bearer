---
title: Configuration
layout: layouts/doc.njk
---

# Configuring Bearer

Configuration of Bearer can be done with [flags on the `scan` command](/reference/commands/#scan), or by using a `bearer.yml` file in the project directory.

To initialize the config file, run the following from your project directory:

```bash
bearer init
```

This creates a config file in your current directory. Below is an annotated version of that file.

```yml
# Report settings
report:
    # Specify report format (json, yaml)
    format: ""
    # Specify the output path for the report.
    output: ""
    # Specify the type of report (summary, privacy). 
    report: summary
    # Specify which severities are included in the report as a comma separated string
    severity: "critical,high,medium,low,warning"
# Rule settings
rule:
    # Specify the comma-separated ids of the rules you would like to run; 
    # skips all other rules.
    only-rule: []
    # Specify the comma-separated ids of the rules you would like to skip; 
    # runs all other rules.
    skip-rule: []
# Scan settings
scan:
    # Expand context of schema classification 
    # For example, "health" will include data types particular to health
    context: ""
    # Override default data subject mapping by providing a path to a custom mapping JSON file
    data-subject-mapping: ""
    # Enable debug logs
    debug: false
    # Do not attempt to resolve detected domains during classification.
    disable-domain-resolution: true
    # Set timeout when attempting to resolve detected domains during classification.
    domain-resolution-timeout: 3s
    # Specify directories paths that contain yaml files with external rules configuration.
    external-rule-dir: []
    # Disable the cache and runs the detections again every time scan runs.
    force: false
    # Define regular expressions for better classification of private or unreachable domains
    # e.g., ".*.my-company.com,private.sh"
    internal-domains: []
    # Suppress non-essential messages
    quiet: false
    # Specify the comma separated files and directories to skip. Supports * syntax.
    skip-path: []
```

## Utilizing a custom config

By default, Bearer will look for a `bearer.yml` file in the project directory where the scan is run. Alternatively, you can use the `--config-file` flag with the scan command to reference a config file that is outside the project directory.