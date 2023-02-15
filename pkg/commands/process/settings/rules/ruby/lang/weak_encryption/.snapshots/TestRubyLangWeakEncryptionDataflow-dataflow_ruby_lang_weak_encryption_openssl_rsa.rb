risks:
    - detector_id: openssl_rsa_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/openssl_rsa.rb
          line_number: 1
          parent:
            line_number: 1
            content: OpenSSL::PKey::RSA.new(File.read('rsa.pem'))
          content: |
            OpenSSL::PKey::RSA.new()
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/openssl_rsa.rb
          line_number: 4
          parent:
            line_number: 4
            content: OpenSSL::PKey::RSA.new(2048)
          content: |
            OpenSSL::PKey::RSA.new()
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/openssl_rsa.rb
          line_number: 7
          parent:
            line_number: 7
            content: OpenSSL::PKey::RSA.new(2048)
          content: |
            OpenSSL::PKey::RSA.new()
    - detector_id: ruby_lang_weak_encryption
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/openssl_rsa.rb
          line_number: 1
          parent:
            line_number: 1
            content: OpenSSL::PKey::RSA.new(File.read('rsa.pem')).private_encrypt("test")
          content: |
            $<VAR>.$<METHOD>($<_>)
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/openssl_rsa.rb
          line_number: 5
          parent:
            line_number: 5
            content: rsa_encrypt.export(cipher, "hello world")
          content: |
            $<VAR>.$<METHOD>($<_>)
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/openssl_rsa.rb
          line_number: 7
          parent:
            line_number: 7
            content: OpenSSL::PKey::RSA.new(2048).to_pem(cipher, "hello world")
          content: |
            $<VAR>.$<METHOD>($<_>)
components: []


--

