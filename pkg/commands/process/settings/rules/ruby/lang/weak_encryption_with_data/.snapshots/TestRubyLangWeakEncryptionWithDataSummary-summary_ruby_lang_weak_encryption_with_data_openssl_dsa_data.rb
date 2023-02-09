critical:
    - rule_dsrid: DSR-7
      rule_display_id: ruby_lang_weak_encryption_with_data
      rule_description: Do not use weak encryption libraries to encrypt sensitive data.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_weak_encryption_with_data
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/openssl_dsa_data.rb
      category_groups:
        - PII
      parent_line_number: 3
      parent_content: dsa_encrypt.export(cipher, user.email)
    - rule_dsrid: DSR-7
      rule_display_id: ruby_lang_weak_encryption_with_data
      rule_description: Do not use weak encryption libraries to encrypt sensitive data.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_weak_encryption_with_data
      line_number: 5
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/openssl_dsa_data.rb
      category_groups:
        - PII
      parent_line_number: 5
      parent_content: OpenSSL::PKey::RSA.new(2048).to_pem(cipher, user.first_name)


--

