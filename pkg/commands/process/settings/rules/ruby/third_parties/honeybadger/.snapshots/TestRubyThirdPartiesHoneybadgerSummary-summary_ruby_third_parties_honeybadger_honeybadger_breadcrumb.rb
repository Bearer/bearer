critical:
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_honeybadger
      policy_description: Do not send sensitive data to Honeybadger.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_breadcrumb.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: 'Honeybadger.add_breadcrumb("Email Sent", metadata: { user_email: current_user.email, message: message })'


--

