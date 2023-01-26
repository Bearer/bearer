low:
    - policy_name: ""
      policy_dsrid: DSR-2
      policy_display_id: ruby_lang_ssl_verification
      policy_description: Enable SSL Certificate Verification.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/lang/ssl_verification/testdata/verification_disabled.rb
      parent_line_number: 1
      parent_content: http.verify_mode = OpenSSL::SSL::VERIFY_NONE
    - policy_name: ""
      policy_dsrid: DSR-2
      policy_display_id: ruby_lang_ssl_verification
      policy_description: Enable SSL Certificate Verification.
      line_number: 4
      filename: pkg/commands/process/settings/rules/ruby/lang/ssl_verification/testdata/verification_disabled.rb
      parent_line_number: 4
      parent_content: |-
        Net::HTTP.start(uri.host, uri.port, :use_ssl => true, :verify_mode => OpenSSL::SSL::VERIFY_NONE) do |http|
          Net::HTTP::Get.new uri
        end


--

