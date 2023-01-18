---
title: Configuration
layout: layouts/doc.njk
---

# Configuring Curio

Configuration of Curio can be done with [flags on the `scan` command](/reference/commands/#scan), or by using a `curio.yml` file in the project directory.

To initialize the config file, run the following from your project directory:

```bash
curio init
```

This creates a config file in your current directory. Below is an annotated version of that file.

```yml
# Detector settings
detector:
    # Specify the comma-separated ids of the detectors you would like to run. 
    # Skips all other detectors.
    only-detector: []
    # Specify the comma-separated ids of the detectors you would like to skip. 
    # Runs all other detectors.
    skip-detector: []
# Policy settings
policy:
    # Specify the comma-separated ids of the policies you would like to run. 
    # Skips all other policies.
    only-policy: []
    # Specify the comma-separated ids of the policies you would like to skip. 
    # Runs all other policies.
    skip-policy: []
# Report settings
report:
    # Specify report format (json, yaml)
    format: ""
    # Specify the output path for the report.
    output: ""
    # Specify the type of report (detectors, dataflow, policies, stats). 
    report: policies
# Scan settings
scan:
    # Expand context of schema classification 
    # For example, "health" will include data types particular to health
    context: ""
    # Enable debug logs
    debug: false
    # Do not attempt to resolve detected domains during classification.
    disable-domain-resolution: true
    # Set timeout when attempting to resolve detected domains during classification.
    domain-resolution-timeout: 3s
    # Specify directories paths that contain .yaml files with external custom detectors configuration.
    external-detector-dir: []
    # Specify directories paths that contain .rego files with external policies configuration.
    external-policy-dir: []
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

By default, Curio will look for a `curio.yml` file in the project directory where the scan is run. Alternately, you can use the `--config-file` flag with the scan command to reference a config file that is outside the project directory.