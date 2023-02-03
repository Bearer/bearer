critical:
    - policy_name: ""
      policy_dsrid: DSR-5
      policy_display_id: ruby_lang_exception
      policy_description: Do not send sensitive data to exceptions.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/lang/exception/testdata/datatype_leak.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: raise CustomException.new(user.email)
    - policy_name: ""
      policy_dsrid: DSR-5
      policy_display_id: ruby_lang_exception
      policy_description: Do not send sensitive data to exceptions.
      line_number: 2
      filename: pkg/commands/process/settings/rules/ruby/lang/exception/testdata/datatype_leak.rb
      category_groups:
        - PII
      parent_line_number: 2
      parent_content: 'raise "User doesn''t exist #{user.email}"'


--

