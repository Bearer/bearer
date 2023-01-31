critical:
    - policy_name: ""
      policy_dsrid: DSR-7
      policy_display_id: ruby_lang_weak_encryption_with_data
      policy_description: Do not use weak encryption libraries to encrypt sensitive data.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/digest_sha1.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: Digest::SHA1.hexidigest(user.first_name)


--

