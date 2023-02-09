low:
    - rule_dsrid: DSR-2
      rule_display_id: ruby_lang_ssl_verification
      rule_description: Enable SSL Certificate Verification.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_ssl_verification
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/lang/ssl_verification/testdata/verification_disabled.rb
      parent_line_number: 1
      parent_content: http.verify_mode = OpenSSL::SSL::VERIFY_NONE
    - rule_dsrid: DSR-2
      rule_display_id: ruby_lang_ssl_verification
      rule_description: Enable SSL Certificate Verification.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_ssl_verification
      line_number: 4
      filename: pkg/commands/process/settings/rules/ruby/lang/ssl_verification/testdata/verification_disabled.rb
      parent_line_number: 4
      parent_content: |-
        Net::HTTP.start(uri.host, uri.port, :use_ssl => true, :verify_mode => OpenSSL::SSL::VERIFY_NONE) do |http|
          Net::HTTP::Get.new uri
        end


--

