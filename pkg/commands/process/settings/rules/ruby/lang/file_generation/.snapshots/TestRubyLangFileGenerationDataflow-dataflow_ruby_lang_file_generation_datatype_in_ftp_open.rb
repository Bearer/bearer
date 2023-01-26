data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 1
              field_name: email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              field_name: email
              object_name: user
    - name: Emails
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 3
              field_name: emails
              object_name: user
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              field_name: first_name
              object_name: user
    - name: Lastname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              field_name: last_name
              object_name: user
components: []


--

