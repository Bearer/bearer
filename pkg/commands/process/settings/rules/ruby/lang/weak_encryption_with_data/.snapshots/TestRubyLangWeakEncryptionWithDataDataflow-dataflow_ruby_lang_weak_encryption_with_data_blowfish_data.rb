data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
              line_number: 9
              field_name: email
              object_name: user
              subject_name: User
    - name: Passwords
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
              line_number: 2
              field_name: password
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
              line_number: 6
              field_name: password
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_lang_weak_encryption_with_data
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
              line_number: 9
              parent:
                line_number: 9
                content: Crypt::Blowfish.new("your-key").encrypt_string(user.email)
              field_name: email
              object_name: user
              subject_name: User
        - name: Passwords
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
              line_number: 2
              parent:
                line_number: 1
                content: |-
                    Crypt::Blowfish.new("insecure").encrypt_block { |user|
                      user.password
                    }
              field_name: password
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
              line_number: 6
              parent:
                line_number: 5
                content: |-
                    Crypt::Blowfish.new("insecure").encrypt_block do |user|
                      user.password
                    end
              field_name: password
              object_name: user
              subject_name: User
components: []


--

