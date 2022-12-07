class Users < ApplicationRecord
  # Detected
  session[:current_user] = user.email

  token = JWT.encode user.address, nil, "none"
  
  logger.info(user.first_name)
  # Not detected
  session[:user_name] = "mish bear"
end