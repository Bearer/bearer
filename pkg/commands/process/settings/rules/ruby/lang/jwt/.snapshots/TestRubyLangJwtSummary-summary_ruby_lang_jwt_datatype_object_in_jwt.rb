critical:
    - rule_dsrid: DSR-3
      rule_display_id: ruby_lang_jwt
      rule_description: Do not store sensitive data in JWTs.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_jwt
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/lang/jwt/testdata/datatype_object_in_jwt.rb
      category_groups:
        - PII
      parent_line_number: 6
      parent_content: JWT.encode(payload, ENV.fetch("SECRET_KEY"))


--

