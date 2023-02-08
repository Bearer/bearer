critical:
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_open_telemetry
      policy_description: Do not send sensitive data to Open Telemetry.
      line_number: 7
      filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_attributes.rb
      category_groups:
        - PII
      parent_line_number: 5
      parent_content: |-
        current_span.add_attributes({
            "user.id" => user.id,
            "user.first_name" => user.first_name
          })
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_open_telemetry
      policy_description: Do not send sensitive data to Open Telemetry.
      line_number: 13
      filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_attributes.rb
      category_groups:
        - PII
      parent_line_number: 14
      parent_content: current_span.set_attribute("current_users", users)


--

