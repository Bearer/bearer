---
title: Create a custom rule
---

# How to create a custom rule

Bearer rules are ways to ensure your codebase meets a set of standards for securing your code. Bearer ships with [many rules by default](/reference/rules/), but custom rules allow you to add specific requirements to suit your organization's needs.

## Getting started

Before you begin, you’ll want to have Bearer installed, and have run it successfully on a project. You can use the Bear Publishing test repo too. You’ll also want to be comfortable writing in the language of your codebase and have some familiarity with [YAML](https://yaml.org/). To make setup easier, we've also included a [starter yml template](#rule-starter) at the end of this guide.

## The rule configuration file

Each rule is a unique `yml` file. Custom rules share the same format as internal rules, so it can be helpful when creating rules to reference similar concepts in the [rules directory on GitHub]({{meta.sourcePath}}/tree/main/pkg/commands/process/settings/rules).

To better understand the structure of a rule file, let’s look at each key:

- `patterns`: See the section below for the Pattern Syntax.
- `languages`: An array of the languages the rule applies to. Available values are: `ruby`, `javascript`
- `trigger`: Defines under which conditions the rule should raise a failure. Optional.
  - `match_on`: Refers to the rule's pattern matches.
    - `presence`: Triggers if the rule's pattern is detected. (Default)
    - `absence`: Rule triggers on the absence of a pattern, but the presence of a `required_detection`. Examples include best practices such as missing configuration like forcing SSL communication. Note: rules that match on `absence` need a `required_detection` to be set.
  - `required_detection`:  Used with the `match_on: absence` trigger. Indicates which rule is required to activate the failure on the absence of the main rule.
  - `data_types_required`: Sometimes we may want a rule to trigger only for applications that process sensitive data. One example is password strength, where the rule only triggers if sensitive data types are found in the application.
    - `false`: Default. Rule triggers whether or not any data types have been detected in the application.
    - `true`: Rule only triggers if at least one data type is detected in the application.
- `severity`: This sets the lowest severity level of the rule, by default at `low`. The severity level can automatically increase based on the data type categories (PHI, PD, PDS, PII) detected depending on the rule `trigger` type. A severity level of `warning`, however, will never increase. Bearer groups rule failures by severity, and you can configure the summary report to only fail on specific severity thresholds. Severity is set for each data type group, each of which takes a severity level of `warning`, `low`, `medium`, `high`, or `critical`. A severity level of `warning` won’t cause CI to fail.
- `metadata`: Rule metadata is used for output to the summary report, and documentation for the internal rules.
  - `id`: A unique identifier. Internal rules are named `lang_framework_rule_name`. For rules targeting the language core, `lang` is used instead of a framework name. For example `ruby_lang_logger` and `ruby_rails_logger`. For custom rules, you may consider appending your org name.
  - `description`: A brief, one-sentence description of the rule. The best practice is to make this an actionable “rule” phrase, such as “Do X” or “Do not do X in Y”.
  - `cwe_id`: The associated list of [CWE](https://cwe.mitre.org/) identifers. (Optional)
  - `associated_recipe`: Links the rule to a [recipe]({{meta.sourcePath}}/tree/main/pkg/classification/db/recipes). Useful for associating a rule with a third party. Example: “Sentry” (Optional)
  - `remediation_message`: Used for internal rules, this builds the documentation page for a rule. (Optional)
  - `documentation_url`: Used to pass custom documentation URL for the summary report. This can be useful for linking to your own internal documentation or policies. By default, all rules in the main repo will automatically generate a link to the rule on [docs.bearer.com](/). (Optional)
- `auxiliary`: Allows you to define helper rules and detectors to make pattern-building more robust. Auxiliary rules contain a unique `id` and their own `patterns` in the same way rules do. You’re unlikely to use this regularly. See the [weak_encryption]({{meta.sourcePath}}/blob/a55ff8cf6334a541300b0e7dc3903d022987afb6/pkg/commands/process/settings/rules/ruby/lang/weak_encryption.yml) rule for examples. (Optional)
- `skip_data_types`: Allows you to prevent the specified data types from triggering this rule. Takes an array of strings matching the data type names. Example: “Passwords”. (Optional)
- `only_data_types`: Allows you to limit the specified data types that trigger this rule. Takes an array of strings matching the data type names. Example: “Passwords”. (Optional)


## Patterns

Patterns allow rules to look for matches in your code, much like regular expressions, but they take advantage of Bearer’s underlying data type detection capabilities.

In their most simple form, patterns look for a code match. As an example, let’s try to match the use of an unsecured FTP connection in ruby:

```yaml
patterns:
	- |
		Net::FTP.new()
	- |
		Net::FTP.open()
```

In the YAML above, we’re using two patterns. One for a new FTP connection and one for opening an FTP connection. But what if we want to match more dynamic code and check for data type access? That’s where **variables** and **filters** come in.

**Variables** add unknowns to the patterns. **Filters** add conditions and describe how to interpret variables. Let’s use them to enhance the FTP patterns above:

```yaml
 patterns:
	- |
		Net::FTP.new()
	- |
		Net::FTP.open()
	- pattern: |
			Net::FTP.open() do
				$<DATA_TYPE>
			end
		filters:
			- variable: DATA_TYPE
				detection: datatype
```

A new pattern appears! This time, it looks for sensitive data types inside the `Net::FTP.open()` block, using Bearer’s built-in `datatype` detection. To better understand what’s happening, let’s examine variables and filters in more detail.

*Note: in the example above, the third pattern uses a different YAML syntax and the `pattern` key. This is required to define filters for a pattern.*

### Variables

In the code above, `$<DATA_TYPE>` is a variable. All variables use the `$<>` syntax. The following are supported variable types:

`$<VARNAME>`: Like the example above, this is a named variable that you can link to a filter. You can use any value for the variable name. In the code below, `METHOD` is the variable name.

```yaml
patterns:
	- pattern: |
			logger.$<METHOD>()
```

`$<_>`: This format targets a portion of code where the specifics aren’t important, but the syntax is required to make a match. For example, we want a class declaration but aren’t concerned with its name as we care about what’s inside the class. For example:

```yaml
patterns:
  - pattern: |
      class $<_>
        validates :password, length: { minimum: $<LENGTH> }
      end
```

`$<...>`: Abstracts away a series of arguments, statements, fields, or characters. It can be used to “sandwich” another variable that might exist in a function argument. For example:

```yaml
patterns:
	- pattern: |
			$<CLIENT>.get($<...>$<DATA_TYPE>$<...>)
```

`$<!>`: In some instances, a pattern requires some wider context to match the exact line of code where the rule occurs. For those cases, use this variable type to explicitly mark the line for Bearer to highlight in the summary report. You’ll mostly need this for rules that target configuration files and settings, rather than logic-related code. For example:

```yaml
patterns:
  - |
    Rails.application.configure do
      $<!>config.force_ssl = false
    end
```

`$<VARNAME:type>`: This is a special type of named variable that helps Bearer's underlying engine by explicitly stating the node type in an AST. This is usually only used in special circumstances, like Ruby's method/block arguments. In this example, we need to mark `CONFIG` as an `identifier`.

```yaml
patterns:
	- pattern: |
			Devise.setup do |$<CONFIG:identifier>|
				$<CONFIG>.password_length = $<LENGTH>
			end
```

### Filters

**Filters** partner with named variables by applying conditions to them. Each filter is made up of the following keys:

- `variable`: The name of the variable. This is required, even in patterns that contain a single variable. (Required)
- Comparison keys: Use these on their own with or nested inside `either`.
  - `values`: Provide an array of values to match a variable against. Useful for specific method names and known options.
  - `less_than`: Compare the variable to the number provided with a **less than** statement.
  - `less_than_or_equal`: Compare the variable to the number provided with a *less than or equal* statement.
  - `greater_than`: Compare the variable to the number provided with a *greater than* statement.
  - `greater_than_or_equal`: Compare the variable to the number provided with a *greater than or equal* statement.
  - `regex`: Applies a regular expression test against the linked variable. This uses the [RE2 syntax](https://github.com/google/re2/wiki/Syntax).
- `not`: Inverts the results of another filter. Can be used with a single comparison key by nesting the key below `not`, or with an `either` block by nesting the block below `not`.
- `either`: Allows for multiple conditional checks. It behaves like an OR condition. You can nest any filter inside of `either`, such as `values`, `detection`, etc.
- `detection`: Detection filters rely on existing filter types, so they handle much of the logic for you.
  - `datatype`: This is the detection type you’ll most often see. It uses Bearer's scan to match any data type.
  - `insecure_url`: Useful for instances where you want to prevent unsecured HTTP requests. It explicitly matches `http://`.
  - `<auxiliary-detection-id>`: This allows you to link external and custom detection types by their id. See the `auxiliary` description in the rule config at the top of this page for more details, and the [weak_encryption rule]({{meta.sourcePath}}/blob/a55ff8cf6334a541300b0e7dc3903d022987afb6/pkg/commands/process/settings/rules/ruby/lang/weak_encryption.yml) for an example.

To better understand how filters and variables interact, see the pattern examples below.

### Pattern examples

Let’s look at some example patterns from Bearer’s core rules that use these filter and variable concepts.

In this example from `ruby_lang_cookies`, there are four patterns. They each use the `datatype` detection to check if a known data type exists in the patterns by matching against `$<DATA_TYPE>`. The second pattern uses a second variable, `$<METHOD>`, and filters it to only match the values of `permanent` or `signed`. Note that while the patterns are not connected, so you need to repeat the data type detection filter for each pattern.

```yaml
patterns:
  - pattern: |
      cookies[] = $<DATA_TYPE>
    filters:
      - variable: DATA_TYPE
        detection: datatype
  - pattern: |
      cookies.$<METHOD>[] = $<DATA_TYPE>
    filters:
      - variable: METHOD
        values:
          - permanent
          - signed
      - variable: DATA_TYPE
        detection: datatype
  - pattern: |
      cookies.permanent.signed[] = $<DATA_TYPE>
    filters:
      - variable: DATA_TYPE
        detection: datatype
  - pattern: |
      cookies.signed.permanent[] = $<DATA_TYPE>
    filters:
      - variable: DATA_TYPE
        detection: datatype
```

This next example uses a combination of the variable types. The pattern introduces the `either` filter, where it checks if `MAX_LENGTH` is less than `35` OR `MIN_LENGTH` is less than `8`.

```yaml
patterns:
  - pattern: |
      class $<_>
        $<!>devise password_length: $<MIN_LENGTH>..$<MAX_LENGTH>
      end
    filters:
      - either:
          - variable: MAX_LENGTH
            less_than: 35
          - variable: MIN_LENGTH
            less_than: 8
```

## How to run a custom rule.

Once you’ve written a custom rule, there are a few ways to tell Bearer about it.

Run scans with the `--external-rule-dir` flag.

```bash
bearer scan . --external-rule-dir /path/to/rules/
```

Add the rule to your bearer config file.

```yaml
external-rule-dir: /path/to/rules/
```

*Note: Including an external rules directory adds custom rules to the summary report. To only run custom rules, you’ll need to use the `only-rule` flag or configuration setting and pass it the IDs of your custom rule.*

## Rule best practices

1. Matching patterns in a rule cause *rule failures*. Depending on the severity level, failures can cause CI to exit and will display in the summary report. Keep this in mind when writing patterns so you don’t match a best practice condition and trigger a failed scan.
2. Lean on the built-in resources, like the data type detectors and recipes before creating complex rules.

## Rule starter

Below is the minimum-viable YAML file for creating your first custom rule. Copy it, customize it, and drop it in a directory.

```yaml
patterns:
  - pattern: |
      # YOUR CODE HERE
languages:
  - ruby
severity: high
metadata:
  id: custom_rule_name
  description: "This is an example rule created based on the tutorial."
```

## Need some help?

If you’re running into any problems or need some help, check out the [Discord Community]({{meta.links.discord}}). You can also [create a new issue]({{meta.links.issues}}) on GitHub.