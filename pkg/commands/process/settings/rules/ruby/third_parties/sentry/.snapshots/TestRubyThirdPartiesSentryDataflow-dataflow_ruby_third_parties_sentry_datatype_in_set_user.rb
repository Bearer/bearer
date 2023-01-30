data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_user.rb
              line_number: 3
              field_name: email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_user.rb
              line_number: 6
              field_name: email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_user.rb
              line_number: 10
              field_name: email
              object_name: user
risks:
    - detector_id: ruby_third_parties_sentry
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_user.rb
              line_number: 3
              parent:
                line_number: 3
                content: 'Sentry.set_user(email: user.email)'
              field_name: email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_user.rb
              line_number: 6
              parent:
                line_number: 6
                content: 'scope.set_user(email: user.email)'
              field_name: email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_user.rb
              line_number: 10
              parent:
                line_number: 10
                content: 'scope.set_user(email: user.email)'
              field_name: email
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_user.rb
              line_number: 3
              parent:
                line_number: 3
                content: 'Sentry.set_user(email: user.email)'
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_user.rb
              line_number: 6
              parent:
                line_number: 6
                content: 'scope.set_user(email: user.email)'
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_user.rb
              line_number: 10
              parent:
                line_number: 10
                content: 'scope.set_user(email: user.email)'
              object_name: user
components: []


--

