critical:
    - policy_name: ""
      policy_dsrid: DSR-3
      policy_display_id: ruby_lang_jwt
      policy_description: Do not store sensitive data in JWTs.
      line_number: 2
      filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatypes_with_encrypted_jwt.rb
      category_groups:
        - PII
      parent_line_number: 2
      parent_content: 'JWT.encode({ user: current_user.email }, private_key, ''HS256'', {})'
    - policy_name: ""
      policy_dsrid: DSR-3
      policy_display_id: ruby_lang_jwt
      policy_description: Do not store sensitive data in JWTs.
      line_number: 4
      filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatypes_with_encrypted_jwt.rb
      category_groups:
        - PII
      parent_line_number: 4
      parent_content: 'JWT.encode({ user: current_user.email }, ENV["SECRET_KEY"])'
    - policy_name: ""
      policy_dsrid: DSR-3
      policy_display_id: ruby_lang_jwt
      policy_description: Do not store sensitive data in JWTs.
      line_number: 6
      filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatypes_with_encrypted_jwt.rb
      category_groups:
        - PII
      parent_line_number: 6
      parent_content: 'JWT.encode({ user_name: user.name }, Rails.application.secret_key_base)'


--

