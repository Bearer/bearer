---
title: Configure the scan command
---

# Configure the scan to meet your needs

Bearer CLI offers a variety of ways to configure the core `scan` command to best meet your needs. Here are some common situations. For a full list of options, see the [commands reference](/reference/commands/). For many of the command flags listed below, you can also define them in your `bearer.yml` [config file](/reference/config/).

## Select a report type

There are a variety of [report types](/explanations/reports/) to choose from. Bearer CLI defaults to the Security report, but you can select any other type with the `--report` flag.

```bash
bearer scan . --report privacy
```

## Select a scanner type

Did you know that Bearer CLI can also detect hard-coded secrets in your code? In addition to the default SAST scanner, there's a built-in secrets scanner. Use the `--scanner` flag to change [scanner types](/explanations/scanners/).

```bash
bearer scan . --scanner secrets
```

## Skip or ignore specific rules

Sometimes you want to ignore one or more rules, either for the entire scan or for individual blocks of code. Rules are identified by their id, for example: `ruby_lang_exception`.

### Skip rules for the entire scan

To ignore rules for the entire scan you can use the `--skip-rule` flag with the `scan` command.

Using `--skip-rule`:

```bash
# skip a single rule
bearer scan . --skip-rule ruby_lang_exception

# skip multiple rules
bearer scan . --skip-rule ruby_lang_exception,ruby_lang_cookies
```

Using `bearer.yml`

```yaml
rule:
  skip-rule: [ruby_lang_exception, ruby_lang_cookies]
```

### Skip rules for individual code blocks

Bearer CLI supports comment-based rule skipping using the `bearer:disable` comment. To ignore a block of code, place the comment immediately before the block.

In ruby:

```ruby
# bearer:disable ruby_lang_logger, ruby_lang_http_insecure
Net::HTTP.start("http://my.api.com/users/search) do
  logger.warn("Searching for #{current_user.email}")
  ...
end
```

In javascript:

```javascript
// bearer:disable javascript_lang_logger
function logUser(user) {
  log.info(user.name)
}
```

To ignore an individual line of code, place the comment immediately before the line.

```ruby
def my_func
 # bearer:disable ruby_rails_logger
  Rails.logger(current_user.email)
end
```

```javascript
function logUser(user) {
  log.info(user.name)
  // bearer:disable javascript_lang_logger
  log.info(user.uuid)
}
```

## Run only specified rules

Similar to how you can skip rules, you can also tell the scan to only run specific rules. To do so, specify the rule IDs with the `--only-rule` flag.

```bash
bearer scan . --only-rule ruby_lang_cookies
```

## Change the output format

Each [report type](/explanations/reports/) has a default output format, but in general you're able to also select between `json` and `yaml` with the `--format` flag.

```bash
bearer scan . --format yaml
```

## Output to a file

Sometimes you'll want to hand off the report, and while you could pipe the results to another command, we've included the `--output` flag to make it easier. Specify the path to the output file.

```bash
bearer scan . --report dataflow --output dataflow.json
```

## Limit severity levels

Depending on how you're using Bearer CLI, you may want to limit the severity levels that show up in the report. This can be useful for triaging only the most critical issues. Use the `--severity` flag to define which levels to include from the list of critical, high, medium, low, and warning.

```bash
bearer scan . --severity critical,high
```