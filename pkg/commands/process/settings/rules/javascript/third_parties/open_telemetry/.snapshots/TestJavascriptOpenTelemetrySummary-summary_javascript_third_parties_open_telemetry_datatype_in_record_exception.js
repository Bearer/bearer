high:
    - rule_dsrid: DSR-1
      rule_display_id: javascript_third_parties_open_telemetry
      rule_description: Do not send sensitive data to Open Telemetry.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_open_telemetry
      line_number: 9
      filename: pkg/commands/process/settings/rules/javascript/third_parties/open_telemetry/testdata/datatype_in_record_exception.js
      category_groups:
        - Personal Data
      parent_line_number: 9
      parent_content: span.recordException(currentUser.ipAddress)


--

