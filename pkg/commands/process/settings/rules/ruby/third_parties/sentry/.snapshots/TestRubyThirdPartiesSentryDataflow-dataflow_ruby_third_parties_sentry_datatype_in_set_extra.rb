data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_extra.rb
              line_number: 2
              field_name: email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_extra.rb
              line_number: 6
              field_name: email
              object_name: user
risks:
    - detector_id: ruby_third_parties_sentry
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_extra.rb
              line_number: 2
              parent:
                line_number: 2
                content: scope.set_extra(:email, user.email)
              field_name: email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_extra.rb
              line_number: 6
              parent:
                line_number: 6
                content: scope.set_extra(:email, user.email)
              field_name: email
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_extra.rb
              line_number: 2
              parent:
                line_number: 2
                content: scope.set_extra(:email, user.email)
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_extra.rb
              line_number: 6
              parent:
                line_number: 6
                content: scope.set_extra(:email, user.email)
              object_name: user
components: []


--

