data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_tags.rb
              line_number: 1
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_tags.rb
              line_number: 4
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_tags.rb
              line_number: 8
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_sentry
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_tags.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'Sentry.set_tags(email: user.email)'
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_tags.rb
              line_number: 4
              parent:
                line_number: 4
                content: 'scope.set_tags(email: user.email)'
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_tags.rb
              line_number: 8
              parent:
                line_number: 8
                content: 'scope.set_tags(email: user.email)'
              field_name: email
              object_name: user
              subject_name: User
components: []


--

