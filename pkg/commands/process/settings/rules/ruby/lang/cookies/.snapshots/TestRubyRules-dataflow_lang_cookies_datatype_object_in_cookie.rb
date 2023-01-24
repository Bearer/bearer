data_types:
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/cookies/testdata/datatype_object_in_cookie.rb
              line_number: 2
              field_name: first_name
              object_name: user
    - name: Lastname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/cookies/testdata/datatype_object_in_cookie.rb
              line_number: 3
              field_name: last_name
              object_name: user
risks:
    - detector_id: ruby_lang_cookies
      data_types:
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/cookies/testdata/datatype_object_in_cookie.rb
              line_number: 2
              parent:
                line_number: 5
                content: 'cookies[:login] = { value: user.to_json, expires: 1.hour, secure: true }'
              field_name: first_name
              object_name: user
        - name: Lastname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/cookies/testdata/datatype_object_in_cookie.rb
              line_number: 3
              parent:
                line_number: 5
                content: 'cookies[:login] = { value: user.to_json, expires: 1.hour, secure: true }'
              field_name: last_name
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/cookies/testdata/datatype_object_in_cookie.rb
              line_number: 1
              parent:
                line_number: 5
                content: 'cookies[:login] = { value: user.to_json, expires: 1.hour, secure: true }'
              object_name: user
components: []


--

