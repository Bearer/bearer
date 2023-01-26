data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/cookies/testdata/ok_encrypted_cookies.rb
              line_number: 2
              field_name: email
              object_name: user
    - name: Physical Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/cookies/testdata/ok_encrypted_cookies.rb
              line_number: 1
              field_name: address
              object_name: user
components: []


--

