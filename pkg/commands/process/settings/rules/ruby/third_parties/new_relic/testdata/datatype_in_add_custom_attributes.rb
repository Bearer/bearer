user = { email: "user@example.com" }
NewRelic::Agent.add_custom_attributes(user)
NewRelic::Agent.add_custom_attributes(a: "test", user: { email: "user@example.com" }, other: 42)
