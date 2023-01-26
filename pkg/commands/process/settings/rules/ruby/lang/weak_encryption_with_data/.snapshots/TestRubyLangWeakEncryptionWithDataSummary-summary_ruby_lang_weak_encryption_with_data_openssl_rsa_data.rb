critical:
    - policy_name: ""
      policy_dsrid: DSR-7
      policy_display_id: ruby_lang_weak_encryption_with_data
      policy_description: Do not use weak encryption libraries to encrypt sensitive data.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/openssl_rsa_data.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: OpenSSL::PKey::RSA.new(File.read('rsa.pem')).private_encrypt(user.password)
    - policy_name: ""
      policy_dsrid: DSR-7
      policy_display_id: ruby_lang_weak_encryption_with_data
      policy_description: Do not use weak encryption libraries to encrypt sensitive data.
      line_number: 5
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/openssl_rsa_data.rb
      category_groups:
        - PII
      parent_line_number: 5
      parent_content: rsa_encrypt.export(cipher, user.password)
    - policy_name: ""
      policy_dsrid: DSR-7
      policy_display_id: ruby_lang_weak_encryption_with_data
      policy_description: Do not use weak encryption libraries to encrypt sensitive data.
      line_number: 7
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/openssl_rsa_data.rb
      category_groups:
        - PII
      parent_line_number: 7
      parent_content: OpenSSL::PKey::RSA.new(2048).to_pem(cipher, user.first_name)
low:
    - policy_name: ""
      policy_dsrid: DSR-7
      policy_display_id: ruby_lang_weak_encryption
      policy_description: Avoid weak encryption libraries.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/openssl_rsa_data.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: OpenSSL::PKey::RSA.new(File.read('rsa.pem')).private_encrypt(user.password)
    - policy_name: ""
      policy_dsrid: DSR-7
      policy_display_id: ruby_lang_weak_encryption
      policy_description: Avoid weak encryption libraries.
      line_number: 5
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/openssl_rsa_data.rb
      category_groups:
        - PII
      parent_line_number: 5
      parent_content: rsa_encrypt.export(cipher, user.password)
    - policy_name: ""
      policy_dsrid: DSR-7
      policy_display_id: ruby_lang_weak_encryption
      policy_description: Avoid weak encryption libraries.
      line_number: 7
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/openssl_rsa_data.rb
      category_groups:
        - PII
      parent_line_number: 7
      parent_content: OpenSSL::PKey::RSA.new(2048).to_pem(cipher, user.first_name)


--

