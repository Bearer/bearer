high:
    - policy_name: ""
      policy_dsrid: DSR-8
      policy_display_id: ruby_rails_password_length
      policy_description: Enforce stronger password requirements.
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/rails/password_length/testdata/password_too_short.rb
      category_groups:
        - PII
      parent_line_number: 3
      parent_content: 'validates :password, length: { minimum: 6 }'


--

