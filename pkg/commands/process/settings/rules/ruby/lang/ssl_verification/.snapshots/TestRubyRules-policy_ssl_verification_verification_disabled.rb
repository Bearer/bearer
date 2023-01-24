medium:
    - policy_name: ""
      policy_dsrid: DSR-2
      policy_display_id: ruby_lang_ssl_verification
      policy_description: Do not disable SSL certificate verification.
      line_number: 4
      filename: pkg/commands/process/settings/rules/ruby/lang/ssl_verification/testdata/verification_disabled.rb
      category_groups:
        - PII
        - Personal Data (Sensitive)
      parent_line_number: 4
      parent_content: http.verify_mode = OpenSSL::SSL::VERIFY_NONE


--

