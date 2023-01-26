risks:
    - detector_id: openssl_dsa_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/openssl_dsa.rb
          line_number: 2
          parent:
            line_number: 2
            content: OpenSSL::PKey::DSA.new(2048)
          content: |
            OpenSSL::PKey::DSA.new()
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/openssl_dsa.rb
          line_number: 2
          parent:
            line_number: 2
            content: OpenSSL::PKey::DSA.new(2048)
          content: |
            OpenSSL::PKey::DSA.new()
    - detector_id: ruby_lang_weak_encryption
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/openssl_dsa.rb
          line_number: 3
          parent:
            line_number: 3
            content: dsa_encrypt.export(cipher, "hello world")
          content: |
            $<VAR>.$<METHOD>($<_>)
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/openssl_dsa.rb
          line_number: 5
          parent:
            line_number: 5
            content: OpenSSL::PKey::DSA.new(2048).to_pem(cipher, "hello world")
          content: |
            $<VAR>.$<METHOD>($<_>)
components: []


--

