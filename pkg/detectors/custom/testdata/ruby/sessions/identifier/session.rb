# Detected
session[:current_user] = user

# should enrich detection with it since it is part of user
user.email.domain = "gmail.com"

# should ignore since none are part of user
admin.username = "admin"
