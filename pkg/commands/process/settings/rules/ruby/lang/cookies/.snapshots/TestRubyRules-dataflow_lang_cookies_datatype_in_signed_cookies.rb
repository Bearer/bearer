data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/cookies/testdata/datatype_in_signed_cookies.rb
              line_number: 1
              field_name: email
              object_name: user
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/cookies/testdata/datatype_in_signed_cookies.rb
              line_number: 2
              field_name: first_name
              object_name: user
risks:
    - detector_id: ruby_lang_cookies
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/cookies/testdata/datatype_in_signed_cookies.rb
              line_number: 1
              parent:
                line_number: 1
                content: cookies.signed[:info] = user.email
              field_name: email
              object_name: user
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/cookies/testdata/datatype_in_signed_cookies.rb
              line_number: 2
              parent:
                line_number: 2
                content: cookies.permanent.signed[:secret] = user.first_name
              field_name: first_name
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/cookies/testdata/datatype_in_signed_cookies.rb
              line_number: 1
              parent:
                line_number: 1
                content: cookies.signed[:info] = user.email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/lang/cookies/testdata/datatype_in_signed_cookies.rb
              line_number: 2
              parent:
                line_number: 2
                content: cookies.permanent.signed[:secret] = user.first_name
              object_name: user
components: []


--

