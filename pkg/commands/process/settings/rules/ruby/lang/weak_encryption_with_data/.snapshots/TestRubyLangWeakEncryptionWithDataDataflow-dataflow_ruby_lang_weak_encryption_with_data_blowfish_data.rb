data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
              line_number: 5
              field_name: email
              object_name: user
    - name: Passwords
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
              line_number: 2
              field_name: password
              object_name: user
risks:
    - detector_id: ruby_lang_weak_encryption_with_data
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
              line_number: 5
              parent:
                line_number: 5
                content: Crypt::Blowfish.new("your-key").encrypt_string(user.email)
              field_name: email
              object_name: user
        - name: Passwords
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
              line_number: 2
              parent:
                line_number: 1
                content: |-
                    Crypt::Blowfish.new("insecure").encrypt_block do |user|
                      user.password
                    end
              field_name: password
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
              line_number: 2
              parent:
                line_number: 1
                content: |-
                    Crypt::Blowfish.new("insecure").encrypt_block do |user|
                      user.password
                    end
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
              line_number: 5
              parent:
                line_number: 5
                content: Crypt::Blowfish.new("your-key").encrypt_string(user.email)
              object_name: user
components: []


--

