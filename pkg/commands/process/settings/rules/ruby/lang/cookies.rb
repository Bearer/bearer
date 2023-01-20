# trigger:data type
session[:current_user] = user.email
# trigger:data type
session[:user_id] = current_user.user_id
# trigger:data type
cookies.signed[:info] = user.email
# trigger:data type
cookies.permanent.signed[:secret] = user.address
# trigger:object with data type(s)
user_1 = {
  first_name: "John",
	last_name: "Doe"
}
cookies[:login] = { value: user_1.to_json, expires: 1.hour, secure: true }
# trigger:data type
cookies.permanent[:user_id] = current_user.user_id

# ok:no data type
session[:user_name] = "mish bear"
cookies[:user_name] = "david"
cookies[:lat_lon] = JSON.generate([47.68, -122.37])
# ok:encrypted
cookies.permanent.encrypted[:secret] = user.address
cookies.encrypted[:full_name] = "John Doe"
# ok:no detectable data type
cookies.signed[:user_email] = "mish@bearer.sh"
cookies.permanent[:first_name] = "John"
# ok:object with no detectable data type(s)
cookies.signed.permanent[:user] = JSON.generate({first_name: "John", last_name: "Doe"})
cookies.signed.permanent[:user] = {
	first_name: "John",
	last_name: "Doe"
}.to_json