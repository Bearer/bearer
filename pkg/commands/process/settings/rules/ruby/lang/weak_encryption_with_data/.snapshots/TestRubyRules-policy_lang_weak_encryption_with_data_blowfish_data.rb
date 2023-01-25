critical:
    - policy_name: ""
      policy_dsrid: DSR-7
      policy_display_id: ruby_lang_weak_encryption
      policy_description: Do not use weak encryption libraries.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: |-
        Crypt::Blowfish.new("insecure").encrypt_block do |user|
          user.password
        end
    - policy_name: ""
      policy_dsrid: DSR-7
      policy_display_id: ruby_lang_weak_encryption_with_data
      policy_description: Do not use weak encryption libraries to encrypt sensitive data.
      line_number: 2
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/blowfish_data.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: |-
        Crypt::Blowfish.new("insecure").encrypt_block do |user|
          user.password
        end


--

