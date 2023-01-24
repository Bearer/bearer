data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatypes_with_encrypted_jwt.rb
              line_number: 2
              field_name: email
              object_name: current_user
            - filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatypes_with_encrypted_jwt.rb
              line_number: 4
              field_name: email
              object_name: current_user
    - name: Fullname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatypes_with_encrypted_jwt.rb
              line_number: 6
              field_name: name
              object_name: user
risks:
    - detector_id: ruby_lang_jwt
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatypes_with_encrypted_jwt.rb
              line_number: 2
              parent:
                line_number: 2
                content: 'JWT.encode({ user: current_user.email }, private_key, ''HS256'', {})'
              field_name: email
              object_name: current_user
            - filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatypes_with_encrypted_jwt.rb
              line_number: 4
              parent:
                line_number: 4
                content: 'JWT.encode({ user: current_user.email }, ENV["SECRET_KEY"])'
              field_name: email
              object_name: current_user
        - name: Fullname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatypes_with_encrypted_jwt.rb
              line_number: 6
              parent:
                line_number: 6
                content: 'JWT.encode({ user_name: user.name }, Rails.application.secret_key_base)'
              field_name: name
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatypes_with_encrypted_jwt.rb
              line_number: 2
              parent:
                line_number: 2
                content: 'JWT.encode({ user: current_user.email }, private_key, ''HS256'', {})'
              object_name: current_user
            - filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatypes_with_encrypted_jwt.rb
              line_number: 4
              parent:
                line_number: 4
                content: 'JWT.encode({ user: current_user.email }, ENV["SECRET_KEY"])'
              object_name: current_user
            - filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatypes_with_encrypted_jwt.rb
              line_number: 6
              parent:
                line_number: 6
                content: 'JWT.encode({ user_name: user.name }, Rails.application.secret_key_base)'
              object_name: user
components: []


--

