data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_record_exception.rb
              line_number: 7
              field_name: email
              object_name: current_user
              subject_name: User
    - name: IP address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_record_exception.rb
              line_number: 17
              field_name: ip_address
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_open_telemetry
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_record_exception.rb
              line_number: 7
              parent:
                line_number: 7
                content: 'current_span.status = OpenTelemetry::Trace::Status.error("error for user #{current_user.email}")'
              field_name: email
              object_name: current_user
              subject_name: User
        - name: IP address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_record_exception.rb
              line_number: 17
              parent:
                line_number: 17
                content: 'current_span.record_exception(ex, attributes: { "user.ip" => user.ip_address })'
              field_name: ip_address
              object_name: user
              subject_name: User
    - detector_id: open_telemetry_current_span
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_record_exception.rb
          line_number: 2
          parent:
            line_number: 2
            content: OpenTelemetry::Trace.current_span
          content: |
            OpenTelemetry::Trace.current_span
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatype_in_record_exception.rb
          line_number: 11
          parent:
            line_number: 11
            content: OpenTelemetry::Trace.current_span
          content: |
            OpenTelemetry::Trace.current_span
components: []


--

