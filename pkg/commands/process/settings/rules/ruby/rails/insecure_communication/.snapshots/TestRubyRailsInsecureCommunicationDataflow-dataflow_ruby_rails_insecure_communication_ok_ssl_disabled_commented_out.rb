data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_communication/testdata/ok_ssl_disabled_commented_out.rb
              line_number: 3
              field_name: email
              object_name: User
    - name: Fullname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_communication/testdata/ok_ssl_disabled_commented_out.rb
              line_number: 3
              field_name: name
              object_name: User
    - name: Passwords
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_communication/testdata/ok_ssl_disabled_commented_out.rb
              line_number: 3
              field_name: password
              object_name: User
components: []


--

