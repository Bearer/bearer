critical:
    - rule_dsrid: DSR-7
      rule_display_id: ruby_lang_weak_encryption_with_data
      rule_description: Do not use weak encryption libraries to encrypt sensitive data.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_weak_encryption_with_data
      line_number: 2
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: |-
        Crypt::Blowfish.new("insecure").encrypt_block { |user|
          user.password
        }
    - rule_dsrid: DSR-7
      rule_display_id: ruby_lang_weak_encryption_with_data
      rule_description: Do not use weak encryption libraries to encrypt sensitive data.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_weak_encryption_with_data
      line_number: 6
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
      category_groups:
        - PII
      parent_line_number: 5
      parent_content: |-
        Crypt::Blowfish.new("insecure").encrypt_block do |user|
          user.password
        end
    - rule_dsrid: DSR-7
      rule_display_id: ruby_lang_weak_encryption_with_data
      rule_description: Do not use weak encryption libraries to encrypt sensitive data.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_weak_encryption_with_data
      line_number: 9
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
      category_groups:
        - PII
      parent_line_number: 9
      parent_content: Crypt::Blowfish.new("your-key").encrypt_string(user.email)


--

