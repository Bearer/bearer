critical:
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_notify.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: Airbrake.notify(user.first_name)
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 4
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_notify.rb
      category_groups:
        - PII
      parent_line_number: 3
      parent_content: |-
        Airbrake.notify('App crashed!', {
          current_user: user.email
        })
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 8
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_notify.rb
      category_groups:
        - PII
      parent_line_number: 7
      parent_content: |-
        Airbrake.notify('App crashed') do |notice|
          notice[:params][:email] = customer.email
        end


--

