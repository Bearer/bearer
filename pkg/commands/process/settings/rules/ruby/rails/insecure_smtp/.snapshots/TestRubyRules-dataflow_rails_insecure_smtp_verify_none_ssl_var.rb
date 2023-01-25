risks:
    - detector_id: ruby_rails_insecure_smtp
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_smtp/testdata/verify_none_ssl_var.rb
          line_number: 3
          parent:
            line_number: 3
            content: 'openssl_verify_mode: OpenSSL::SSL::VERIFY_NONE'
          content: |
            Rails.application.configure do
              config.action_mailer.smtp_settings = {
                $<!>openssl_verify_mode: OpenSSL::SSL::VERIFY_NONE
              }
            end
components: []


--

