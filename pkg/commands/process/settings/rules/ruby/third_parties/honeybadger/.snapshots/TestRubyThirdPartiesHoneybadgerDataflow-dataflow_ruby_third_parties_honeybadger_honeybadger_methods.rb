data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_methods.rb
              line_number: 3
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_honeybadger
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_methods.rb
              line_number: 3
              parent:
                line_number: 2
                content: |-
                    def to_honeybadger_context
                        { user: { id: id, email: email } }
                      end
              field_name: email
              object_name: user
              subject_name: User
components: []


--

