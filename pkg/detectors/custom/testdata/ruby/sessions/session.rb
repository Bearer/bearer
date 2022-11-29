class Users < ApplicationRecord
  # Detected
  session[:current_user] = user.email
  # Not detected
  session[:user_name] = "mish bear"
end
