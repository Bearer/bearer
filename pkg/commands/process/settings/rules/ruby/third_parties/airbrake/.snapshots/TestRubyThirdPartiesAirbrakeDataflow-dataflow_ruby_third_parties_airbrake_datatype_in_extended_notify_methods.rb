data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 13
              field_name: email
              object_name: current_user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 34
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 53
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 71
              field_name: email
              object_name: user
              subject_name: User
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 4
              field_name: first_name
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 22
              field_name: first_name
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 43
              field_name: first_name
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_airbrake
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 13
              parent:
                line_number: 11
                content: |-
                    Airbrake.notify_request(
                      {
                        current_user: current_user.email
                      },
                      request_id: 123
                    )
              field_name: email
              object_name: current_user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 34
              parent:
                line_number: 32
                content: |-
                    Airbrake.notify_query(
                      {
                        user: user.email
                      },
                      request_id: 123
                    )
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 53
              parent:
                line_number: 51
                content: |-
                    Airbrake.notify_performance_breakdown(
                      {
                        user: user.email
                      },
                      request_id: 123
                    )
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 71
              parent:
                line_number: 69
                content: |-
                    Airbrake.notify_queue(
                      {
                        user: user.email
                      },
                      job_id: 123
                    )
              field_name: email
              object_name: user
              subject_name: User
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 4
              parent:
                line_number: 2
                content: |-
                    Airbrake.notify_request(
                      method: 'GET',
                      route: "/users/#{user.first_name}",
                      status_code: 200,
                      timing: 123.45 # ms
                    )
              field_name: first_name
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 22
              parent:
                line_number: 20
                content: |-
                    Airbrake.notify_query(
                      method: 'GET',
                      route: "/users/#{user.first_name}",
                      query: 'SELECT * FROM foos',
                      func: 'foo', # optional
                      file: 'foo.rb', # optional
                      line: 123, # optional
                      timing: 123.45 # ms
                    )
              field_name: first_name
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 43
              parent:
                line_number: 41
                content: |-
                    Airbrake.notify_performance_breakdown(
                      method: 'GET',
                      route: "/users/#{user.first_name}",
                      response_type: 'json',
                      groups: { db: 24.0, view: 0.4 }, # ms
                      timing: 123.45 # ms
                    )
              field_name: first_name
              object_name: user
              subject_name: User
components: []


--

