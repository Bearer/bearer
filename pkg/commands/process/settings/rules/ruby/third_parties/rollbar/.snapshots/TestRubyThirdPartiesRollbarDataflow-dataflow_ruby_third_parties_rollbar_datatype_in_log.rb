data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log.rb
              line_number: 1
              field_name: email
              object_name: user
              subject_name: User
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log.rb
              line_number: 3
              field_name: first_name
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_rollbar
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'Rollbar.log("error", "oops #{user.email}")'
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log.rb
              line_number: 2
              parent:
                line_number: 2
                content: 'Rollbar.log("error", "oops", user: { email: "someone@example.com" })'
              field_name: email
              object_name: user
              subject_name: User
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log.rb
              line_number: 3
              parent:
                line_number: 3
                content: 'Rollbar.log("error", "oops", { user: { first_name: "someone" } })'
              field_name: first_name
              object_name: user
              subject_name: User
components: []


--

