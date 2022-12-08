class Users < ApplicationRecord
  # Detected
  session[:current_user] = user.email

  # should enrich detection with it since it is part of user.email
  user.email.domain = "gmail.com"

  # should ignore since none are part of user.email
  token = JWT.encode user.address, nil, "none"
  logger.info(user.first_name)
  session[:user_name] = "mish bear"
  admin.username = "admin"
end