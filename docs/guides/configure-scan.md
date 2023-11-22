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

## Only report new findings on a branch

{% callout %}
Differential scanning avoids scanning your entire codebase and drastically
reduces scan times. We recommended that you use this feature when possible.
{% endcallout %}

When scanning a Git repository, you can choose to only report new findings that
have been introduced, relative to a base branch. Any findings that already
existed in the base branch will not be reported.

Use the `--diff` flag to enable differential scanning. The repository's default
branch will be used as the base branch for comparison. You can override this by
setting the `BEARER_DIFF_BASE_BRANCH` environment variable.

```bash
git checkout my-feature
BEARER_DIFF_BASE_BRANCH=base-branch bearer scan --diff .
```

If the base branch is not available in the git repository, it's head will be
fetched by Bearer CLI (a shallow fetch of depth 1).

See our [guide to using the GitHub action](/guides/github-action/#pull-request-diff) and
[guide to using GitLab](/guides/gitlab/#gitlab-merge-request-diff) for
information on using this feature with those services.

## Ignore specific findings

Every finding is associated with a unique fingerprint visible directly in the CLI output, for example:

```bash
HIGH: SQL injection vulnerability detected. [CWE-89]
https://docs.bearer.com/reference/rules/javascript_lang_sql_injection
To ignore this finding, run: bearer ignore add 4b0883d52334dfd9a4acce2fcf810121_0
...
```

If a finding is not relevant, you can ignore it automatically from future scans using the `bearer ignore add` command. This adds the finding's fingerprint to your ignore file. You can also provide optional author information or a comment:

```bash
bearer ignore add 4b0883d52334dfd9a4acce2fcf810121_0 \
  --author="Mish Bear" \
  --comment="Ignore this finding"
  --false-positive
```

<br/>
{% callout "info" %} If you're looking for more options when it comes to managing findings, take a look at <a href="/guides/bearer-cloud">Bearer Cloud</a>. For ignored findings in particular, see <a href="/guides/bearer-cloud/#ignored-findings-in-bearer-cloud">Ignored findings in Bearer Cloud</a>. {% endcallout %}

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

## Limit severity levels

Depending on how you're using Bearer CLI, you may want to limit the severity levels that show up in the report. This can be useful for triaging only the most critical issues. Use the `--severity` flag to define which levels to include from the list of critical, high, medium, low, and warning.

```bash
bearer scan . --severity critical,high
```

## Force a given exit code for the scan command

If you want to force a successful exit code even when findings are reported, use the `--exit-code` flag and set it to 0. It's particularly useful if you want to perform a scan and report findings without failing your CI or CD pipeline.

```bash
bearer scan . --exit-code 0
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

## Generate a SARIF report

Bearer CLI offers SARIF output for tools that make use of the standard. To generate a security report in SARIF and write it to disk, use the `--format` and `--output` flags.

```bash
bearer scan . --format sarif --output sarif-report.sarif
```

If you're using GitHub or GitLab, you can use our [integrations](/guides/ci-setup/) to send SARIF reports directly to those services.

## Format a report as HTML

Sometimes it's useful to have a nicely formatted HTML file to hand off to others. Security and privacy reports support the `html` format type. Pair the `--format` and `--output` flags to create and write an HTML file. It looks like this:

![Preview of html output](/assets/img/bearer-output-html.png)

Run the commands together, replacing the scan location and the output path to match your needs:

```bash
bearer scan . --format html --output path/to/security-scan.html
```

## Send report to Bearer Cloud

If you're looking to manage product and application code security at scale, [Bearer Cloud](https://www.bearer.com/bearer-cloud) offers a platform for teams that syncs with Bearer CLI's output.

Learn how to [send your report](/guides/bearer-cloud) to Bearer Cloud.

![Cloud dashboard](/assets/img/cloud-dashboard.jpg)

## Next steps

For more ways to make the most of our Bearer CLI, check out the [commands reference](/reference/commands/). Need additional help? [Open an issue]({{meta.links.issues}}) or join our [Discord community]({{meta.links.discord}}).
