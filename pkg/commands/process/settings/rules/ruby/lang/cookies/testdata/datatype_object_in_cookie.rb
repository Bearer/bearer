user = {
  first_name: "John",
	last_name: "Doe"
}
cookies[:login] = { value: user.to_json, expires: 1.hour, secure: true }