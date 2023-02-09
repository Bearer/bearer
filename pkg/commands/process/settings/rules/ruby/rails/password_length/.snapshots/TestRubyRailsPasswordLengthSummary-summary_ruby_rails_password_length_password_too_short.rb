high:
    - rule_dsrid: DSR-8
      rule_display_id: ruby_rails_password_length
      rule_description: Enforce stronger password requirements.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_rails_password_length
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/rails/password_length/testdata/password_too_short.rb
      category_groups:
        - PII
      parent_line_number: 3
      parent_content: 'validates :password, length: { minimum: 6 }'


--

