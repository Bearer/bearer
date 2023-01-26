risks:
    - detector_id: blowfish_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/blowfish.rb
          line_number: 1
          parent:
            line_number: 1
            content: Crypt::Blowfish.new("insecure")
          content: |
            Crypt::Blowfish.new()
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/blowfish.rb
          line_number: 1
          parent:
            line_number: 1
            content: Crypt::Blowfish.new("insecure")
          content: |
            Crypt::Blowfish.new()
    - detector_id: ruby_lang_weak_encryption
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/blowfish.rb
          line_number: 2
          parent:
            line_number: 2
            content: |-
                blowfish.encrypt_block do
                  "hello world"
                end
          content: |
            $<VAR>.$<METHOD> do
              $<_>
            end
components: []


--

