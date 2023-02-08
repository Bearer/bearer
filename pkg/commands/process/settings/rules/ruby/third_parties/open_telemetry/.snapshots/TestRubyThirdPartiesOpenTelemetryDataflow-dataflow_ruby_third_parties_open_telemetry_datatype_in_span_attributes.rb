data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_attributes.rb
              line_number: 12
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_attributes.rb
              line_number: 18
              field_name: email
              object_name: user
              subject_name: User
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_attributes.rb
              line_number: 7
              field_name: first_name
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_open_telemetry
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_attributes.rb
              line_number: 12
              parent:
                line_number: 12
                content: |-
                    Tracer.in_span("data leaking", attributes: { "current_user" => user.email, "date" => DateTime.now }) do |span|
                      puts "in the span block"
                    end
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_attributes.rb
              line_number: 18
              parent:
                line_number: 19
                content: current_span.set_attribute("current_users", users)
              field_name: email
              object_name: admin_user
              subject_name: User
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_attributes.rb
              line_number: 7
              parent:
                line_number: 5
                content: |-
                    current_span.add_attributes({
                        "user.id" => user.id,
                        "user.first_name" => user.first_name
                      })
              field_name: first_name
              object_name: user
              subject_name: User
    - detector_id: open_telemetry_current_span
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_attributes.rb
          line_number: 3
          parent:
            line_number: 3
            content: OpenTelemetry::Trace.current_span
          content: |
            OpenTelemetry::Trace.current_span
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_attributes.rb
          line_number: 17
          parent:
            line_number: 17
            content: OpenTelemetry::Trace.current_span
          content: |
            OpenTelemetry::Trace.current_span
components: []


--

