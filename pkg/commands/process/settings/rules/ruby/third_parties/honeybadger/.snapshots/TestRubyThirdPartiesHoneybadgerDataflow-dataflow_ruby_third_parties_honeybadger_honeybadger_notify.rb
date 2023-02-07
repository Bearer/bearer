data_types:
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_notify.rb
              line_number: 9
              field_name: first_name
              object_name: current_user
              subject_name: User
    - name: Fullname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_notify.rb
              line_number: 2
              field_name: name
              object_name: user
              subject_name: User
    - name: Gender
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_notify.rb
              line_number: 29
              field_name: gender
              object_name: user
              subject_name: User
    - name: Lastname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_notify.rb
              line_number: 9
              field_name: last_name
              object_name: current_user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_honeybadger
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_notify.rb
              line_number: 22
              parent:
                line_number: 28
                content: |-
                    Honeybadger.notify(
                      "Something is wrong here for " + user.gender,
                      class_name: "MyError",
                      error_message: error_message,
                      tags: tags,
                      context: context,
                      parameters: parameters,
                    )
              field_name: email
              object_name: user
              subject_name: User
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_notify.rb
              line_number: 9
              parent:
                line_number: 28
                content: |-
                    Honeybadger.notify(
                      "Something is wrong here for " + user.gender,
                      class_name: "MyError",
                      error_message: error_message,
                      tags: tags,
                      context: context,
                      parameters: parameters,
                    )
              field_name: first_name
              object_name: current_user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_notify.rb
              line_number: 13
              parent:
                line_number: 28
                content: |-
                    Honeybadger.notify(
                      "Something is wrong here for " + user.gender,
                      class_name: "MyError",
                      error_message: error_message,
                      tags: tags,
                      context: context,
                      parameters: parameters,
                    )
              field_name: first_name
              object_name: user
              subject_name: User
        - name: Fullname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_notify.rb
              line_number: 2
              parent:
                line_number: 28
                content: |-
                    Honeybadger.notify(
                      "Something is wrong here for " + user.gender,
                      class_name: "MyError",
                      error_message: error_message,
                      tags: tags,
                      context: context,
                      parameters: parameters,
                    )
              field_name: name
              object_name: user
              subject_name: User
        - name: Gender
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_notify.rb
              line_number: 29
              parent:
                line_number: 28
                content: |-
                    Honeybadger.notify(
                      "Something is wrong here for " + user.gender,
                      class_name: "MyError",
                      error_message: error_message,
                      tags: tags,
                      context: context,
                      parameters: parameters,
                    )
              field_name: gender
              object_name: user
              subject_name: User
        - name: Lastname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_notify.rb
              line_number: 9
              parent:
                line_number: 28
                content: |-
                    Honeybadger.notify(
                      "Something is wrong here for " + user.gender,
                      class_name: "MyError",
                      error_message: error_message,
                      tags: tags,
                      context: context,
                      parameters: parameters,
                    )
              field_name: last_name
              object_name: current_user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_notify.rb
              line_number: 14
              parent:
                line_number: 28
                content: |-
                    Honeybadger.notify(
                      "Something is wrong here for " + user.gender,
                      class_name: "MyError",
                      error_message: error_message,
                      tags: tags,
                      context: context,
                      parameters: parameters,
                    )
              field_name: last_name
              object_name: user
              subject_name: User
components: []


--

