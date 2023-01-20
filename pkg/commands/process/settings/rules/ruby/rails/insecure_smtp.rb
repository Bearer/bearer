# trigger_condition: application has sensitive data
class User
  attr_reader :name, :email, :password
end

# trigger: SSL::VERIFY_NONE verify mode for SMTP
Rails.application.configure do
  config.action_mailer.smtp_settings = {
    openssl_verify_mode: OpenSSL::SSL::VERIFY_NONE
  }
end

# trigger: "none" verify mode for SMTP
Rails.application.configure do
  config.action_mailer.smtp_settings = {
    openssl_verify_mode: "none"
  }
end

# ok: SSL::VERIFY_PEER verify mode for SMTP
Rails.application.configure do
  config.action_mailer.smtp_settings = {
    openssl_verify_mode: OpenSSL::SSL::VERIFY_PEER
  }
end

# ok: "peer" verify mode for SMTP
Rails.application.configure do
  config.action_mailer.smtp_settings = {
    openssl_verify_mode: "peer"
  }
end
