warning:
    - policy_name: ""
      policy_dsrid: DSR-7
      policy_display_id: ruby_lang_weak_encryption
      policy_description: Avoid weak encryption libraries.
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/openssl_dsa.rb
      parent_line_number: 3
      parent_content: dsa_encrypt.export(cipher, "hello world")
    - policy_name: ""
      policy_dsrid: DSR-7
      policy_display_id: ruby_lang_weak_encryption
      policy_description: Avoid weak encryption libraries.
      line_number: 5
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/openssl_dsa.rb
      parent_line_number: 5
      parent_content: OpenSSL::PKey::DSA.new(2048).to_pem(cipher, "hello world")


--

