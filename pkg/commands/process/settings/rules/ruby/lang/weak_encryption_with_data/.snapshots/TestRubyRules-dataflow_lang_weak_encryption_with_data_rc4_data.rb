data_types:
    - name: Passwords
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/rc4_data.rb
              line_number: 1
              field_name: password
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/rc4_data.rb
              line_number: 4
              field_name: password
              object_name: user
risks:
    - detector_id: ruby_lang_weak_encryption_with_data
      data_types:
        - name: Passwords
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/rc4_data.rb
              line_number: 1
              parent:
                line_number: 1
                content: RC4.new("insecure").encrypt(user.password)
              field_name: password
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/rc4_data.rb
              line_number: 4
              parent:
                line_number: 4
                content: rc4_encrypt.encrypt!(user.password)
              field_name: password
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/rc4_data.rb
              line_number: 1
              parent:
                line_number: 1
                content: RC4.new("insecure").encrypt(user.password)
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/rc4_data.rb
              line_number: 4
              parent:
                line_number: 4
                content: rc4_encrypt.encrypt!(user.password)
              object_name: user
    - detector_id: rc4_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/rc4_data.rb
          line_number: 3
          parent:
            line_number: 3
            content: RC4.new("insecure")
          content: |
            RC4.new()
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/rc4_data.rb
          line_number: 3
          parent:
            line_number: 3
            content: RC4.new("insecure")
          content: |
            RC4.new()
    - detector_id: ruby_lang_weak_encryption
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/rc4_data.rb
          line_number: 1
          parent:
            line_number: 1
            content: RC4.new("insecure").encrypt(user.password)
          content: |
            $<VAR>.$<METHOD>($<_>)
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/rc4_data.rb
          line_number: 4
          parent:
            line_number: 4
            content: rc4_encrypt.encrypt!(user.password)
          content: |
            $<VAR>.$<METHOD>($<_>)
components: []


--

