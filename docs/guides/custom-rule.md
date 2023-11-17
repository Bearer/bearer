---
title: Create a custom rule
---

# How to create a custom rule

Bearer CLI rules are ways to ensure your codebase meets a set of standards for securing your code. Bearer CLI ships with [many rules by default](/reference/rules/), but custom rules allow you to add specific requirements to suit your organization's needs.

## Getting started

Before you begin, you’ll want to have Bearer CLI installed, and have run it successfully on a project. You can use the Bear Publishing test repo too. You’ll also want to be comfortable writing in the language of your codebase and have some familiarity with [YAML](https://yaml.org/). To make setup easier, we've also included a [starter yml template](#rule-starter) at the end of this guide.

## The rule configuration file

Each rule is a unique `yml` file. Custom rules share the same format as internal rules, so it can be helpful when creating rules to reference similar concepts in the [rules repo on GitHub](https://github.com/Bearer/bearer-rules).

To better understand the structure of a rule file, let’s look at each key:

- `patterns`: See the section below for the Pattern Syntax.
- `sanitizer`: The id of an auxiliary rule which is used to restrict the
  main rule. If the sanitizer rule matches then the main rule is disabled inside
  the matched code.
- `languages`: An array of the languages the rule applies to. Available values are: `ruby`, `javascript`, `java`
- `trigger`: Defines under which conditions the rule should raise a result. Optional.
  - `match_on`: Refers to the rule's pattern matches.
    - `presence`: Triggers if the rule's pattern is detected. (Default)
    - `absence`: Rule triggers on the absence of a pattern, but the presence of a `required_detection`. Examples include best practices such as missing configuration like forcing SSL communication. Note: rules that match on `absence` need a `required_detection` to be set.
  - `required_detection`: Used with the `match_on: absence` trigger. Indicates which rule is required to activate the result on the absence of the main rule.
  - `data_types_required`: Sometimes we may want a rule to trigger only for applications that process sensitive data. One example is password strength, where the rule only triggers if sensitive data types are found in the application.
    - `false`: Default. Rule triggers whether or not any data types have been detected in the application.
    - `true`: Rule only triggers if at least one data type is detected in the application.
- `severity`: This sets the lowest severity level of the rule, by default at `low`. The severity level can [automatically increase based on multiple factors](/explanations/severity). A severity level of `warning`, however, will never increase and won’t cause CI to fail.. Bearer CLI groups rule findings by severity, and you can configure the security report to only trigger on specific severity thresholds.
- `metadata`: Rule metadata is used for output to the security report, and documentation for the internal rules.
  - `id`: A unique identifier. Internal rules are named `lang_framework_rule_name`. For rules targeting the language core, `lang` is used instead of a framework name. For example `ruby_lang_logger` and `ruby_rails_logger`. For custom rules, you may consider appending your org name.
  - `description`: A brief, one-sentence description of the rule. The best practice is to make this an actionable “rule” phrase, such as “Do X” or “Do not do X in Y”.
  - `cwe_id`: The associated list of [CWE](https://cwe.mitre.org/) identifiers. (Optional)
  - `associated_recipe`: Links the rule to a [recipe]({{meta.sourcePath}}/tree/main/internal/classification/db/recipes). Useful for associating a rule with a third party. Example: “Sentry” (Optional)
  - `remediation_message`: Used for internal rules, this builds the documentation page for a rule. (Optional)
  - `documentation_url`: Used to pass custom documentation URL for the security report. This can be useful for linking to your own internal documentation or policies. By default, all rules in the main repo will automatically generate a link to the rule on [docs.bearer.com](/). (Optional)
- `auxiliary`: Allows you to define helper rules and detectors to make pattern-building more robust. Auxiliary rules contain a unique `id` and their own `patterns` in the same way rules do. You’re unlikely to use this regularly. See the [weak_encryption](https://github.com/Bearer/bearer-rules/blob/main/ruby/lang/weak_encryption.yml) rule for examples. In addition, see our advice on how to avoid [variable joining](#variable-joining) in auxiliary rules. (Optional)
- `skip_data_types`: Allows you to prevent the specified data types from triggering this rule. Takes an array of strings matching the data type names. Example: “Passwords”. (Optional)
- `only_data_types`: Allows you to limit the specified data types that trigger this rule. Takes an array of strings matching the data type names. Example: “Passwords”. (Optional)

## Patterns

Patterns allow rules to look for matches in your code, much like regular expressions, but they take advantage of Bearer CLI’s underlying data type detection capabilities.

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

A new pattern appears! This time, it looks for sensitive data types inside the `Net::FTP.open()` block, using Bearer CLI’s built-in `datatype` detection. To better understand what’s happening, let’s examine variables and filters in more detail.

_Note: in the example above, the third pattern uses a different YAML syntax and the `pattern` key. This is required to define filters for a pattern._

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

`$<!>`: In some instances, a pattern requires some wider context to match the exact line of code where the rule occurs. For those cases, use this variable type to explicitly mark the line for Bearer CLI to highlight in the security report. You’ll mostly need this for rules that target configuration files and settings, rather than logic-related code. For example:

```yaml
patterns:
  - |
    Rails.application.configure do
      $<!>config.force_ssl = false
    end
```

`$<VARNAME:type>`: This is a special type of named variable that helps Bearer CLI's underlying engine by explicitly stating the node type in an AST. This is usually only used in special circumstances, like Ruby's method/block arguments. In this example, we need to mark `CONFIG` as an `identifier`.

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
  - `filename_regex`: Applies a regular expression test against the filename. This uses the [RE2 syntax](https://github.com/google/re2/wiki/Syntax).
  - `values`: Provide an array of values to match a variable against. Useful for specific method names and known options.
  - `length_less_than`: Compare the length of the (string) variable to the number provided with a **less than** statement.
  - `string_regex`: Applies a regular expression test against the string value of the linked variable. This uses the [RE2 syntax](https://github.com/google/re2/wiki/Syntax).
  - `less_than`: Compare the variable to the number provided with a **less than** statement.
  - `less_than_or_equal`: Compare the variable to the number provided with a _less than or equal_ statement.
  - `greater_than`: Compare the variable to the number provided with a _greater than_ statement.
  - `greater_than_or_equal`: Compare the variable to the number provided with a _greater than or equal_ statement.
  - `regex`: Applies a regular expression test against the code content of the linked variable. This uses the [RE2 syntax](https://github.com/google/re2/wiki/Syntax).
- `not`: Inverts the results of another filter. Can be used with a single comparison key by nesting the key below `not`, or with an `either` block by nesting the block below `not`.
- `either`: Allows for multiple conditional checks. It behaves like an OR condition. You can nest any filter inside of `either`, such as `values`, `detection`, etc.
- `detection`: Detection filters rely on existing filter types, so they handle much of the logic for you.
  - `datatype`: This is the detection type you’ll most often see. It uses Bearer CLI's scan to match any data type.
  - `insecure_url`: Useful for instances where you want to prevent unsecured HTTP requests. It explicitly matches `http://`.
  - `<auxiliary-detection-id>`: This allows you to link external and custom detection types by their id. See the `auxiliary` description in the rule config at the top of this page for more details, and the [weak_encryption rule](https://github.com/Bearer/bearer-rules/blob/main/ruby/lang/weak_encryption.yml) for an example.

To better understand how filters and variables interact, see the pattern examples below.

### Pattern examples

Let’s look at some example patterns from Bearer CLI’s core rules that use these filter and variable concepts.

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

Once you’ve written a custom rule, there are a few ways to tell Bearer CLI about it.

Run scans with the `--external-rule-dir` flag.

```bash
bearer scan . --external-rule-dir /path/to/rules/
```

Add the rule to your bearer config file.

```yaml
external-rule-dir: /path/to/rules/
```

_Note: Including an external rules directory adds custom rules to the security report. To only run custom rules, you’ll need to use the `only-rule` flag or configuration setting and pass it the IDs of your custom rule._

## Rule best practices

1. Matching patterns in a rule cause _rule findings_. Depending on the severity level, findings can cause CI to exit and will display in the security report. Keep this in mind when writing patterns so you don’t match a best practice condition and trigger a failed scan.
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

## Variable joining

When a rule relies on another rule as part of its `filter` declaration, the variables are treated as a single set when matching against the code. Variables with the same name will identify as the same AST node. This can cause problems where the wrong AST nodes are used in the detection, leading to unexpected scan results. Let's look at an example.

We'll start with the code we want to target:

```ruby
# We want to match this:
sql_query("SELECT * FROM #{user_input}")

# but not this:
sql_query("SELECT * FROM #{sanitize(user_input)}")
```

The goal is to target the lack of `sanitize` in SQL queries. Now let's look at a rule that can handle this use case. Don't worry, we'll break down each part after the code.

```yml
# Rule 0
patterns:
  - pattern: |
      sql_query($<QUERY>)
    filters:
      - variable: QUERY
        detection: my_rule_user_input
      - not:
          variable: QUERY
          detection: my_rule_sanitized
auxiliary:
  # Rule 1
  - id: my_rule_sanitized
    patterns:
      - pattern: sanitize($<SANITIZED>)
        filters:
          - variable: SANITIZED
            detection: my_rule_user_input
  # Rule 2
  - id: my_rule_user_input_source
    patterns:
      - user_input
  # Rule 3
  - id: my_rule_user_input
    patterns:
      - pattern: $<USER_INPUT>
        filters:
          - variable: USER_INPUT
            detection: my_rule_user_input_source
            contains: false
languages:
  - ruby
severity: medium
metadata:
  description: "Variable joining example"
  remediation_message: ""
  cwe_id:
    - 601
  id: "my_rule"
```

The `patterns` portion at the beginning should look familiar. The only difference compared to most rules is that it references `detection` types that are auxiliary rules. The main pattern looks for `sql_query`, and then uses filters to tell Bearer CLI to apply each detection to any code it finds in `$<QUERY>`. In this case, it wants to trigger the rule using `my_rule_user_input` but NOT `my_rule_sanitized`. I've labeled the core rule/pattern combo as rule 0, then each Aux rule as rules 1, 2, and 3 to make them easier to follow.

Let's start with Rule 0's positive case and follow the detection:

- **Rule 0**'s positive filter calls on `my_rule_user_input`, **Rule 3**, to handle the detection.
- **Rule 3** uses the `my_rule_user_input_source` detection, which belongs to **Rule 2**. _It sets `contains: false`, but we'll come back to that._
- **Rule 2** is a simple pattern that looks for `user_input`. This is from our initial code target example. Think of it as a variable that was passed to `sql_query`.

That chain of detections will result in a match for the non-sanitized code. Now let's look at the the negative, `not` filter case and see if we notice any overlap.

- **Rule 0**'s negative filter calls on `my_rule_sanitized`, which is **Rule 1**.
- **Rule 1** sets up its own pattern to look for the `sanitize` function, then calls on `my_rule_user_input`, **Rule 3**, to handle the detection.
- **Rule 3**, as we saw before, bounces the detection over to **Rule 2** to handle the code portion.
- **Rule 2** uses its simple pattern to detect the `user_input` variable in the code.

This seems like it would all work fine, but because both rules rely on `my_rule_user_input` we could end up in a situation where both occurrences refer to the same AST node. That's where `contains: false` comes into play. It ensures that the `USER_INPUT` from `$<USER_INPUT>` will always match the specific AST node, and not a parent of it. Without it, the sanitize call would exist inside `USER_INPUT` and Rule 1 would never match.

We recommend using unique variable names in rules that refer to each other when you don't intend them to match the same code locations. If your use case doesn't rely on reusing an input source across multiple, differing patterns, it may be easier to duplicate the detection logic with unique names. Otherwise, make use of `contains: false` to prevent the variable joining.

## Shared rules

You can use shared rules to avoid duplication of auxiliary rules between different rule files. To use one rule from another, it must be of type `shared` and must be imported by the rule that uses it.

As shared rules are only used by other rules and do not result in any findings, they have no associated CWE, severity, etc.

### Example

Shared rule:

```yaml
languages:
  - ruby
type: shared
patterns:
  - params[$<_>]
metadata:
  description: "Ruby user input"
  id: ruby_shared_user_input
```

Main rule:

```yaml
languages:
  - ruby
imports:
  - ruby_shared_user_input
patterns:
  - pattern: unsafe($<USER_INPUT>)
    filters:
      - variable: USER_INPUT
        detection: ruby_shared_user_input
severity: high
metadata:
  description: "Unsafe user input detected."
  remediation_message: "..."
  cwe_id:
    - 601
  id: ruby_lang_unsafe_user_input
```

## Syntax updates

### v1.1 Trigger changes

If you have created a custom rule before v1.1 you will need to make the some small changes

#### Local, Present

If you use `trigger: local` or `trigger: present` you can simply remove the trigger attribute and your rule should work as before.

#### Absence

If you use `trigger: absence`, replace it with the following syntax and remove `trigger_rule_on_presence_of` from your existing rule.

```yaml
trigger:
  match_on: absence
  required_detection: # whatever value you had for `trigger_rule_on_presence_of`
```

#### Global

For `trigger: global` replace it with the following syntax.

```yaml
trigger:
  data_types_required: true
```

## Need some help?

If you’re running into any problems or need some help, check out the [Discord Community]({{meta.links.discord}}). You can also [create a new issue]({{meta.links.issues}}) on GitHub.
