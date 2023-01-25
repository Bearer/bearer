user = { address: "foo" }
NewRelic::Agent.add_custom_parameters(user)
NewRelic::Agent.add_custom_parameters(user: { email: "user@example.com" })
