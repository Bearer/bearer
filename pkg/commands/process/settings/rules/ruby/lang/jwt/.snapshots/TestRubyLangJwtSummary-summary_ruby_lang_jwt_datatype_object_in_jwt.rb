critical:
    - policy_name: ""
      policy_dsrid: DSR-3
      policy_display_id: ruby_lang_jwt
      policy_description: Do not store sensitive data in JWTs.
      line_number: 2
      filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatype_object_in_jwt.rb
      category_groups:
        - PII
      parent_line_number: 6
      parent_content: JWT.encode(payload, ENV.fetch("SECRET_KEY"))
    - policy_name: ""
      policy_dsrid: DSR-3
      policy_display_id: ruby_lang_jwt
      policy_description: Do not store sensitive data in JWTs.
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatype_object_in_jwt.rb
      category_groups:
        - PII
      parent_line_number: 6
      parent_content: JWT.encode(payload, ENV.fetch("SECRET_KEY"))


--

