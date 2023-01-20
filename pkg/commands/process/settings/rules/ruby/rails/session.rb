# trigger:data type in sessions
session[:current_user] = user.email
session[:user_id] = current_user.user_id

# ok:no detectable data type
session[:user_name] = "mish bear"