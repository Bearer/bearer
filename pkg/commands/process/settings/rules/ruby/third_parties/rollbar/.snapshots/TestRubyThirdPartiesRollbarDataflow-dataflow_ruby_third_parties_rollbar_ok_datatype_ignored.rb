risks:
    - detector_id: ruby_third_parties_rollbar
      data_types:
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/ok_datatype_ignored.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'exception.rollbar_context = { user: { user_id: 123 } }'
              field_name: user_id
              object_name: user
              subject_name: User
components: []


--

