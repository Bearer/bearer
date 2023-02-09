critical:
    - rule_dsrid: DSR-7
      rule_display_id: ruby_lang_weak_encryption_with_data
      rule_description: Do not use weak encryption libraries to encrypt sensitive data.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_weak_encryption_with_data
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/rc4_data.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: RC4.new("insecure").encrypt(user.password)
    - rule_dsrid: DSR-7
      rule_display_id: ruby_lang_weak_encryption_with_data
      rule_description: Do not use weak encryption libraries to encrypt sensitive data.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_weak_encryption_with_data
      line_number: 4
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/rc4_data.rb
      category_groups:
        - PII
      parent_line_number: 4
      parent_content: rc4_encrypt.encrypt!(user.password)


--

