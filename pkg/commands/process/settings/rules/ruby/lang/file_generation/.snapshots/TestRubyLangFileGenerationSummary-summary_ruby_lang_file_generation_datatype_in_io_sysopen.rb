critical:
    - policy_name: ""
      policy_dsrid: DSR-4
      policy_display_id: ruby_lang_file_generation
      policy_description: Do not write sensitive data to static files.
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_io_sysopen.rb
      category_groups:
        - PII
      parent_line_number: 3
      parent_content: 'a.puts "Hello, #{user.full_name}!"'


--

