critical:
    - policy_name: ""
      policy_dsrid: DSR-3
      policy_display_id: ruby_lang_jwt
      policy_description: Do not store sensitive data in JWTs.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatype_in_jwt.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: JWT.encode user.address, nil, "none"


--

