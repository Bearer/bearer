data_types:
    - name: Passwords
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/password_length/testdata/ok_password_length.rb
              line_number: 3
              field_name: password
              object_name: Student
              subject_name: User
components: []


--

