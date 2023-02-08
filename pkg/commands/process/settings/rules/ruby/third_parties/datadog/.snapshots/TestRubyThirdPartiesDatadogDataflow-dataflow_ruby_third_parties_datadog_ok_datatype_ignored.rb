data_types:
    - name: Unique Identifier
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/ok_datatype_ignored.rb
              line_number: 2
              field_name: user_id
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_datadog
      data_types:
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/ok_datatype_ignored.rb
              line_number: 2
              parent:
                line_number: 3
                content: c.tags = user
              field_name: user_id
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/ok_datatype_ignored.rb
              line_number: 6
              parent:
                line_number: 6
                content: |-
                    Datadog::Tracing.trace("web.request", tags: { user_id: 42 }) do |span, trace|
                      call
                    end
              field_name: user_id
              object_name: tags
    - detector_id: ruby_third_parties_datadog_span
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/ok_datatype_ignored.rb
          line_number: 6
          parent:
            line_number: 6
            content: span
          content: |
            Datadog::Tracing.trace() { |$<!>$<SPAN:identifier>$<...>| }
components: []


--

