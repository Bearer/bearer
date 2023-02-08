data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 17
              field_name: email
              object_name: current_user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 23
              field_name: email
              object_name: current_user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 43
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 49
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 73
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 79
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 101
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 107
              field_name: email
              object_name: user
              subject_name: User
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 3
              field_name: first_name
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 9
              field_name: first_name
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 31
              field_name: first_name
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 36
              field_name: first_name
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 57
              field_name: first_name
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 64
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
              line_number: 17
              parent:
                line_number: 15
                content: |-
                    Airbrake.notify_request_sync(
                      {
                        current_user: current_user.email
                      },
                      request_id: 123
                    )
              field_name: email
              object_name: current_user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 23
              parent:
                line_number: 21
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
              line_number: 43
              parent:
                line_number: 41
                content: |-
                    Airbrake.notify_query_sync(
                      {
                        user: user.email
                      },
                      request_id: 123
                    )
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 49
              parent:
                line_number: 47
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
              line_number: 73
              parent:
                line_number: 71
                content: |-
                    Airbrake.notify_performance_breakdown_sync(
                      {
                        user: user.email
                      },
                      request_id: 123
                    )
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 79
              parent:
                line_number: 77
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
              line_number: 101
              parent:
                line_number: 99
                content: |-
                    Airbrake.notify_queue_sync(
                      {
                        user: user.email
                      },
                      job_id: 123
                    )
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 107
              parent:
                line_number: 105
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
              line_number: 3
              parent:
                line_number: 1
                content: |-
                    Airbrake.notify_request_sync(
                      method: 'GET',
                      route: "/users/#{user.first_name}",
                      status_code: 200,
                      timing: 123.45 # ms
                    )
              field_name: first_name
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 9
              parent:
                line_number: 7
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
              line_number: 31
              parent:
                line_number: 29
                content: |-
                    Airbrake.notify_query_sync(
                      method: 'GET',
                      route: "/users/#{user.first_name}",
                      query: 'SELECT * FROM foos'
                    )
              field_name: first_name
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 36
              parent:
                line_number: 34
                content: |-
                    Airbrake.notify_query(
                      method: 'GET',
                      route: "/users/#{user.first_name}",
                      query: 'SELECT * FROM foos'
                    )
              field_name: first_name
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 57
              parent:
                line_number: 55
                content: |-
                    Airbrake.notify_performance_breakdown_sync(
                      method: 'GET',
                      route: "/users/#{user.first_name}",
                      response_type: 'json',
                      groups: { db: 24.0, view: 0.4 }, # ms
                      timing: 123.45 # ms
                    )
              field_name: first_name
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
              line_number: 64
              parent:
                line_number: 62
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

