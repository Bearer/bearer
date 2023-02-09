critical:
    - rule_dsrid: DSR-2
      rule_display_id: ruby_rails_insecure_smtp
      rule_description: Only communicate with secure SMTP connections.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_rails_insecure_smtp
      line_number: 8
      filename: pkg/commands/process/settings/rules/ruby/rails/insecure_smtp/testdata/verify_none.rb
      category_groups:
        - PII
      parent_line_number: 8
      parent_content: 'openssl_verify_mode: "none"'


--

