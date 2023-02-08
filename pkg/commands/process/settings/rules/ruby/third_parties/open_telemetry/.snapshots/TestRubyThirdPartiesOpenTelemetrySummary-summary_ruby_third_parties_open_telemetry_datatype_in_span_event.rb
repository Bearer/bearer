critical:
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_open_telemetry
      policy_description: Do not send sensitive data to Open Telemetry.
      line_number: 2
      filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_event.rb
      category_groups:
        - PII
      parent_line_number: 2
      parent_content: 'span.add_event("Schedule job for user: #{current_user.email}")'
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_open_telemetry
      policy_description: Do not send sensitive data to Open Telemetry.
      line_number: 4
      filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_event.rb
      category_groups:
        - PII
      parent_line_number: 3
      parent_content: |-
        span.add_event("Cancel job for user", attributes: {
          "current_user" => current_user.email
        })


--

