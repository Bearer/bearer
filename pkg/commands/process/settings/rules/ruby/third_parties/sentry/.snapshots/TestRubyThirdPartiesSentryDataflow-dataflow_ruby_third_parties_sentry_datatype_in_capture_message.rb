data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_capture_message.rb
              line_number: 1
              field_name: email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_capture_message.rb
              line_number: 2
              field_name: email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_capture_message.rb
              line_number: 3
              field_name: email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_capture_message.rb
              line_number: 4
              field_name: email
              object_name: user
risks:
    - detector_id: ruby_third_parties_sentry
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_capture_message.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'Sentry.capture_message("test: #{user.email}")'
              field_name: email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_capture_message.rb
              line_number: 2
              parent:
                line_number: 2
                content: 'Sentry.capture_message("test", extra: { email: user.email })'
              field_name: email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_capture_message.rb
              line_number: 3
              parent:
                line_number: 3
                content: 'Sentry.capture_message("test", tags: { email: user.email })'
              field_name: email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_capture_message.rb
              line_number: 4
              parent:
                line_number: 4
                content: 'Sentry.capture_message("test", user: { email: user.email })'
              field_name: email
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_capture_message.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'Sentry.capture_message("test: #{user.email}")'
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_capture_message.rb
              line_number: 2
              parent:
                line_number: 2
                content: 'Sentry.capture_message("test", extra: { email: user.email })'
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_capture_message.rb
              line_number: 3
              parent:
                line_number: 3
                content: 'Sentry.capture_message("test", tags: { email: user.email })'
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_capture_message.rb
              line_number: 4
              parent:
                line_number: 4
                content: 'Sentry.capture_message("test", user: { email: user.email })'
              object_name: user
components: []


--

