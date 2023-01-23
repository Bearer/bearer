critical:
    - policy_name: ""
      policy_dsrid: DSR-4
      policy_display_id: ruby_lang_file_generation
      policy_description: Do not write sensitive data to static files.
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_io_sysopen.rb
      category_groups:
        - PII
      parent_line_number: 2
      parent_content: |-
        IO.open(fd,"w") do |a|
          a.puts "Hello, #{user_6.full_name}!"
        end


--

