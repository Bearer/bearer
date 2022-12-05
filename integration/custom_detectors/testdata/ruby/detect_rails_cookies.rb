# Detected
cookies[:user_name] = user.name
cookies.signed[:user_email] = user.email
cookies.encrypted[:full_name] = user.name
cookies.permanent[:first_name] = user.first_name
cookies.signed.permanent[:user_email] = user.email
cookies.permanent.encrypted[:full_name] = user.name
cookies.signed.permanent[:first_name] = user.first_name

user_1 = {
  first_name: "John",
	last_name: "Doe"
}
cookies[:login] = { value: user_1.to_json, expires: 1.hour, secure: true }

# Not detected
cookies[:user_name] = "david"
cookies.signed[:user_email] = "mish@bearer.sh"
cookies.encrypted[:full_name] = "John Doe"
cookies.permanent[:first_name] = "John"
cookies[:lat_lon] = JSON.generate([47.68, -122.37])

cookies.signed.permanent[:user] = JSON.generate({first_name: "John", last_name: "Doe"})
cookies.signed.permanent[:user] = {
	first_name: "John",
	last_name: "Doe"
}.to_json