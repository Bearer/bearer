data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_scope.rb
              line_number: 1
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_scope.rb
              line_number: 3
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_rollbar
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_scope.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'Rollbar.scope!({ user: { email: "someone@example.com" }})'
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_scope.rb
              line_number: 3
              parent:
                line_number: 5
                content: Rollbar.scope(user)
              field_name: email
              object_name: user
              subject_name: User
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_scope.rb
              line_number: 7
              parent:
                line_number: 7
                content: 'notifier.scope(user: { first_name: "someone" })'
              field_name: first_name
              object_name: user
              subject_name: User
    - detector_id: ruby_third_parties_rollbar_scope
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_scope.rb
          line_number: 5
          parent:
            line_number: 5
            content: Rollbar.scope(user)
          content: |
            Rollbar.scope()
components: []


--

