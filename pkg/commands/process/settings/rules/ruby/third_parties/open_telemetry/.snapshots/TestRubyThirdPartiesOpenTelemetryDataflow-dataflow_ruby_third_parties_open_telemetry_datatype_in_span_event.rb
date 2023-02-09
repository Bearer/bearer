data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_event.rb
              line_number: 2
              field_name: email
              object_name: current_user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_event.rb
              line_number: 4
              field_name: email
              object_name: current_user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_open_telemetry
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_event.rb
              line_number: 2
              parent:
                line_number: 2
                content: 'span.add_event("Schedule job for user: #{current_user.email}")'
              field_name: email
              object_name: current_user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_event.rb
              line_number: 4
              parent:
                line_number: 3
                content: |-
                    span.add_event("Cancel job for user", attributes: {
                      "current_user" => current_user.email
                    })
              field_name: email
              object_name: current_user
              subject_name: User
    - detector_id: open_telemetry_current_span
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_span_event.rb
          line_number: 1
          parent:
            line_number: 1
            content: OpenTelemetry::Trace.current_span
          content: |
            OpenTelemetry::Trace.current_span
components: []


--

