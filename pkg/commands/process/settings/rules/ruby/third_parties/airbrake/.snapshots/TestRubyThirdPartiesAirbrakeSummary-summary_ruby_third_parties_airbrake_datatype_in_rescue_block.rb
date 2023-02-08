critical:
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_airbrake
      policy_description: Do not send sensitive data to Airbrake.
      line_number: 4
      filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_rescue_block.rb
      category_groups:
        - PII
      parent_line_number: 4
      parent_content: Airbrake.notify_sync(current_user.first_name)


--

