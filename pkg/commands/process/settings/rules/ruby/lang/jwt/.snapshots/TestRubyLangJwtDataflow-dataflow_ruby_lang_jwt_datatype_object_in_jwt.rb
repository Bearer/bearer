data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatype_object_in_jwt.rb
              line_number: 3
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_lang_jwt
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatype_object_in_jwt.rb
              line_number: 3
              parent:
                line_number: 6
                content: JWT.encode(payload, ENV.fetch("SECRET_KEY"))
              field_name: email
              object_name: user
              subject_name: User
components: []


--

