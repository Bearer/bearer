risks:
    - detector_id: rc4_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/rc4_encrypt.rb
          line_number: 3
          parent:
            line_number: 3
            content: RC4.new("insecure")
          content: |
            RC4.new()
    - detector_id: ruby_lang_weak_encryption
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/rc4_encrypt.rb
          line_number: 1
          parent:
            line_number: 1
            content: RC4.new("insecure").encrypt("hello world")
          content: |
            $<VAR>.$<METHOD>($<_>)
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/rc4_encrypt.rb
          line_number: 4
          parent:
            line_number: 4
            content: rc4_encrypt.encrypt!("hello world")
          content: |
            $<VAR>.$<METHOD>($<_>)
components: []


--

