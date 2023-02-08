critical:
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_rollbar
      policy_description: Do not send sensitive data to Rollbar.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_scope.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: 'Rollbar.scope!({ user: { email: "someone@example.com" }})'
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_rollbar
      policy_description: Do not send sensitive data to Rollbar.
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_scope.rb
      category_groups:
        - PII
      parent_line_number: 5
      parent_content: Rollbar.scope(user)
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_rollbar
      policy_description: Do not send sensitive data to Rollbar.
      line_number: 7
      filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_scope.rb
      category_groups:
        - PII
      parent_line_number: 7
      parent_content: 'notifier.scope(user: { first_name: "someone" })'


--

