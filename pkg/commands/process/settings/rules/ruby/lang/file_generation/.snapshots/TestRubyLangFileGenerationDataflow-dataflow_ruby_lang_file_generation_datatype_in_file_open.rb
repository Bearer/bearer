data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_file_open.rb
              line_number: 1
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_file_open.rb
              line_number: 5
              field_name: email
              object_name: user
              subject_name: User
    - name: Emails
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_file_open.rb
              line_number: 3
              field_name: emails
              object_name: user
              subject_name: User
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_file_open.rb
              line_number: 5
              field_name: first_name
              object_name: user
              subject_name: User
    - name: Lastname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_file_open.rb
              line_number: 5
              field_name: last_name
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_lang_file_generation
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_file_open.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'f.write "#{Time.now} - User #{user.email} logged in\n"'
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_file_open.rb
              line_number: 5
              parent:
                line_number: 5
                content: f.write "#{user.email},#{user.first_name},#{user.last_name}"
              field_name: email
              object_name: user
              subject_name: User
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_file_open.rb
              line_number: 5
              parent:
                line_number: 5
                content: f.write "#{user.email},#{user.first_name},#{user.last_name}"
              field_name: first_name
              object_name: user
              subject_name: User
        - name: Lastname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_file_open.rb
              line_number: 5
              parent:
                line_number: 5
                content: f.write "#{user.email},#{user.first_name},#{user.last_name}"
              field_name: last_name
              object_name: user
              subject_name: User
    - detector_id: ruby_lang_file_generation_file
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_file_open.rb
          line_number: 1
          parent:
            line_number: 1
            content: f
          content: |
            File.open() { |$<!>$<_:identifier>| }
        - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_file_open.rb
          line_number: 3
          parent:
            line_number: 3
            content: f
          content: |
            File.open() { |$<!>$<_:identifier>| }
components: []


--

