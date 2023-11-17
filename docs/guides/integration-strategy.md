---
title: Strategy on how to integrate Bearer into your workflow
---

# Strategy on how to integrate Bearer into your workflow

Welcome to this guide on integrating Bearer into your development process.

Follow these steps to seamlessly incorporate Bearer at every stage of your workflow. We understand that adopting new practices can take time, so we provide this guide as a strategy blueprint to gradually implement these changes, allowing you and your team to adjust easily.

## Step 1: Understand the current state of your application

Before integrating Bearer into your CI/CD pipeline, it's essential to assess the existing landscape. To do so, follow these steps:

- Install Bearer on your local machine (see [doc](/quickstart)).
- Run the scan command on your repository main branch.

This initial scan will provide insights into your application's current state. Let's dive in!

```bash
$ bearer scan myRepo/
…
MEDIUM: Weak hashing library (MD5) detected. [CWE-327, CWE-328]
https://docs.bearer.com/reference/rules/javascript_lang_weak_hash_md5
To ignore this finding, run: bearer ignore add 42a76a8c10a52b38c1b8729a2f211830_0

File: lib/insecurity.ts:43

 43 export const hash = (data: string) => crypto.createHash('md5').update(data).digest('hex')
=====================================

72 checks, 55 findings

CRITICAL: 0
HIGH: 21 (CWE-22, CWE-312, CWE-79, CWE-798, CWE-89, CWE-918)
MEDIUM: 34 (CWE-327, CWE-328, CWE-525, CWE-548, CWE-73, CWE-79)
LOW: 0
WARNING: 0
```

If your code triggers any findings, as show above, read the next section. Otherwise, skip directly to step 3.

## Step 2: Triage and remediate existing findings

We need to start triaging findings from the CLI on your local machine.

Here is our recommended strategy:

1. Start by addressing "Critical" findings, then "High", "Medium", and finally "Low" findings.
2. Remove any irrelevant findings for each severity level by using the `bearer ignore` command and categorize them as _false positive_ (see [doc](/guides/configure-scan/#ignore-specific-findings)).
3. If a rule is problematic for your codebase, skip it entirely using the `--skip-rule` option (see [doc](/guides/configure-scan/#skip-rules-for-the-entire-scan)).
4. Review the remaining findings with your team and begin fixing them.
5. Commit your changes before running another scan to ensure the scanner picks up the changes.
6. If you have unresolved findings, choose one of two strategies:
   - Postpone handling them for now, but note that they will appear on future scans.
   - Ignore them using the `bearer ignore` command, categorizing them as _allowed_ and providing a comment explaining why.

{% callout "info" %}
If you have many findings and need assistance from your team to triage and remediate them, consider using <a href="/guides/bearer-cloud">Bearer Cloud</a>, a UI interface that complements Bearer CLI for faster resolution.
{% endcallout %}

## Step 3: Minimize new issues from being introduced

To prevent the introduction of new issues in your codebase, it is crucial to identify and address them before developers merge their code into the main branch, as part of your CI.

Start configuring a Bearer scan when developers create a new PR/MR. This allows for immediate feedback and provides the necessary context for any detected findings:

- For GitHub Actions, refer to [using GitHub Action](/guides/github-action/#pull-request-diff).
- For GitLab CI, refer to the [using GitLab Action](/guides/gitlab/#gitlab-merge-request-diff).
- For CircleCI, refer to [using CircleCI](/guides/ci-setup/#circleci).
- For other CI, refer to the [universal setup](/guides/ci-setup/#universal-setup).

Also, depending on your GitHub or GitLab configuration, you can choose to block MR/PR until all findings have been resolved. However, it is important to consider the impact this may have on your development team before implementing this setting.

When integrating Bearer into your CI, we recommend the following configuration:

1. Start by only triggering critical and high findings using the `--severity` option (see [doc](/guides/configure-scan/#limit-severity-levels)).
2. If you find that certain rules are not accurately identifying issues for your stack, you can disable them using the `--skip-rule` option (see [doc](/guides/configure-scan/#skip-rules-for-the-entire-scan)).
3. After a few weeks of using Bearer in your CI, fine-tune these settings and consider including more severity-level findings.

{% callout "info" %}
If you are concerned about displaying a failed status in your CI, you have the option to force a successful exit code regardless of the presence of findings (see <a href="/guides/configure-scan/#force-a-given-exit-code-for-the-scan-command">doc</a>).
{% endcallout %}

## Step 4: Check for critical findings before deploying your code

When integrating Bearer into your CI, it performs differential scans that only analyze the code changes. This approach provides developers with timely feedback, only on their code, and ensures fast scans without unnecessary waiting.

However, some findings may persist. This could be because they haven't been addressed during the previous step (_if you don't enforce PR/PR checks_), or because a complete code scan is required to identify specific findings spanning across unmodified code.

To prevent significant security vulnerabilities in production, we recommend conducting a full scan as part of your CD process or when merging code into the main branch. Since breaking deployments is a serious matter, we advise limiting scanning to certain severity levels, such as critical, or to specific rules using the `--only-rule` option (see [doc](/guides/configure-scan/#run-only-specified-rules)).

## Step 5: Schedule regular complete scans

If you've followed our previous setup recommendations, you may have noticed that certain findings will persist in production. This is particularly true for low-priority findings or rules that you have low confidence in and have disabled.

To address these findings and maintain effective control over your security posture, we recommend conducting a full scan of your codebase at least once a month. This scan should allow all severity levels and enable all rules to ensure a comprehensive review.

Once the scan is complete, you can employ the same logic as step 2) to prioritize and resolve any identified findings.

## Conclusion

As you can see, our recommendations aim to strike a balance between tight security and seamless integration into your development workflow.

Feel free to adapt them to your own experience and specific organization requirements. For instance, you have the option to opt for a more rigid or flexible approach.
