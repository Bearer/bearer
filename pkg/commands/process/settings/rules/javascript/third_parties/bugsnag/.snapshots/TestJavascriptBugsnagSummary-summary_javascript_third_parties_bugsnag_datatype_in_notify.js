high:
    - rule_dsrid: DSR-1
      rule_display_id: javascript_third_parties_bugsnag
      rule_description: Do not send sensitive data to Bugsnag.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_bugsnag
      line_number: 5
      filename: pkg/commands/process/settings/rules/javascript/third_parties/bugsnag/testdata/datatype_in_notify.js
      category_groups:
        - Personal Data
      parent_line_number: 5
      parent_content: 'Bugsnag.notify(user.ip_address + " : " + e)'


--

