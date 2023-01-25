critical:
    - policy_name: ""
      policy_dsrid: DSR-2
      policy_display_id: ruby_lang_http_post_insecure_with_data
      policy_description: Do not send sensitive data through unsecure HTTP calls.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/lang/http_post_insecure_with_data/testdata/insecure_post_form_with_datatype.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: 'Net::HTTP.post_form("http://my.api.com/users/search", email: user.email)'
low:
    - policy_name: ""
      policy_dsrid: DSR-2
      policy_display_id: ruby_lang_http_insecure
      policy_description: Do not perform unsecure HTTP calls.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/lang/http_post_insecure_with_data/testdata/insecure_post_form_with_datatype.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: 'Net::HTTP.post_form("http://my.api.com/users/search", email: user.email)'


--

