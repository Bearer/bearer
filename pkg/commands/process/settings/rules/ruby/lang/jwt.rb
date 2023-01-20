# trigger:data type in JWT
JWT.encode user.address, nil, "none"
JWT.encode(user.email, nil, "none")

# trigger:data type in JWT
payload = { email: current_user.email }
JWT.encode(payload, ENV.fetch("SECRET_KEY"))

# trigger:data type in JWT
JWT.encode({ user: current_user.email }, ENV["SECRET_KEY"])

# trigger:data type in JWT
private_key = ENV.fetch("PRIVATE_JWT_KEY")
JWT.encode({ secret: "stuff", email: current_user.email }, private_key, 'HS256', {})

# trigger:data type in JWT
JWT.encode({ user_name: user.name }, Rails.application.secret_key_base)

# ok:undetectable data type
payload = {
	user: {
		first_name: "John",
		last_name: "Doe",
	}
}
token1 = JWT.encode payload, nil, "none"