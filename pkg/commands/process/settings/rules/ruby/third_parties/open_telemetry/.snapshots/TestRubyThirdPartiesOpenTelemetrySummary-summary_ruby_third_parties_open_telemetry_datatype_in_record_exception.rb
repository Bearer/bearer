critical:
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_open_telemetry
      policy_description: Do not send sensitive data to Open Telemetry.
      line_number: 7
      filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_record_exception.rb
      category_groups:
        - PII
        - Personal Data
      parent_line_number: 7
      parent_content: 'current_span.status = OpenTelemetry::Trace::Status.error("error for user #{current_user.email}")'
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_open_telemetry
      policy_description: Do not send sensitive data to Open Telemetry.
      line_number: 17
      filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_record_exception.rb
      category_groups:
        - PII
        - Personal Data
      parent_line_number: 17
      parent_content: 'current_span.record_exception(ex, attributes: { "user.ip" => user.ip_address })'


--

