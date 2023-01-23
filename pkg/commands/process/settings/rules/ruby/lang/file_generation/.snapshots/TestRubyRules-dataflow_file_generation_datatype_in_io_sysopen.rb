data_types:
    - name: Fullname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_io_sysopen.rb
              line_number: 3
              field_name: full_name
              object_name: user_6
risks:
    - detector_id: ruby_lang_file_generation
      data_types:
        - name: Fullname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_io_sysopen.rb
              line_number: 3
              parent:
                line_number: 2
                content: |-
                    IO.open(fd,"w") do |a|
                      a.puts "Hello, #{user_6.full_name}!"
                    end
              field_name: full_name
              object_name: user_6
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_io_sysopen.rb
              line_number: 3
              parent:
                line_number: 2
                content: |-
                    IO.open(fd,"w") do |a|
                      a.puts "Hello, #{user_6.full_name}!"
                    end
              object_name: user_6
components: []


--

