# critical risk: application has sensitive data
class User
  attr_reader :name, :email, :password
end

Rails.application.configure do
  config.force_ssl = false
end

