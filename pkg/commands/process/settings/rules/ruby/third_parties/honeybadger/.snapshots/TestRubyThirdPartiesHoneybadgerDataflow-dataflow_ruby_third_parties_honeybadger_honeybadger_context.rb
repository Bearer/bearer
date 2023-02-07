data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_context.rb
              line_number: 8
              field_name: email
              object_name: current_user
              subject_name: User
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_context.rb
              line_number: 1
              field_name: first_name
              object_name: current_user
              subject_name: User
    - name: Lastname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_context.rb
              line_number: 1
              field_name: last_name
              object_name: current_user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_honeybadger
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_context.rb
              line_number: 8
              parent:
                line_number: 7
                content: |-
                    Honeybadger.context({
                      my_data: current_user.email
                    })
              field_name: email
              object_name: current_user
              subject_name: User
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_context.rb
              line_number: 1
              parent:
                line_number: 3
                content: |-
                    Honeybadger.context({
                      tags: tags
                    })
              field_name: first_name
              object_name: current_user
              subject_name: User
        - name: Lastname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_context.rb
              line_number: 1
              parent:
                line_number: 3
                content: |-
                    Honeybadger.context({
                      tags: tags
                    })
              field_name: last_name
              object_name: current_user
              subject_name: User
components: []


--

