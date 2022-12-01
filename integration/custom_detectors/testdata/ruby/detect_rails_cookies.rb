# Detected
cookies[:user_name] = user.name
cookies.signed[:user_email] = current_user.email
cookies.encrypted[:full_name] = customer.name
cookies.permanent[:first_name] = current_user.first_name
cookies.signed.permanent[:user_email] = current_user.email
cookies.permanent.encrypted[:full_name] = customer.name
cookies.signed.permanent[:first_name] = current_user.first_name

# Not detected
cookies[:user_name] = "david"
cookies.signed[:user_email] = "mish@bearer.sh"
cookies.encrypted[:full_name] = "John Doe"
cookies.permanent[:first_name] = "John"
cookies[:lat_lon] = JSON.generate([47.68, -122.37])

# Not detected (yet - variable reconciliation)
# user = {
#   first_name: "John",
# 	last_name: "Doe"
# }
# cookies[:login] = { value: user.to_json, expires: 1.hour, secure: true }
# cookies[:login] = { value: user.login, expires: 1.hour, secure: true }

cookies.signed.permanent[:user] = JSON.generate({first_name: "John", last_name: "Doe"})
cookies.signed.permanent[:user] = {
	first_name: "John",
	last_name: "Doe"
}.to_json