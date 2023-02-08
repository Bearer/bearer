data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_notify.rb
              line_number: 4
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_notify.rb
              line_number: 8
              field_name: email
              object_name: customer
              subject_name: User
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_notify.rb
              line_number: 1
              field_name: first_name
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_airbrake
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_notify.rb
              line_number: 4
              parent:
                line_number: 3
                content: |-
                    Airbrake.notify('App crashed!', {
                      current_user: user.email
                    })
              field_name: email
              object_name: user
              subject_name: User
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_notify.rb
              line_number: 1
              parent:
                line_number: 1
                content: Airbrake.notify(user.first_name)
              field_name: first_name
              object_name: user
              subject_name: User
components: []


--

