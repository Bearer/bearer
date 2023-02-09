critical:
    - rule_dsrid: DSR-3
      rule_display_id: ruby_rails_session
      rule_description: Do not store sensitive data in session cookies.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_rails_session
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/rails/session/testdata/datatype_in_session.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: session[:current_user] = user.email


--

