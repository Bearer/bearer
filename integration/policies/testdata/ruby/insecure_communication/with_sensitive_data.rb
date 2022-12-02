# Insecure communication

class User
  attr_reader :name, :email, :password, :ethnicity
end

# Should match
Rails.application.configure do
  config.force_ssl = false
end

# Should not match
Rails.application.configure do
  config.force_ssl = true
end

Rails.application.configure do
  # config.force_ssl = false
end
