data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_open.rb
              line_number: 5
              field_name: email
              object_name: user
              subject_name: User
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_open.rb
              line_number: 6
              field_name: first_name
              object_name: user
              subject_name: User
    - name: Lastname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_open.rb
              line_number: 7
              field_name: last_name
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_lang_file_generation_file
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_open.rb
          line_number: 1
          parent:
            line_number: 1
            content: csv
          content: |
            CSV.open() { |$<!>$<_:identifier>| }
components: []


--

