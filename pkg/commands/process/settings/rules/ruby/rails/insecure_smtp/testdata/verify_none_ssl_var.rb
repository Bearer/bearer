Rails.application.configure do
  config.action_mailer.smtp_settings = {
    :openssl_verify_mode => OpenSSL::SSL::VERIFY_NONE
  }
end
