critical:
    - rule_dsrid: DSR-1
      rule_display_id: javascript_third_parties_open_telemetry
      rule_description: Do not send sensitive data to Open Telemetry.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_open_telemetry
      line_number: 5
      filename: pkg/commands/process/settings/rules/javascript/third_parties/open_telemetry/testdata/datatype_in_add_event.js
      category_groups:
        - PII
      parent_line_number: 4
      parent_content: |-
        currentSpan.addEvent('my-event', {
          'event.metadata': customer.emailAddress
        })


--

