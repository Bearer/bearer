# Insecure SMTP

class User
  attr_reader :name, :email, :password
end

## Detected
Rails.application.configure do
  config.action_mailer.smtp_settings = {
    openssl_verify_mode: OpenSSL::SSL::VERIFY_NONE
  }
end

Rails.application.configure do
  config.action_mailer.smtp_settings = {
    openssl_verify_mode: "none"
  }
end

## Not Detected
Rails.application.configure do
  config.action_mailer.smtp_settings = {
    openssl_verify_mode: OpenSSL::SSL::VERIFY_PEER
  }
end

Rails.application.configure do
  config.action_mailer.smtp_settings = {
    openssl_verify_mode: "peer"
  }
end
