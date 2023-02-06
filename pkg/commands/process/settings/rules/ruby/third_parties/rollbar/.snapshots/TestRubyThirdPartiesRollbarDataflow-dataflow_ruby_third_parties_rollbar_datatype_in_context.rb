risks:
    - detector_id: ruby_third_parties_rollbar
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_context.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'exception.rollbar_context = { user: { email: "someone@example.com" } }'
              field_name: email
              object_name: user
              subject_name: User
components: []


--

