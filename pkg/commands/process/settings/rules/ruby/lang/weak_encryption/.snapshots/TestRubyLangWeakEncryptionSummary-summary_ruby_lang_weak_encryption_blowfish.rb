warning:
    - policy_name: ""
      policy_dsrid: DSR-7
      policy_display_id: ruby_lang_weak_encryption
      policy_description: Avoid weak encryption libraries.
      line_number: 2
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/blowfish.rb
      parent_line_number: 2
      parent_content: |-
        blowfish.encrypt_block do
          "hello world"
        end


--

