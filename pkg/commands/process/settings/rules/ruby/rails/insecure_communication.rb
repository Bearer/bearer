# trigger_condition: application has sensitive data
class User
  attr_reader :name, :email, :password
end

# trigger:ssl disabled
Rails.application.configure do
  config.force_ssl = false
end

# ok:ssl enabled
Rails.application.configure do
  config.force_ssl = true
end

# ok:commented out code
Rails.application.configure do
  # config.force_ssl = false
end
