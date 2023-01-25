data_types:
    - name: Physical Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatype_in_jwt.rb
              line_number: 1
              field_name: address
              object_name: user
risks:
    - detector_id: ruby_lang_jwt
      data_types:
        - name: Physical Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatype_in_jwt.rb
              line_number: 1
              parent:
                line_number: 1
                content: JWT.encode user.address, nil, "none"
              field_name: address
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatype_in_jwt.rb
              line_number: 1
              parent:
                line_number: 1
                content: JWT.encode user.address, nil, "none"
              object_name: user
components: []


--

