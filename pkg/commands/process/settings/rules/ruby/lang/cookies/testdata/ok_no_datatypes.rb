# no (detectable) data types
cookies[:user_name] = "david"
cookies[:lat_lon] = JSON.generate([47.68, -122.37])

cookies.signed[:user_email] = "mish@bearer.sh"
cookies.permanent[:first_name] = "John"

cookies.signed.permanent[:user] = JSON.generate({first_name: "John", last_name: "Doe"})
cookies.signed.permanent[:user] = {
	first_name: "John",
	last_name: "Doe"
}.to_json