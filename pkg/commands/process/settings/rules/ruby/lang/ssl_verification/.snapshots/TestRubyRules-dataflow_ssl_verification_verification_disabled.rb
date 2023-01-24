data_types:
    - name: Gender identity
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/ssl_verification/testdata/verification_disabled.rb
              line_number: 2
              field_name: gender_identity
              object_name: user
risks:
    - detector_id: ruby_lang_ssl_verification
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/ssl_verification/testdata/verification_disabled.rb
          line_number: 4
          parent:
            line_number: 4
            content: http.verify_mode = OpenSSL::SSL::VERIFY_NONE
          content: |
            $<_>.verify_mode = OpenSSL::SSL::VERIFY_NONE
components: []


--

