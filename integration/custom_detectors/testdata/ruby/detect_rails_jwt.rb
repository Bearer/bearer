payload = {
	user: {
		first_name: "John",
		last_name: "Doe",
	}
}

token1 = JWT.encode payload, nil, "none"

token2 = JWT.encode user.address, nil, "none"

token3 = JWT.encode(user.email, nil, "none")