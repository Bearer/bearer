critical:
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: |-
        Airbrake.notify_request_sync(
          method: 'GET',
          route: "/users/#{user.first_name}",
          status_code: 200,
          timing: 123.45 # ms
        )
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 9
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
      category_groups:
        - PII
      parent_line_number: 7
      parent_content: |-
        Airbrake.notify_request(
          method: 'GET',
          route: "/users/#{user.first_name}",
          status_code: 200,
          timing: 123.45 # ms
        )
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 17
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
      category_groups:
        - PII
      parent_line_number: 15
      parent_content: |-
        Airbrake.notify_request_sync(
          {
            current_user: current_user.email
          },
          request_id: 123
        )
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 23
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
      category_groups:
        - PII
      parent_line_number: 21
      parent_content: |-
        Airbrake.notify_request(
          {
            current_user: current_user.email
          },
          request_id: 123
        )
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 31
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
      category_groups:
        - PII
      parent_line_number: 29
      parent_content: |-
        Airbrake.notify_query_sync(
          method: 'GET',
          route: "/users/#{user.first_name}",
          query: 'SELECT * FROM foos'
        )
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 36
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
      category_groups:
        - PII
      parent_line_number: 34
      parent_content: |-
        Airbrake.notify_query(
          method: 'GET',
          route: "/users/#{user.first_name}",
          query: 'SELECT * FROM foos'
        )
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 43
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
      category_groups:
        - PII
      parent_line_number: 41
      parent_content: |-
        Airbrake.notify_query_sync(
          {
            user: user.email
          },
          request_id: 123
        )
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 49
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
      category_groups:
        - PII
      parent_line_number: 47
      parent_content: |-
        Airbrake.notify_query(
          {
            user: user.email
          },
          request_id: 123
        )
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 57
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
      category_groups:
        - PII
      parent_line_number: 55
      parent_content: |-
        Airbrake.notify_performance_breakdown_sync(
          method: 'GET',
          route: "/users/#{user.first_name}",
          response_type: 'json',
          groups: { db: 24.0, view: 0.4 }, # ms
          timing: 123.45 # ms
        )
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 64
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
      category_groups:
        - PII
      parent_line_number: 62
      parent_content: |-
        Airbrake.notify_performance_breakdown(
          method: 'GET',
          route: "/users/#{user.first_name}",
          response_type: 'json',
          groups: { db: 24.0, view: 0.4 }, # ms
          timing: 123.45 # ms
        )
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 73
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
      category_groups:
        - PII
      parent_line_number: 71
      parent_content: |-
        Airbrake.notify_performance_breakdown_sync(
          {
            user: user.email
          },
          request_id: 123
        )
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 79
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
      category_groups:
        - PII
      parent_line_number: 77
      parent_content: |-
        Airbrake.notify_performance_breakdown(
          {
            user: user.email
          },
          request_id: 123
        )
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 101
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
      category_groups:
        - PII
      parent_line_number: 99
      parent_content: |-
        Airbrake.notify_queue_sync(
          {
            user: user.email
          },
          job_id: 123
        )
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 107
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_extended_notify_methods.rb
      category_groups:
        - PII
      parent_line_number: 105
      parent_content: |-
        Airbrake.notify_queue(
          {
            user: user.email
          },
          job_id: 123
        )


--

