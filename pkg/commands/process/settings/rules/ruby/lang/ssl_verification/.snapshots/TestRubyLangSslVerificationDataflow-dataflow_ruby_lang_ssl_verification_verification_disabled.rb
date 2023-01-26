risks:
    - detector_id: ruby_lang_ssl_verification
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/ssl_verification/testdata/verification_disabled.rb
          line_number: 1
          parent:
            line_number: 1
            content: http.verify_mode = OpenSSL::SSL::VERIFY_NONE
          content: |
            $<_>.verify_mode = OpenSSL::SSL::VERIFY_NONE
        - filename: pkg/commands/process/settings/rules/ruby/lang/ssl_verification/testdata/verification_disabled.rb
          line_number: 4
          parent:
            line_number: 4
            content: |-
                Net::HTTP.start(uri.host, uri.port, :use_ssl => true, :verify_mode => OpenSSL::SSL::VERIFY_NONE) do |http|
                  Net::HTTP::Get.new uri
                end
          content: |
            Net::HTTP.start(:verify_mode => OpenSSL::SSL::VERIFY_NONE)$<...>
components: []


--

