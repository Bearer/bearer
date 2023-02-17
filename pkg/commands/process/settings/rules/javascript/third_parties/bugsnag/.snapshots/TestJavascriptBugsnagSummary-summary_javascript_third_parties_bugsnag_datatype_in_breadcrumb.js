critical:
    - rule_dsrid: DSR-1
      rule_display_id: javascript_third_parties_bugsnag
      rule_description: Do not send sensitive data to Bugsnag.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_bugsnag
      line_number: 1
      filename: pkg/commands/process/settings/rules/javascript/third_parties/bugsnag/testdata/datatype_in_breadcrumb.js
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: Bugsnag.leaveBreadcrumb(user.email)


--

