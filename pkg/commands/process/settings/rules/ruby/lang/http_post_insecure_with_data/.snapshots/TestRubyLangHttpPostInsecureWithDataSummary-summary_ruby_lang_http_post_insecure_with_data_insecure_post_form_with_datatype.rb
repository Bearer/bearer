critical:
    - rule_dsrid: DSR-2
      rule_display_id: ruby_lang_http_post_insecure_with_data
      rule_description: Only send sensitive data through HTTPS connections.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_http_post_insecure_with_data
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/lang/http_post_insecure_with_data/testdata/insecure_post_form_with_datatype.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: 'Net::HTTP.post_form("http://my.api.com/users/search", email: user.email)'


--

