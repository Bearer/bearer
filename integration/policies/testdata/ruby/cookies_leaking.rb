# Detected
cookies.signed[:info] = user.email
cookies.permanent.encrypted[:secret] = user.address
user_1 = {
  first_name: "John",
  last_name: "Doe"
}
cookies[:login] = { value: user_1.to_json, expires: 1.hour, secure: true }
cookies.permanent[:user_id] = current_user.user_id

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
