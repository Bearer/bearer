# critical risk: application has sensitive data
class User
  attr_reader :name, :email, :password
end

Rails.application.configure do
  config.action_mailer.smtp_settings = {
    openssl_verify_mode: "none"
  }
end