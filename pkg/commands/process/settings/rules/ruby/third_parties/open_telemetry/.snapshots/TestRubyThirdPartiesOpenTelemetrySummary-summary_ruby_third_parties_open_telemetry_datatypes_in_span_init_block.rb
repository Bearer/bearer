critical:
    - rule_dsrid: DSR-1
      rule_display_id: ruby_third_parties_open_telemetry
      rule_description: Do not send sensitive data to Open Telemetry.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_open_telemetry
      line_number: 2
      filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatypes_in_span_init_block.rb
      category_groups:
        - PII
      parent_line_number: 2
      parent_content: |-
        Tracer.in_span("data leaking", attributes: { "current_user" => user.email, "date" => DateTime.now }) do |span|
          puts "in the span block"
        end
    - rule_dsrid: DSR-1
      rule_display_id: ruby_third_parties_open_telemetry
      rule_description: Do not send sensitive data to Open Telemetry.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_open_telemetry
      line_number: 6
      filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatypes_in_span_init_block.rb
      category_groups:
        - PII
      parent_line_number: 6
      parent_content: |-
        SomeOtherTracer.in_span("data leaking", attributes: { "current_user" => user.email, "date" => DateTime.now }) do |span|
          span.add_attributes(user.email)
        end
    - rule_dsrid: DSR-1
      rule_display_id: ruby_third_parties_open_telemetry
      rule_description: Do not send sensitive data to Open Telemetry.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_open_telemetry
      line_number: 7
      filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatypes_in_span_init_block.rb
      category_groups:
        - PII
      parent_line_number: 7
      parent_content: span.add_attributes(user.email)
    - rule_dsrid: DSR-1
      rule_display_id: ruby_third_parties_open_telemetry
      rule_description: Do not send sensitive data to Open Telemetry.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_open_telemetry
      line_number: 11
      filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatypes_in_span_init_block.rb
      category_groups:
        - PII
      parent_line_number: 11
      parent_content: 'span.add_event("leaking data for #{user.email}")'


--

