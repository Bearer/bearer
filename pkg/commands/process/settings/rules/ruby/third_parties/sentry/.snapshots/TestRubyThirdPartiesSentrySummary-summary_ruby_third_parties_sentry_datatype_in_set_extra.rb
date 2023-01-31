critical:
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_sentry
      policy_description: Do not send sensitive data to Sentry.
      line_number: 2
      filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_extra.rb
      category_groups:
        - PII
      parent_line_number: 2
      parent_content: scope.set_extra(:email, user.email)
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_sentry
      policy_description: Do not send sensitive data to Sentry.
      line_number: 6
      filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_extra.rb
      category_groups:
        - PII
      parent_line_number: 6
      parent_content: scope.set_extra(:email, user.email)


--

