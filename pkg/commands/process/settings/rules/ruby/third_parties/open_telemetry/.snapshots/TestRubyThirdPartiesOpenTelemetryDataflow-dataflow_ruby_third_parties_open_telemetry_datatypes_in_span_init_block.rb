data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatypes_in_span_init_block.rb
              line_number: 2
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatypes_in_span_init_block.rb
              line_number: 6
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatypes_in_span_init_block.rb
              line_number: 7
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatypes_in_span_init_block.rb
              line_number: 11
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_open_telemetry
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatypes_in_span_init_block.rb
              line_number: 2
              parent:
                line_number: 2
                content: |-
                    Tracer.in_span("data leaking", attributes: { "current_user" => user.email, "date" => DateTime.now }) do |span|
                      puts "in the span block"
                    end
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatypes_in_span_init_block.rb
              line_number: 6
              parent:
                line_number: 6
                content: |-
                    SomeOtherTracer.in_span("data leaking", attributes: { "current_user" => user.email, "date" => DateTime.now }) do |span|
                      span.add_attributes(user.email)
                    end
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatypes_in_span_init_block.rb
              line_number: 7
              parent:
                line_number: 7
                content: span.add_attributes(user.email)
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatypes_in_span_init_block.rb
              line_number: 11
              parent:
                line_number: 11
                content: 'span.add_event("leaking data for #{user.email}")'
              field_name: email
              object_name: user
              subject_name: User
    - detector_id: open_telemetry_current_span
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatypes_in_span_init_block.rb
          line_number: 2
          parent:
            line_number: 2
            content: span
          content: |
            $<_>.in_span() { |$<!>$<_:identifier>| }
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatypes_in_span_init_block.rb
          line_number: 6
          parent:
            line_number: 6
            content: span
          content: |
            $<_>.in_span() { |$<!>$<_:identifier>| }
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/open_telemetry/testdata/datatypes_in_span_init_block.rb
          line_number: 10
          parent:
            line_number: 10
            content: span
          content: |
            $<_>.in_span() { |$<!>$<_:identifier>| }
components: []


--

