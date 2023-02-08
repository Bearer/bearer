risks:
    - detector_id: ruby_third_parties_rollbar
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_scoped.rb
              line_number: 1
              parent:
                line_number: 3
                content: |-
                    Rollbar.scoped(scope) do
                      call
                    end
              field_name: email
              object_name: person
              subject_name: User
components: []


--

