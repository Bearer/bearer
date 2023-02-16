high:
    - rule_dsrid: DSR-1
      rule_display_id: javascript_third_parties_new_relic
      rule_description: Do not send sensitive data to New Relic.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_new_relic
      line_number: 7
      filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_notice_error.js
      category_groups:
        - Personal Data
      parent_line_number: 7
      parent_content: newrelic.noticeError(err, customer.ip_address)


--

