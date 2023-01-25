critical:
    - policy_name: ""
      policy_dsrid: DSR-3
      policy_display_id: ruby_rails_session
      policy_description: Do not store sensitive data in session cookies.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/rails/session/testdata/datatype_in_session.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: session[:current_user] = user.email


--

