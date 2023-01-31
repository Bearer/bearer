data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/openssl_dsa_data.rb
              line_number: 3
              field_name: email
              object_name: user
              subject_name: User
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/openssl_dsa_data.rb
              line_number: 5
              field_name: first_name
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_lang_weak_encryption_with_data
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/openssl_dsa_data.rb
              line_number: 3
              parent:
                line_number: 3
                content: dsa_encrypt.export(cipher, user.email)
              field_name: email
              object_name: user
              subject_name: User
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/openssl_dsa_data.rb
              line_number: 5
              parent:
                line_number: 5
                content: OpenSSL::PKey::RSA.new(2048).to_pem(cipher, user.first_name)
              field_name: first_name
              object_name: user
              subject_name: User
    - detector_id: openssl_dsa_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/openssl_dsa_data.rb
          line_number: 2
          parent:
            line_number: 2
            content: OpenSSL::PKey::DSA.new(2048)
          content: |
            OpenSSL::PKey::DSA.new()
components: []


--

