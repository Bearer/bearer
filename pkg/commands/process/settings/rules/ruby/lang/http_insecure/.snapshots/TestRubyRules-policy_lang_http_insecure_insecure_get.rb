low:
    - policy_name: ""
      policy_dsrid: DSR-2
      policy_display_id: ruby_lang_http_insecure
      policy_description: Do not perform unsecure HTTP calls.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/lang/http_insecure/testdata/insecure_get.rb
      parent_line_number: 1
      parent_content: Faraday.get("http://api.insecure.com")


--

