# Insecure communication
## Detected
Rails.application.configure do
  config.force_ssl = false
end

## Not Detected
Rails.application.configure do
  config.force_ssl = true
end

Rails.application.configure do
  # config.force_ssl = false
end
