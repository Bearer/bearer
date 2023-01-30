data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_smtp/testdata/ok_ssl_verify_peer_mode.rb
              line_number: 3
              field_name: email
              object_name: User
              subject_name: User
    - name: Fullname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_smtp/testdata/ok_ssl_verify_peer_mode.rb
              line_number: 3
              field_name: name
              object_name: User
              subject_name: User
    - name: Passwords
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_smtp/testdata/ok_ssl_verify_peer_mode.rb
              line_number: 3
              field_name: password
              object_name: User
              subject_name: User
components: []


--

