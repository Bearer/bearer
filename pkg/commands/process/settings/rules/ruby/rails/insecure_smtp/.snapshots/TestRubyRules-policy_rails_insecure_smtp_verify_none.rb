critical:
    - policy_name: ""
      policy_dsrid: DSR-2
      policy_display_id: ruby_rails_insecure_smtp
      policy_description: Only communicate with secure SMTP connections.
      line_number: 8
      filename: pkg/commands/process/settings/rules/ruby/rails/insecure_smtp/testdata/verify_none.rb
      category_groups:
        - PII
      parent_line_number: 8
      parent_content: 'openssl_verify_mode: "none"'


--

