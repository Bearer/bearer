# detected
cookies[:user_name] = user.name
# Not detected
cookies[:user_name] = "david"

# cookies[:lat_lon] = JSON.generate([47.68, -122.37])

# user = {
#   first_name: "John",
# 	last_name: "Doe"
# }
# cookies[:login] = { value: user.to_json, expires: 1.hour, secure: true }
# cookies[:login] = { value: user.login, expires: 1.hour, secure: true }

# # Signed
# cookies.signed[:user_id] = current_user.id

# # Encrypted
# cookies.encrypted[:full_name] = "John Doe"

# # Permanent
# cookies.permanent[:first_name] = "John"

# # Permanent & Encrypted
# cookies.signed.permanent[:user] = JSON.generate({first_name: "John", last_name: "Doe"})
# cookies.signed.permanent[:user] = {
# 	first_name: "John",
# 	last_name: "Doe"
# }.to_json