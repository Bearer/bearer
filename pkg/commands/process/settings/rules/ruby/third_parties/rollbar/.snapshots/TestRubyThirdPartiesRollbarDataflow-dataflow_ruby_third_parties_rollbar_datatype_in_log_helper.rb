data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log_helper.rb
              line_number: 1
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log_helper.rb
              line_number: 2
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log_helper.rb
              line_number: 6
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log_helper.rb
              line_number: 8
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log_helper.rb
              line_number: 10
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log_helper.rb
              line_number: 12
              field_name: email
              object_name: user
              subject_name: User
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log_helper.rb
              line_number: 4
              field_name: first_name
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_rollbar
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log_helper.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'Rollbar.critical("oops #{user.email}")'
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log_helper.rb
              line_number: 2
              parent:
                line_number: 2
                content: 'Rollbar.critical(e, "oops #{user.email}")'
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log_helper.rb
              line_number: 3
              parent:
                line_number: 3
                content: 'Rollbar.critical(e, user: { email: "someone@example.com" })'
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log_helper.rb
              line_number: 6
              parent:
                line_number: 6
                content: 'Rollbar.error("oops #{user.email}")'
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log_helper.rb
              line_number: 8
              parent:
                line_number: 8
                content: 'Rollbar.debug("oops #{user.email}")'
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log_helper.rb
              line_number: 10
              parent:
                line_number: 10
                content: 'Rollbar.info("oops #{user.email}")'
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log_helper.rb
              line_number: 12
              parent:
                line_number: 12
                content: 'Rollbar.warning("oops #{user.email}")'
              field_name: email
              object_name: user
              subject_name: User
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_log_helper.rb
              line_number: 4
              parent:
                line_number: 4
                content: 'Rollbar.critical(e, { user: { first_name: "someone" } })'
              field_name: first_name
              object_name: user
              subject_name: User
components: []


--

