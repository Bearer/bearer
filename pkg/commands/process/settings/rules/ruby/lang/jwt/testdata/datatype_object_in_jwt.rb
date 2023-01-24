payload = {
	user: {
		email: user.email
	}
}
JWT.encode(payload, ENV.fetch("SECRET_KEY"))
