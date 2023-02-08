data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/datatype_in_tags.rb
              line_number: 2
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/datatype_in_tags.rb
              line_number: 7
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/datatype_in_tags.rb
              line_number: 9
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/datatype_in_tags.rb
              line_number: 10
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/datatype_in_tags.rb
              line_number: 12
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/datatype_in_tags.rb
              line_number: 17
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_datadog
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/datatype_in_tags.rb
              line_number: 2
              parent:
                line_number: 3
                content: c.tags = user
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/datatype_in_tags.rb
              line_number: 7
              parent:
                line_number: 7
                content: span.set_tag('user.email', user.email)
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/datatype_in_tags.rb
              line_number: 9
              parent:
                line_number: 9
                content: Datadog::Tracing.active_span&.set_tag('customer.id', user.email)
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/datatype_in_tags.rb
              line_number: 10
              parent:
                line_number: 10
                content: Datadog::Tracing.active_span.set_tag('customer.id', user.email)
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/datatype_in_tags.rb
              line_number: 12
              parent:
                line_number: 12
                content: |-
                    Datadog::Tracing.trace("web.request", tags: { email: user.email }) do |span, trace|
                      call
                    end
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/datatype_in_tags.rb
              line_number: 17
              parent:
                line_number: 17
                content: span.set_tag('user.email', user.email)
              field_name: email
              object_name: user
              subject_name: User
    - detector_id: ruby_third_parties_datadog_span
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/datatype_in_tags.rb
          line_number: 6
          parent:
            line_number: 6
            content: Datadog.configuration[:cucumber][:tracer].active_span
          content: |
            Datadog.configuration[$<_>][$<_>].active_span
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/datatype_in_tags.rb
          line_number: 12
          parent:
            line_number: 12
            content: span
          content: |
            Datadog::Tracing.trace() { |$<!>$<SPAN:identifier>$<...>| }
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/datatype_in_tags.rb
          line_number: 16
          parent:
            line_number: 16
            content: span
          content: |
            Datadog::Tracing.trace() { |$<!>$<SPAN:identifier>$<...>| }
components: []


--

